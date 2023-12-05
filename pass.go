package hikaku

import (
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
}

func newAttributeData(name string, path PathIdentifier) *AttributeData {
	return &AttributeData{
		Name: name,
		Path: path,
	}
}

// A string or series of identifiers that uniquely locate an element within a nested structure, often used to pinpoint where a difference occurs
type PathIdentifier string

type AttributeMap map[PathIdentifier]AttributeData

func newAttributeMap() *AttributeMap {
	return &AttributeMap{}
}

type optsAttributeData func(c *AttributeData)

func withTag(v string) optsAttributeData {
	return func(c *AttributeData) {
		c.Tag = v
	}
}

func applyOptsAttr(c *AttributeData, opts ...optsAttributeData) *AttributeData {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	return c
}

func (m *AttributeMap) Add(path PathIdentifier, value reflect.Value, opts ...optsAttributeData) *AttributeData {
	attr := *applyOptsAttr(newAttributeData(value.Type().Name(), newPath(path, value)), opts...)
	(*m)[path] = attr
	return &attr
}

func newPath(parent PathIdentifier, value reflect.Value) PathIdentifier {
	switch value.Kind() {
	// we need to loop through all structfields
	case reflect.Struct:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Bool:
	case reflect.String:
	case reflect.Complex64:
	case reflect.Uintptr:
	case reflect.Complex128:
	case reflect.Interface:
	case reflect.UnsafePointer:
		return PathIdentifier(fmt.Sprintf("%v.%v", parent, value.Type().Name()))
	case reflect.Slice:
	case reflect.Array:
		return PathIdentifier(fmt.Sprintf("%v.[%v]", parent, value.Type().Name()))
	case reflect.Func:
		return PathIdentifier(fmt.Sprintf("%v.(%v)", parent, value.Type().Name()))
	case reflect.Chan:
		return PathIdentifier(fmt.Sprintf("%v.chan(%v)", parent, value.Type().Name()))
	case reflect.Map:
		return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value.Type().Name()))
	case reflect.Invalid:
		return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value.Type().Name()))
	default:
		return PathIdentifier(fmt.Sprintf("%v.?(%v)", parent, value.Type().Name()))
	}
	return ""
}
