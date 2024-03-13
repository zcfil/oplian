package config

type Server struct {
	JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	// gorm
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	// oss
	Local Local `mapstructure:"local" json:"local" yaml:"local"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}

type ServerRoom struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Web      Web      `mapstructure:"web" json:"web" yaml:"web"`
	Gateway  GateWay  `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	Op       Op       `mapstructure:"op" json:"op" yaml:"op"`
	OpC2     OpC2     `mapstructure:"opc-c2" json:"opc-c2" yaml:"opc-c2"`
	Download Download `mapstructure:"download" json:"download" yaml:"download"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
}
