package hikaku

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type pass struct {
	level  int
	path   string
	index  int
	parent string
}

type valueDiff struct{}

type differenceContext []valueDiff

func newDifferenceContext() *differenceContext {
	return &differenceContext{}
}

type executionBuffer []func() error

func (e *executionBuffer) Add(cb func() error) {
	(*e) = append((*e), cb)
}

func (e *executionBuffer) Pop() func() error {
	var x func() error
	x, (*e) = (*e)[0], (*e)[1:]
	return x
}

func (e *executionBuffer) Len() int {
	return len(*e)
}

func newExecutionBuffer() *executionBuffer {
	return &executionBuffer{}
}

type AttributeData struct {
	Name       string
	Path       PathIdentifier
	Tag        string
	ParentPath PathIdentifier
	TypeInfo   string
}

func newAttributeData(name string) *AttributeData {
	return &AttributeData{
		Name: name,
		// Path: path,
	}
}

// A string or series of identifiers that uniquely locate an element within a nested structure, often used to pinpoint where a difference occurs
type PathIdentifier string

type AttributeMap map[PathIdentifier]AttributeData

func newAttributeMap() *AttributeMap {
	return &AttributeMap{}
}

// type optsAttributeData func(c *AttributeData)

// func withTag(v string) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.Tag = v
// 	}
// }

// func applyOptsAttr(c *AttributeData, opts ...optsAttributeData) *AttributeData {
// 	for i := 0; i < len(opts); i++ {
// 		opts[i](c)
// 	}
// 	return c
// }

// func optPath(path PathIdentifier) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.Path = path
// 	}
// }

// func optTypeInfo(typeInfo string) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.TypeInfo = typeInfo
// 	}
// }

// func optAttrParent(parentPath PathIdentifier) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.ParentPath = parentPath
// 	}
// }

// func optTag(tag string) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.Tag = tag
// 	}
// }

// func optName(name string) optsAttributeData {
// 	return func(c *AttributeData) {
// 		c.Name = name
// 	}
// }

// func (m *AttributeMap) Add(parentPath PathIdentifier, value reflect.Value, opts ...optsAttributeData) *AttributeData {
// 	// by default, type name, can be override by options
// 	attr := *applyOptsAttr(newAttributeData(value.Type().Name()), opts...)
// 	realPath := newPath(parentPath, value.Kind(), attr.Name)
// 	(*m)[realPath] = attr
// 	return &attr
// }

// Generate a unique but comprehensible path pattern
func newPath(parent PathIdentifier, kind reflect.Kind, value string) PathIdentifier {
	if parent == "." {
		parent = ""
	}
	switch kind {
	case reflect.Struct:
		return PathIdentifier(fmt.Sprintf("%v.%v", parent, value))
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.String, reflect.Complex64, reflect.Uintptr, reflect.Complex128, reflect.Interface, reflect.UnsafePointer:
		return PathIdentifier(fmt.Sprintf("%v.%v", parent, value))
	case reflect.Slice, reflect.Array:
		return PathIdentifier(fmt.Sprintf("%v.[%v]", parent, value))
	case reflect.Func:
		return PathIdentifier(fmt.Sprintf("%v.(%v)", parent, value))
	case reflect.Chan:
		return PathIdentifier(fmt.Sprintf("%v.chan(%v)", parent, value))
	case reflect.Map:
		return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value))
	case reflect.Invalid:
		return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value))
	default:
		return PathIdentifier(fmt.Sprintf("%v.?{%v}", parent, value))
	}
}

func processBuffer(s chan error, c context.Context) {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(c); err != nil {
		s <- err
		close(s)
		return
	}
	finished := false
	for !finished {
		fn := exe.Pop()
		if err = fn(); err != nil {
			s <- err
			return
		}
		// put an end to it
		if exe.Len() == 0 {
			finished = true
		}
	}
	close(s)
}

func deepDifferenceWait(sigA chan error, sigB chan error) error {
	err := errors.New("deep difference: ")
	hasError := false

	for e := range sigA {
		err = errors.Join(err, e)
		hasError = true
	}

	select {
	case <-sigA:
		fmt.Println("closed")
	}

	for e := range sigB {
		err = errors.Join(err, e)
		hasError = true
	}

	select {
	case <-sigB:
		fmt.Println("closed")
	}

	if hasError {
		return err
	}

	return nil
}
