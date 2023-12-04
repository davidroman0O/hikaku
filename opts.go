package hikaku

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
	isPointer bool
}

func newValueOptions() *valueOptions {
	return &valueOptions{}
}

type optsValueOptions func(c *valueOptions)

func fromPointer() optsValueOptions {
	return func(c *valueOptions) {
		c.isPointer = true
	}
}

func applyValueOptions(c *valueOptions, opts ...optsValueOptions) *valueOptions {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	return c
}
