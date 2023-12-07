package hikaku

import "reflect"

type config struct {
	tag string
}

func newConfig() *config {
	return &config{
		tag: "json",
	}
}

func applyOpts(c *config, opts ...optsConfig) *config {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	return c
}

type optsConfig func(c *config)

func WithTag(tagName string) optsConfig {
	return func(c *config) {
		c.tag = tagName
	}
}

type valueOptions struct {
	isPointer  bool
	parent     PathIdentifier
	path       string
	tag        reflect.StructTag
	parentType reflect.Type
	fieldIndex int
	fieldName  string
	typeInfo   string
}

func newValueOptions() *valueOptions {
	return &valueOptions{
		fieldIndex: -1,
	}
}

type optsValueOptions func(c *valueOptions)

func fromPointer() optsValueOptions {
	return func(c *valueOptions) {
		c.isPointer = true
	}
}

func fromTypeInfo(typeInfo string) optsValueOptions {
	return func(c *valueOptions) {
		c.typeInfo = typeInfo
	}
}

func fromParentPath(path PathIdentifier) optsValueOptions {
	return func(c *valueOptions) {
		c.parent = path
	}
}

func fromParentType(parentType reflect.Type) optsValueOptions {
	return func(c *valueOptions) {
		c.parentType = parentType
	}
}

func fromTag(tag reflect.StructTag) optsValueOptions {
	return func(c *valueOptions) {
		c.tag = tag
	}
}

func fromPath(path string) optsValueOptions {
	return func(c *valueOptions) {
		c.path = path
	}
}

func fromFieldName(name string) optsValueOptions {
	return func(c *valueOptions) {
		c.fieldName = name
	}
}

func fromFieldIndex(index int) optsValueOptions {
	return func(c *valueOptions) {
		c.fieldIndex = index
	}
}

func applyValueOptions(c *valueOptions, opts ...optsValueOptions) *valueOptions {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	return c
}
