package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"oplian/global"
	sysModel "oplian/model/system"
	"oplian/service/system"
	"time"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ZC_MODEL: global.ZC_MODEL{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "base", Method: "POST", Path: "/base/login", Description: "用户登录(必选)"},

		{ZC_MODEL: global.ZC_MODEL{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ZC_MODEL: global.ZC_MODEL{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ZC_MODEL: global.ZC_MODEL{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ZC_MODEL: global.ZC_MODEL{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ZC_MODEL: global.ZC_MODEL{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 9, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ZC_MODEL: global.ZC_MODEL{ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},

		{ZC_MODEL: global.ZC_MODEL{ID: 13, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "创建api"},
		{ZC_MODEL: global.ZC_MODEL{ID: 14, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "删除Api"},
		{ZC_MODEL: global.ZC_MODEL{ID: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "更新Api"},
		{ZC_MODEL: global.ZC_MODEL{ID: 16, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "获取api列表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 17, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},
		{ZC_MODEL: global.ZC_MODEL{ID: 18, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "获取api详细信息"},
		{ZC_MODEL: global.ZC_MODEL{ID: 19, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "批量删除api"},

		{ZC_MODEL: global.ZC_MODEL{ID: 20, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ZC_MODEL: global.ZC_MODEL{ID: 21, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ZC_MODEL: global.ZC_MODEL{ID: 22, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ZC_MODEL: global.ZC_MODEL{ID: 23, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ZC_MODEL: global.ZC_MODEL{ID: 24, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 25, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "角色", Method: "POST", Path: "/authority/setDataAuthority", Description: "设置角色资源权限"},

		{ZC_MODEL: global.ZC_MODEL{ID: 26, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ZC_MODEL: global.ZC_MODEL{ID: 27, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ZC_MODEL: global.ZC_MODEL{ID: 28, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/addBaseMenu", Description: "新增菜单"},
		{ZC_MODEL: global.ZC_MODEL{ID: 29, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenu", Description: "获取菜单树(必选)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 30, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/deleteBaseMenu", Description: "删除菜单"},
		{ZC_MODEL: global.ZC_MODEL{ID: 31, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/updateBaseMenu", Description: "更新菜单"},
		{ZC_MODEL: global.ZC_MODEL{ID: 32, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuById", Description: "根据id获取菜单"},
		{ZC_MODEL: global.ZC_MODEL{ID: 33, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuList", Description: "分页获取基础menu列表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 34, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "获取用户动态路由"},
		{ZC_MODEL: global.ZC_MODEL{ID: 35, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuAuthority", Description: "获取指定角色menu"},
		{ZC_MODEL: global.ZC_MODEL{ID: 36, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "菜单", Method: "POST", Path: "/menu/addMenuAuthority", Description: "增加menu和角色关联关系"},

		{ZC_MODEL: global.ZC_MODEL{ID: 37, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "分片上传", Method: "GET", Path: "/fileUploadAndDownload/findFile", Description: "寻找目标文件（秒传）"},
		{ZC_MODEL: global.ZC_MODEL{ID: 38, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinue", Description: "断点续传"},
		{ZC_MODEL: global.ZC_MODEL{ID: 39, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinueFinish", Description: "断点续传完成"},
		{ZC_MODEL: global.ZC_MODEL{ID: 40, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "分片上传", Method: "POST", Path: "/fileUploadAndDownload/removeChunk", Description: "上传完成移除文件"},

		{ZC_MODEL: global.ZC_MODEL{ID: 41, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "文件上传示例"},
		{ZC_MODEL: global.ZC_MODEL{ID: 42, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "删除文件"},
		{ZC_MODEL: global.ZC_MODEL{ID: 43, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "文件名或者备注编辑"},
		{ZC_MODEL: global.ZC_MODEL{ID: 44, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "文件上传与下载", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "获取上传文件列表"},

		{ZC_MODEL: global.ZC_MODEL{ID: 45, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统服务", Method: "POST", Path: "/system/getServerInfo", Description: "获取服务器信息"},
		{ZC_MODEL: global.ZC_MODEL{ID: 46, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统服务", Method: "POST", Path: "/system/getSystemConfig", Description: "获取配置文件内容"},
		{ZC_MODEL: global.ZC_MODEL{ID: 47, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统服务", Method: "POST", Path: "/system/setSystemConfig", Description: "设置配置文件内容"},

		{ZC_MODEL: global.ZC_MODEL{ID: 48, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "客户", Method: "PUT", Path: "/customer/customer", Description: "更新客户"},
		{ZC_MODEL: global.ZC_MODEL{ID: 49, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "客户", Method: "POST", Path: "/customer/customer", Description: "创建客户"},
		{ZC_MODEL: global.ZC_MODEL{ID: 50, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "客户", Method: "DELETE", Path: "/customer/customer", Description: "删除客户"},
		{ZC_MODEL: global.ZC_MODEL{ID: 51, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "客户", Method: "GET", Path: "/customer/customer", Description: "获取单一客户"},
		{ZC_MODEL: global.ZC_MODEL{ID: 52, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "客户", Method: "GET", Path: "/customer/customerList", Description: "获取客户列表"},

		{ZC_MODEL: global.ZC_MODEL{ID: 53, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getDB", Description: "获取所有数据库"},
		{ZC_MODEL: global.ZC_MODEL{ID: 54, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getTables", Description: "获取数据库表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 55, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/createTemp", Description: "自动化代码"},
		{ZC_MODEL: global.ZC_MODEL{ID: 56, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/preview", Description: "预览自动化代码"},
		{ZC_MODEL: global.ZC_MODEL{ID: 57, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "GET", Path: "/autoCode/getColumn", Description: "获取所选table的所有字段"},
		{ZC_MODEL: global.ZC_MODEL{ID: 58, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/createPlug", Description: "自动创建插件包"},
		{ZC_MODEL: global.ZC_MODEL{ID: 59, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器", Method: "POST", Path: "/autoCode/installPlugin", Description: "安装插件"},

		{ZC_MODEL: global.ZC_MODEL{ID: 60, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "包（pkg）生成器", Method: "POST", Path: "/autoCode/createPackage", Description: "生成包(package)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 61, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "包（pkg）生成器", Method: "POST", Path: "/autoCode/getPackage", Description: "获取所有包(package)"},
		{ZC_MODEL: global.ZC_MODEL{ID: 62, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "包（pkg）生成器", Method: "POST", Path: "/autoCode/delPackage", Description: "删除包(package)"},

		{ZC_MODEL: global.ZC_MODEL{ID: 63, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/getMeta", Description: "获取meta信息"},
		{ZC_MODEL: global.ZC_MODEL{ID: 64, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/rollback", Description: "回滚自动生成代码"},
		{ZC_MODEL: global.ZC_MODEL{ID: 65, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/getSysHistory", Description: "查询回滚记录"},
		{ZC_MODEL: global.ZC_MODEL{ID: 66, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "代码生成器历史", Method: "POST", Path: "/autoCode/delSysHistory", Description: "删除回滚记录"},

		{ZC_MODEL: global.ZC_MODEL{ID: 67, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典详情", Method: "PUT", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Description: "更新字典内容"},
		{ZC_MODEL: global.ZC_MODEL{ID: 68, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典详情", Method: "POST", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Description: "新增字典内容"},
		{ZC_MODEL: global.ZC_MODEL{ID: 69, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典详情", Method: "DELETE", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Description: "删除字典内容"},
		{ZC_MODEL: global.ZC_MODEL{ID: 70, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Description: "根据ID获取字典内容"},
		{ZC_MODEL: global.ZC_MODEL{ID: 71, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Description: "获取字典内容列表"},

		{ZC_MODEL: global.ZC_MODEL{ID: 72, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典", Method: "POST", Path: "/sysDictionary/createSysDictionary", Description: "新增字典"},
		{ZC_MODEL: global.ZC_MODEL{ID: 73, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典", Method: "DELETE", Path: "/sysDictionary/deleteSysDictionary", Description: "删除字典"},
		{ZC_MODEL: global.ZC_MODEL{ID: 74, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典", Method: "PUT", Path: "/sysDictionary/updateSysDictionary", Description: "更新字典"},
		{ZC_MODEL: global.ZC_MODEL{ID: 75, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典", Method: "GET", Path: "/sysDictionary/findSysDictionary", Description: "根据ID获取字典"},
		{ZC_MODEL: global.ZC_MODEL{ID: 76, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "系统字典", Method: "GET", Path: "/sysDictionary/getSysDictionaryList", Description: "获取字典列表"},

		{ZC_MODEL: global.ZC_MODEL{ID: 77, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ZC_MODEL: global.ZC_MODEL{ID: 78, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ZC_MODEL: global.ZC_MODEL{ID: 79, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ZC_MODEL: global.ZC_MODEL{ID: 80, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ZC_MODEL: global.ZC_MODEL{ID: 81, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ZC_MODEL: global.ZC_MODEL{ID: 82, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "断点续传(插件版)", Method: "POST", Path: "/simpleUploader/upload", Description: "插件版分片上传"},
		{ZC_MODEL: global.ZC_MODEL{ID: 83, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/simpleUploader/checkFileMd5", Description: "文件完整度验证"},
		{ZC_MODEL: global.ZC_MODEL{ID: 84, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/simpleUploader/mergeFileMd5", Description: "上传完成合并文件"},

		{ZC_MODEL: global.ZC_MODEL{ID: 85, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "email", Method: "POST", Path: "/email/emailTest", Description: "发送测试邮件"},
		{ZC_MODEL: global.ZC_MODEL{ID: 86, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "email", Method: "POST", Path: "/email/emailSend", Description: "发送邮件示例"},

		{ZC_MODEL: global.ZC_MODEL{ID: 87, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "excel", Method: "POST", Path: "/excel/importExcel", Description: "导入excel"},
		{ZC_MODEL: global.ZC_MODEL{ID: 88, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "excel", Method: "GET", Path: "/excel/loadExcel", Description: "下载excel"},
		{ZC_MODEL: global.ZC_MODEL{ID: 89, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "excel", Method: "POST", Path: "/excel/exportExcel", Description: "导出excel"},
		{ZC_MODEL: global.ZC_MODEL{ID: 90, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "excel", Method: "GET", Path: "/excel/downloadTemplate", Description: "下载excel模板"},

		{ZC_MODEL: global.ZC_MODEL{ID: 91, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "设置按钮权限"},
		{ZC_MODEL: global.ZC_MODEL{ID: 92, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "获取已有按钮权限"},
		{ZC_MODEL: global.ZC_MODEL{ID: 93, CreatedAt: time.Now(), UpdatedAt: time.Now()}, ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "删除按钮"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/authorityBtn/canRemoveAuthorityBtn", "POST").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
