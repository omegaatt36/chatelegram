package repository

// ClientOption defines how to inject option.
type ClientOption interface {
	injectOption(*common)
}

var _ ClientOption = &WithMaxToken{}

// WithMaxToken helps inject max token.
type WithMaxToken struct {
	MaxToken int
}

func (o WithMaxToken) injectOption(c *common) {
	c.maxToken = o.MaxToken
}

var _ ClientOption = &WithCompletionsEngine{}

// WithCompletionsEngine helps inject completions engine.
type WithCompletionsEngine struct {
	Engine string
}

func (o WithCompletionsEngine) injectOption(c *common) {
	c.engine = o.Engine
}
