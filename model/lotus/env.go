package lotus

import "oplian/global"

type LotusEnv struct {
	global.ZC_MODEL
	Ekey    string        `json:"ekey"  gorm:"index;comment:环境变量key"`
	Evalue  string        `json:"evalue"  gorm:"comment:环境变量value"`
	MinerId string        `json:"minerId"  gorm:"comment:节点ID"`
	Etype   EnvConfigType `json:"Etype"  gorm:"comment:配置类型"`
	Remark  string        `json:"remark"  gorm:"comment:remark"`
}

func (LotusEnv) TableName() string {
	return "lotus_env"
}
