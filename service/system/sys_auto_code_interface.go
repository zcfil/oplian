package system

import (
	"oplian/global"
	"oplian/model/system/response"
)

type Database interface {
	GetDB(businessDB string) (data []response.Db, err error)
	GetTables(businessDB string, dbName string) (data []response.Table, err error)
	GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error)
}

func (autoCodeService *AutoCodeService) Database(businessDB string) Database {

	if businessDB == "" {
		switch global.ZC_CONFIG.System.DbType {
		case "mysql":
			return AutoCodeMysql
		default:
			return AutoCodeMysql
		}
	} else {
		for _, info := range global.ZC_CONFIG.DBList {
			if info.AliasName == businessDB {

				switch info.Type {
				case "mysql":
					return AutoCodeMysql
				default:
					return AutoCodeMysql
				}
			}
		}
		return AutoCodeMysql
	}

}
