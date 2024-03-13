package config

type Timer struct {
	Start       bool     `mapstructure:"start" json:"start" yaml:"start"`
	Spec        string   `mapstructure:"spec" json:"spec" yaml:"spec"`
	WithSeconds bool     `mapstructure:"with_seconds" json:"with_seconds" yaml:"with_seconds"`
	Detail      []Detail `mapstructure:"detail" json:"detail" yaml:"detail"`
}

type Detail struct {
	TableName    string `mapstructure:"tableName" json:"tableName" yaml:"tableName"`
	CompareField string `mapstructure:"compareField" json:"compareField" yaml:"compareField"`
	Interval     string `mapstructure:"interval" json:"interval" yaml:"interval"`
}
