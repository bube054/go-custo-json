package gocustojson

type Config struct {
	AllowExtraWS bool

	AllowHexNumbers       bool
	AllowPointEdgeNumbers bool
	AllowInfinity         bool
	AllowNaN              bool

	AllowUnquoted         bool
	AllowSingleQuotes     bool
	AllowNewlineInStrings bool
	AllowEscapeChars      bool

	AllowTrailingCommaArray  bool
	AllowTrailingCommaObject bool

	AllowLineComments  bool
	AllowBlockComments bool

	AllowMalformedInput bool
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

func WithAllowInfinity(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowInfinity = allow
	}
}

func WithAllowNaN(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNaN = allow
	}
}

func WithAllowUnquoted(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowUnquoted = allow
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

func withMalformedInput(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowMalformedInput = allow
	}
}
