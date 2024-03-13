package initialize

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"oplian/global"
	"oplian/initialize/internal"
)

// GormMysql Initialize the Mysql database

func GormMysql() *gorm.DB {
	m := global.ZC_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
