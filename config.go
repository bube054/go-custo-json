package gocustojson

type Config struct {
	AllowExtraWS bool

	AllowHexNumbers       bool
	AllowPointEdgeNumbers bool
	AllowPosInfinity      bool
	AllowNegInfinity      bool
	AllowNaN              bool

	AllowSingleQuotes     bool
	AllowNewlineInStrings bool
	AllowEscapeChars      bool

	AllowTrailingCommaArray  bool
	AllowTrailingCommaObject bool

	AllowLineComments  bool
	AllowBlockComments bool
}

func NewConfig(opts ...func(*Config)) *Config {
	cfg := &Config{}

	for _, o := range opts {
		o(cfg)
	}

	return cfg
}

func WithAllowExtraWS(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowExtraWS = allow
	}
}

func WithAllowHexNumbers(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowHexNumbers = allow
	}
}

func WithAllowPointEdgeNumbers(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowPointEdgeNumbers = allow
	}
}

func WithAllowPosInfinity(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowPosInfinity = allow
	}
}

func WithAllowNegInfinity(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNegInfinity = allow
	}
}

func WithAllowNaN(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNaN = allow
	}
}

func WithAllowSingleQuotes(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowSingleQuotes = allow
	}
}

func WithAllowNewlineInStrings(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNewlineInStrings = allow
	}
}

func WithAllowEscapeChars(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowEscapeChars = allow
	}
}

func WithAllowTrailingCommaArray(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowTrailingCommaArray = allow
	}
}

func WithAllowTrailingCommaObject(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowTrailingCommaObject = allow
	}
}

func WithAllowLineComments(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowLineComments = allow
	}
}

func WithAllowBlockComments(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowBlockComments = allow
	}
}
