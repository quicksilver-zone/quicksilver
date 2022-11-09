package types

type Config struct {
	SourceChain string            `yaml:"source_chain"`
	SourceLcd   string            `yaml:"source_lcd"`
	Chains      map[string]string `yaml:"chains"`
}
