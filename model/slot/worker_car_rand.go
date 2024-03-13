package slot

import (
	"oplian/global"
)

type WorkerCarRand struct {
	global.ZC_MODEL
	NumIndex int `gorm:"comment:随机数索引" json:"numIndex"` // 随机数索引
	Number   int `gorm:"comment:随机数" json:"number"`     // 随机数

}

func (WorkerCarRand) TableName() string {
	return "worker_car_rand"
}
