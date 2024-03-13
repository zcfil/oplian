package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/model/system/request"
	"os"
	"sort"
)

const (
	Mysql           = "mysql"
	InitSuccess     = "\n[%v] --> 初始数据成功!\n"
	InitDataExist   = "\n[%v] --> %v 的初始数据已存在!\n"
	InitDataFailed  = "\n[%v] --> %v 初始数据失败! \nerr: %+v\n"
	InitDataSuccess = "\n[%v] --> %v 初始数据成功!\n"
)

const (
	InitOrderSystem   = 10
	InitOrderInternal = 1000
	InitOrderExternal = 100000
)

var (
	ErrMissingDBContext        = errors.New("missing db in context")
	ErrMissingDependentContext = errors.New("missing dependent value in context")
	ErrDBTypeMismatch          = errors.New("db type mismatch")
)

type SubInitializer interface {
	InitializerName() string //
	MigrateTable(ctx context.Context) (next context.Context, err error)
	InitializeData(ctx context.Context) (next context.Context, err error)
	TableCreated(ctx context.Context) bool
	DataInserted(ctx context.Context) bool
}

// TypedDBInitHandler Execute the passed initializer
type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error)
	WriteConfig(ctx context.Context) error
	InitTables(ctx context.Context, inits initSlice) error
	InitData(ctx context.Context, inits initSlice) error
}

// orderedInitializer
type orderedInitializer struct {
	order int
	SubInitializer
}

// initSlice
type initSlice []*orderedInitializer

var (
	initializers initSlice
	cache        map[string]*orderedInitializer
)

// RegisterInit Register the initialization procedure to be performed, which is called at InitDB()
func RegisterInit(order int, i SubInitializer) {

	//log.Println("generate config file")
	isNew, err := OpGenerateConfigFile(define.PathIpfsConfig, ConfigList)
	if err != nil {
		log.Fatal("Failed to generate the corresponding config file")
	} else {
		if isNew {
			log.Fatal("The config file is generated. Configure the config file correctly and start the startup again")
		}
	}

	if initializers == nil {
		initializers = initSlice{}
	}
	if cache == nil {
		cache = map[string]*orderedInitializer{}
	}
	name := i.InitializerName()
	if _, existed := cache[name]; existed {
		panic(fmt.Sprintf("Name conflict on %s", name))
	}
	ni := orderedInitializer{order, i}
	//log.Println("Register：", i.InitializerName())
	initializers = append(initializers, &ni)
	cache[name] = &ni
}

/* ---- * service * ---- */

type InitDBService struct{}

// InitDB Create the database and initialize the general entry
func (initDBService *InitDBService) InitDB(conf request.InitDB) (err error) {
	ctx := context.TODO()
	if len(initializers) == 0 {
		return errors.New("No initialization procedure is available. Please check whether the initialization is complete")
	}
	sort.Sort(&initializers)

	var initHandler TypedDBInitHandler
	switch conf.DBType {
	case "mysql":
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	default:
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	}
	ctx, err = initHandler.EnsureDB(ctx, &conf)
	if err != nil {
		return err
	}

	db := ctx.Value("db").(*gorm.DB)
	global.ZC_DB = db

	if err = initHandler.InitTables(ctx, initializers); err != nil {
		return err
	}

	if err = initHandler.InitData(ctx, initializers); err != nil {
		return err
	}

	if err = initHandler.WriteConfig(ctx); err != nil {
		return err
	}
	initializers = initSlice{}
	cache = map[string]*orderedInitializer{}
	return nil
}

// InitDB Create the database and initialize the general entry
func (initDBService *InitDBService) InitData(DBType string) (err error) {
	ctx := context.TODO()
	if len(initializers) == 0 {
		return errors.New("no initialization procedure is available. Please check whether the initialization is complete")
	}
	sort.Sort(&initializers)

	var initHandler TypedDBInitHandler
	switch DBType {
	case "mysql":
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	default:
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	}
	ctx = context.WithValue(ctx, "db", global.ZC_DB)
	if err = initHandler.InitData(ctx, initializers); err != nil {
		return err
	}

	initializers = initSlice{}
	return nil
}

// createDatabase Create database (called in EnsureDB())
func createDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// createTables Create tables (default dbInitHandler.initTables behavior)
func createTables(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer func(c func()) { c() }(cancel)
	for _, init := range inits {
		if init.TableCreated(next) {
			continue
		}
		if n, err := init.MigrateTable(next); err != nil {
			return err
		} else {
			next = n
		}

	}
	return nil
}

/* -- sortable interface -- */

func (a initSlice) Len() int {
	return len(a)
}

func (a initSlice) Less(i, j int) bool {
	return a[i].order < a[j].order
}

func (a initSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ScriptFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

var ConfigList = []ScriptFile{
	{
		Name:    "config.yaml",
		Content: "autocode:\n  transfer-restart: true\n  root: G:\\project\\zc-admin\n  server: /server\n  server-api: /api/v1/%s\n  server-plug: /plugin/%s\n  server-initialize: /initialize\n  server-model: /model/%s\n  server-request: /model/%s/request/\n  server-router: /router/%s\n  server-service: /service/%s\n  web: /web/src\n  web-api: /api\n  web-form: /view\n  web-table: /view\ncaptcha:\n  key-long: 6\n  img-width: 240\n  img-height: 80\ncors:\n  mode: whitelist\n  whitelist:\n  - allow-origin: example1.com\n    allow-methods: GET, POST\n    allow-headers: content-type\n    expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,\n      Content-Type\n    allow-credentials: true\n  - allow-origin: example2.com\n    allow-methods: GET, POST\n    allow-headers: content-type\n    expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,\n      Content-Type\n    allow-credentials: true\njwt:\n  signing-key: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\n  expires-time: 7d\n  buffer-time: 1d\n  issuer: qmPlus\nlocal:\n  path: uploads/file\n  store-path: uploads/file\nmysql:\n  path: 8.129.83.148\n  port: \"3306\"\n  config: charset=utf8mb4&parseTime=True&loc=Local\n  db-name: zc_test\n  username: root\n  password: ZCXTong@2023!+\n  max-idle-conns: 10\n  max-open-conns: 100\n  log-mode: error\n  log-zap: false\nsystem:\n  env: public\n  addr: 50005\n  db-type: mysql\n  oss-type: local\n  use-multipoint: false\n  iplimit-count: 15000\n  iplimit-time: 3600\ntimer:\n  start: true\n  spec: '@daily'\n  with_seconds: false\n  detail:\n  - tableName: sys_operation_records\n    compareField: created_at\n    interval: 2160h\n  - tableName: jwt_blacklists\n    compareField: created_at\n    interval: 168h",
	},
	{
		Name:    "config_room.yaml",
		Content: "web:\n  addr: 127.0.0.1:50005\n  token: slfsdaklfhasldfjda\ngateway:\n  gateWayId: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\n  port: 50006\n  ip: 127.0.0.1\n  token: slfsdaklfhasldfjda\nop:\n  ip: 127.0.0.1\n  port: 50007\n  opId: 8b356e3d-5e35-47a6-a933-07bbb8cc0bbb1\n  token: sdhskfsdhfkjashflksahfklashfsadd\nopc-c2:\n  port: 50008\n  opId: 8b356e3d-5e35-47a6-a933-07bbb8cc0bbb1\n  token: sdhskfsdhfkjashflksahfklashfsadd1\njwt:\n  signing-key: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb\ndownload:\n  addr: 14.215.165.39:9533\n  user: admin\n  password: iZw03PNfUWu3IIHqpS\nzap:\n  level: info\n  prefix: '[oplian]'\n  format: console\n  director: log\n  encode-level: LowercaseColorLevelEncoder\n  stacktrace-key: stacktrace\n  max-age: 0\n  show-line: true\n  log-in-console: true\n",
	},
}

// OpGenerateConfigFile The op terminal initializes and generates the config file
func OpGenerateConfigFile(filePath string, fileList []ScriptFile) (bool, error) {

	os.MkdirAll(define.PathIpfsProgram, os.ModePerm)
	os.MkdirAll(define.PathIpfsLog, os.ModePerm)
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Println("function os.MkdirAll() Filed", err.Error())
		return false, err
	}

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
