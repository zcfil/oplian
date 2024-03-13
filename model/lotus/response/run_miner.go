package response

type RunMiner struct {
	Ip           string `json:"ip" gorm:"comment:ip"`
	OpId         string `json:"opId" gorm:"comment:设备uuid"`
	GateId       string `json:"gateId" gorm:"comment:机房uuid"`
	Actor        string `json:"actor" gorm:"comment:节点号"`
	Port         string `json:"port" gorm:"comment:端口号"`
	DeployStatus int    `json:"deployStatus" gorm:"comment:部署状态"`
	RunStatus    int    `json:"runStatus" gorm:"comment:运行状态"`
	IsManage     bool   `json:"isManage" gorm:"comment:是否打开调度"`
	IsWdpost     bool   `json:"isWdpost" gorm:"comment:是否打开wdpost"`
	IsWnpost     bool   `json:"isWnpost" gorm:"comment:是否打开wnpost"`
	Partitions   string `json:"partitions" gorm:"comment:wdpost partitions号"`
	LotusToken   string `json:"lotusToken"  gorm:"comment:lotus token"`
	LotusIp      string `json:"lotusIp"  gorm:"comment:lotus ip"`
}
