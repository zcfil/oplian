package initialize

import (
	"log"
	"os"
)

var ConfigList = []ScriptFile{
	{
		Name:    "config.yaml",
		Content: "cors:\n  mode: whitelist\n  whitelist:\n  - allow-origin: example1.com\n    allow-methods: GET, POST\n    allow-headers: content-type\n    expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,\n      Content-Type\n    allow-credentials: true\n  - allow-origin: example2.com\n    allow-methods: GET, POST\n    allow-headers: content-type\n    expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,\n      Content-Type\n    allow-credentials: true\njwt:\n  signing-key: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\n  expires-time: 7d\n  buffer-time: 1d\n  issuer: qmPlus\nlocal:\n  path: uploads/file\n  store-path: uploads/file\nmysql:\n  path: 8.129.83.148\n  port: \"3306\"\n  config: charset=utf8mb4&parseTime=True&loc=Local\n  db-name: zc_test\n  username: root\n  password: ZCXTong@2023!+\n  max-idle-conns: 10\n  max-open-conns: 100\n  log-mode: error\n  log-zap: false\nredis:\n  db: 0\n  addr: 127.0.0.1:6379\n  password: \"\"\nsystem:\n  env: public\n  addr: 50005\n  db-type: mysql\n  oss-type: local\n  use-multipoint: false\n  use-redis: false\n  iplimit-count: 15000\n  iplimit-time: 3600\ntimer:\n  start: true\n  spec: '@daily'\n  with_seconds: false\n  detail:\n  - tableName: sys_operation_records\n    compareField: created_at\n    interval: 2160h\n  - tableName: jwt_blacklists\n    compareField: created_at\n    interval: 168h",
	},
	{
		Name:    "config_room.yaml",
		Content: "web:\n  addr: 127.0.0.1:50005\n  token: slfsdaklfhasldfjda\ngateway:\n  gateWayId: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\n  port: 50006\n  ip: 127.0.0.1\n  token: slfsdaklfhasldfjda\nop:\n  ip: 127.0.0.1\n  port: 50007\n  opId: 8b356e3d-5e35-47a6-a933-07bbb8cc0bbb1\n  token: sdhskfsdhfkjashflksahfklashfsadd\nopc-c2:\n  port: 50008\n  opId: 8b356e3d-5e35-47a6-a933-07bbb8cc0bbb1\n  token: sdhskfsdhfkjashflksahfklashfsadd1\njwt:\n  signing-key: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\n  expires-time: 7d\n  buffer-time: 1d\n  issuer: qmPlus\ndownload:\n  addr: 14.215.165.39:9533\n  user: admin\n  password: iZw03PNfUWu3IIHqpS\nzap:\n  level: info\n  prefix: '[oplian]'\n  format: console\n  director: log\n  encode-level: LowercaseColorLevelEncoder\n  stacktrace-key: stacktrace\n  max-age: 0\n  show-line: true\n  log-in-console: true\n",
	},
}

// OpGenerateConfigFile The op terminal initializes and generates the config file
func OpGenerateConfigFile(filePath string, fileList []ScriptFile) (bool, error) {

	_ = os.MkdirAll(filePath, os.ModePerm)

	var isNew bool

	for _, val := range fileList {
		path := filePath + val.Name
		if _, err := os.Stat(path); os.IsNotExist(err) {
			isNew = true

			out, createErr := os.Create(path)
			if createErr != nil {
				log.Println("function os.Create() Filed", createErr.Error())
				return false, createErr
			}
			defer out.Close() //

			_, writeErr := out.WriteString(val.Content)
			if writeErr != nil {
				log.Println("function os.WriteString() Filed", writeErr.Error())
				return false, createErr
			}
		}
	}
	return isNew, nil
}
