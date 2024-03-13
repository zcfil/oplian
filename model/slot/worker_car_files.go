package slot

import (
	"oplian/global"
)

type WorkerCarFiles struct {
	global.ZC_MODEL
	RelationId  int    `gorm:"comment:关联ID" json:"relationId"`             // 关联ID
	FileName    string `gorm:"type:longtext;comment:文件名称" json:"fileName"` // 文件名称
	FileIndex   int    `gorm:"comment:文件索引" json:"fileIndex"`              // 文件索引
	FileStr     string `gorm:"type:longtext;comment:生成参数" json:"fileStr"`  // 生成参数
	CarFileName string `gorm:"index;comment:Car名称" json:"carFileName"`     // 生成参数
	PieceCid    string `gorm:"index;comment:pieceCid" json:"pieceCid"`     // pieceCid
	PieceSize   int    `gorm:"comment:pieceSize" json:"pieceSize"`         // pieceSize
	CarSize     int    `gorm:"comment:carSize" json:"carSize"`             // carSize
	DataCid     string `gorm:"comment:DataCid" json:"dataCid"`             // dataCid
	InputDir    string `gorm:"comment:输入目录" json:"inputDir"`               // 输入目录
}

func (WorkerCarFiles) TableName() string {
	return "worker_car_files"
}
