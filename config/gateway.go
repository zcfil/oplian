package config

type Web struct {
	Addr  string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
}

type GateWay struct {
	IP    string `mapstructure:"ip" json:"ip" yaml:"ip"`
	Port  string `mapstructure:"port" json:"port" yaml:"port"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
	//GateWayId string `mapstructure:"gateWayId" json:"gateWayId" yaml:"gate-way-id"`
}

type Op struct {
	IP    string `mapstructure:"ip" json:"ip" yaml:"ip"`
	Port  string `mapstructure:"port" json:"host" yaml:"port"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
	OpId  string `mapstructure:"opId" json:"opId" yaml:"op-id"`
}

type OpC2 struct {
	Port  string `mapstructure:"port" json:"host" yaml:"port"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
	OpId  string `mapstructure:"opId" json:"opId" yaml:"op-id"`
}

type Download struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
