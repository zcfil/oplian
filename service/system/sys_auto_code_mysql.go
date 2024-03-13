package system

import (
	"oplian/global"
	"oplian/model/system/response"
)

var AutoCodeMysql = new(autoCodeMysql)

type autoCodeMysql struct{}

// GetDB Gets all database names for the database

func (s *autoCodeMysql) GetDB(businessDB string) (data []response.Db, err error) {
	var entities []response.Db
	sql := "SELECT SCHEMA_NAME AS `database` FROM INFORMATION_SCHEMA.SCHEMATA;"
	if businessDB == "" {
		err = global.ZC_DB.Raw(sql).Scan(&entities).Error
	} else {
		err = global.ZC_DBList[businessDB].Raw(sql).Scan(&entities).Error
	}
	return entities, err
}

// GetTables Gets all the table names for the database

func (s *autoCodeMysql) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	var entities []response.Table
	sql := `select table_name as table_name from information_schema.tables where table_schema = ?`
	if businessDB == "" {
		err = global.ZC_DB.Raw(sql, dbName).Scan(&entities).Error
	} else {
		err = global.ZC_DBList[businessDB].Raw(sql, dbName).Scan(&entities).Error
	}

	return entities, err
}

// GetColumn Gets all field names, type values, and so on for the specified database and specified data table

func (s *autoCodeMysql) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	var entities []response.Column
	sql := `
	SELECT COLUMN_NAME        column_name,
       DATA_TYPE          data_type,
       CASE DATA_TYPE
           WHEN 'longtext' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'varchar' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'double' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'decimal' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'int' THEN c.NUMERIC_PRECISION
           WHEN 'bigint' THEN c.NUMERIC_PRECISION
           ELSE '' END AS data_type_long,
       COLUMN_COMMENT     column_comment
	FROM INFORMATION_SCHEMA.COLUMNS c
	WHERE table_name = ?
	  AND table_schema = ?
	`
	if businessDB == "" {
		err = global.ZC_DB.Raw(sql, tableName, dbName).Scan(&entities).Error
	} else {
		err = global.ZC_DBList[businessDB].Raw(sql, tableName, dbName).Scan(&entities).Error
	}

	return entities, err
}
