package config

type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`
	StorePath string `mapstructure:"store-path" json:"store-path" yaml:"store-path"`
}
