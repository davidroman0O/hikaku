package hikaku

import (
	"fmt"
	"reflect"
)

// `ProbeMap` contains a flat map of all properties which has been identified recurssively with `Probe` data
type ProbeMap map[PathIdentifier]Probe

var keyProbeMapCtx string = "mapProbes"

func newProbeMap() *ProbeMap {
	return &ProbeMap{}
}

func (p *ProbeMap) Add(probe *Probe) {
	(*p)[probe.path] = *probe
}

// `Probe` represent the instance of a data at a specific path at a specific level of a type and extract it's info
// It is accumulating information at runtime which can be costly when doing deep difference of nested structures
// One `Probe` is supposed to be part of a tree of values by knowing it's parent
// TODO @droman: what about imported types?
// TODO @droman: what about non-imported types?
type Probe struct {
	level      int
	path       PathIdentifier
	parentPath PathIdentifier
	// A parent `Probe` will allow us to know all the kind of different types of a datastructure.
	// Which is useful for complex and unusual types.
	parent     *Probe
	isPointer  bool
	kind       reflect.Kind
	tag        reflect.StructTag
	parentType reflect.Type
	value      reflect.Value
	data       interface{}
	fieldIndex int
	fieldName  string
	typeName   string
}

func newProbe() *Probe {
	return &Probe{
		fieldIndex: -1,
	}
}

type optionProbe func(c *Probe)

func probeWithProbe(probe *Probe) optionProbe {
	return func(c *Probe) {
		(*c) = *probe
	}
}

func probeWithParentProbe(probe *Probe) optionProbe {
	return func(c *Probe) {
		c.parent = probe
	}
}

func probeWithPointer() optionProbe {
	return func(c *Probe) {
		c.isPointer = true
	}
}

func probeWithData(data interface{}) optionProbe {
	return func(c *Probe) {
		c.data = data
	}
}

func probeWithValue(value reflect.Value) optionProbe {
	return func(c *Probe) {
		c.value = value
	}
}

func probeWithLevel(level int) optionProbe {
	return func(c *Probe) {
		c.level = level
	}
}

func probeWithTypeName(typeName string) optionProbe {
	return func(c *Probe) {
		c.typeName = typeName
	}
}

func probeWithKind(kind reflect.Kind) optionProbe {
	return func(c *Probe) {
		c.kind = kind
	}
}

func probeWithParentPath(path PathIdentifier) optionProbe {
	return func(c *Probe) {
		c.parentPath = path
	}
}

func probeWithParentType(parentType reflect.Type) optionProbe {
	return func(c *Probe) {
		c.parentType = parentType
	}
}

func probeWithTag(tag reflect.StructTag) optionProbe {
	return func(c *Probe) {
		c.tag = tag
	}
}

func probeWithPath(path PathIdentifier) optionProbe {
	return func(c *Probe) {
		c.path = path
	}
}

func probeWithFieldName(name string) optionProbe {
	return func(c *Probe) {
		c.fieldName = name
	}
}

func probeWithFieldIndex(index int) optionProbe {
	return func(c *Probe) {
		c.fieldIndex = index
	}
}

func applyProbeOptions(c *Probe, opts ...optionProbe) *Probe {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	// after apply every options we should compute the path
	c.path = computePath(c)

	// pp.Println(c)
	return c
}

func computePath(c *Probe) PathIdentifier {
	if c.kind.String() == "" {
		return "unknown"
	}

	switch c.kind {
	case reflect.String:
		switch c.level {
		case 1:
			return PathIdentifier(fmt.Sprintf(".%v", c.fieldName))
		default:
			return PathIdentifier(fmt.Sprintf("%v.%v", c.parentPath, c.fieldName))
		}
	case reflect.Struct:
		switch c.level {
		case 0:
			return PathIdentifier(fmt.Sprintf("%v", c.fieldName))
		default:
			return PathIdentifier(fmt.Sprintf("%v.%v", c.parentPath, c.fieldName))
		}

	}

	return "unknown"
}

// if parent == "." {
// 	parent = ""
// }
// switch kind {
// case reflect.Struct:
// 	return PathIdentifier(fmt.Sprintf("%v.%v", parent, value))
// case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool, reflect.String, reflect.Complex64, reflect.Uintptr, reflect.Complex128, reflect.Interface, reflect.UnsafePointer:
// 	return PathIdentifier(fmt.Sprintf("%v.%v", parent, value))
// case reflect.Slice, reflect.Array:
// 	return PathIdentifier(fmt.Sprintf("%v.[%v]", parent, value))
// case reflect.Func:
// 	return PathIdentifier(fmt.Sprintf("%v.(%v)", parent, value))
// case reflect.Chan:
// 	return PathIdentifier(fmt.Sprintf("%v.chan(%v)", parent, value))
// case reflect.Map:
// 	return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value))
// case reflect.Invalid:
// 	return PathIdentifier(fmt.Sprintf("%v.map(%v)", parent, value))
// default:
// 	return PathIdentifier(fmt.Sprintf("%v.?{%v}", parent, value))
// }
