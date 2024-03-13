package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"oplian/global"
	. "oplian/model/system"
	"oplian/service/system"
	"time"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []SysBaseMenu{

		// oplian
		{ZC_MODEL: global.ZC_MODEL{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: Meta{Title: "homePage", Icon: "pie-chart"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "Configuration", Name: "Configuration", Component: "view/routerHolder.vue", Sort: 2, Meta: Meta{Title: "cmbdManagement", Icon: "operation"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "2", Path: "Config", Name: "Config", Component: "view/config/index.vue", Sort: 1, Meta: Meta{Title: "roomManagement", Icon: "platform"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "2", Path: "Property", Name: "Property", Component: "view/property/index.vue", Sort: 2, Meta: Meta{Title: "assetManagement", Icon: "coin"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "Platform", Name: "Platform", Component: "view/routerHolder.vue", Sort: 3, Meta: Meta{Title: "workPlatform", Icon: "collection"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "5", Path: "HistoryRecord", Name: "HistoryRecord", Component: "view/historyRecord/index.vue", Sort: 3, Meta: Meta{Title: "executionHistory", Icon: "odometer"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "5", Path: "FilePost", Name: "FilePost", Component: "view/filePost/index.vue", Sort: 2, Meta: Meta{Title: "fileDistribution", Icon: "folder-opened"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "5", Path: "ScriptExecution", Name: "ScriptExecution", Component: "view/scriptExecution/index.vue", Sort: 1, Meta: Meta{Title: "scriptExecution", Icon: "guide"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 9, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "5", Path: "FileManage", Name: "FileManage", Component: "view/fileManage/index.vue", Sort: 5, Meta: Meta{Title: "fileManagement", Icon: "files"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "environment", Name: "environment", Component: "view/routerHolder.vue", Sort: 4, Meta: Meta{Title: "testInspection", Icon: "help"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "10", Path: "testManage", Name: "testManage", Component: "view/testManage/index.vue", Sort: 1, Meta: Meta{Title: "testManagement", Icon: "wind-power"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "10", Path: "inspectionManage", Name: "inspectionManage", Component: "view/inspectionManage/index.vue", Sort: 2, Meta: Meta{Title: "inspectionManagement", Icon: "compass"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 13, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "automation", Name: "automation", Component: "view/routerHolder.vue", Sort: 5, Meta: Meta{Title: "deploymentPlatform", Icon: "promotion"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 14, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "13", Path: "lotusManage", Name: "lotusManage", Component: "view/lotusManage/index.vue", Sort: 2, Meta: Meta{Title: "lotusManagement", Icon: "connection"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "13", Path: "NodeManage", Name: "NodeManage", Component: "view/nodeManage/index.vue", Sort: 1, Meta: Meta{Title: "nodeInformation", Icon: "coordinate"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 16, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "13", Path: "minerManage", Name: "minerManage", Component: "view/minerManage/index.vue", Sort: 3, Meta: Meta{Title: "minerManagement", Icon: "set-up"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 17, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "13", Path: "workerManage", Name: "workerManage", Component: "view/workerManage/index.vue", Sort: 4, Meta: Meta{Title: "workerManagement", Icon: "refrigerator"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 18, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "13", Path: "storageManage", Name: "storageManage", Component: "view/storageManage/index.vue", Sort: 5, Meta: Meta{Title: "storageManagement", Icon: "school"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 19, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dispatchCenter", Name: "dispatchCenter", Component: "view/routerHolder.vue", Sort: 14, Meta: Meta{Title: "dispatchingCenter", Icon: "refresh"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 20, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "19", Path: "dispatchManage", Name: "dispatchManage", Component: "view/dispatchManage/index.vue", Sort: 1, Meta: Meta{Title: "schedulManagement", Icon: "flag"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 21, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "19", Path: "sectorsLifecycle", Name: "sectorsLifecycle", Component: "view/sectorsLifecycle/index.vue", Sort: 2, Meta: Meta{Title: "sectorLifeCycle", Icon: "orange"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 22, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "19", Path: "sectorJobList", Name: "sectorJobList", Component: "view/sectorJobList/index.vue", Sort: 3, Meta: Meta{Title: "sectorList", Icon: "cellphone"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 23, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "19", Path: "sectorRecover", Name: "sectorRecover", Component: "view/sectorRecover/index.vue", Sort: 5, Meta: Meta{Title: "sectorRecovery", Icon: "medal"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 24, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "monitoringCenter", Name: "monitoringCenter", Component: "view/routerHolder.vue", Sort: 15, Meta: Meta{Title: "monitoringCenter", Icon: "location-information"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 25, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "24", Path: "Equiment", Name: "Equiment", Component: "view/equipment/index.vue", Sort: 1, Meta: Meta{Title: "equipmentMonitor", Icon: "location-information"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 26, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "24", Path: "businessMonitoring", Name: "businessMonitoring", Component: "view/businessMonitoring/index.vue", Sort: 2, Meta: Meta{Title: "businessMonitor", Icon: "message-box"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 27, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "WarnCenter", Name: "WarnCenter", Component: "view/routerHolder.vue", Sort: 16, Meta: Meta{Title: "alarmCenter", Icon: "warning"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 28, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "27", Path: "Warning", Name: "Warning", Component: "view/warning/index.vue", Sort: 1, Meta: Meta{Title: "alarmCenter", Icon: "warning-filled"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 29, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "27", Path: "WarnList", Name: "WarnList", Component: "view/warnList/index.vue", Sort: 2, Meta: Meta{Title: "alarmList", Icon: "list"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 30, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "27", Path: "WarnTactics", Name: "WarnTactics", Component: "view/warnTactics/index.vue", Sort: 3, Meta: Meta{Title: "alarmManagement", Icon: "grid"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 31, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "colonyPlatform", Name: "colonyPlatform", Component: "view/routerHolder.vue", Sort: 18, Meta: Meta{Title: "c2Platform", Icon: "message-box"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 32, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: true, ParentId: "31", Path: "taskMarket", Name: "taskMarket", Component: "view/taskMarket/index.vue", Sort: 1, Meta: Meta{Title: "c2Market", Icon: "credit-card"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 33, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "31", Path: "myTask", Name: "myTask", Component: "view/myTask/index.vue", Sort: 2, Meta: Meta{Title: "c2Cluster", Icon: "set-up"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 34, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: true, ParentId: "31", Path: "AllTask", Name: "AllTask", Component: "view/allTask/index.vue", Sort: 0, Meta: Meta{Title: "c2Tasks", Icon: "aim"}},

		{ZC_MODEL: global.ZC_MODEL{ID: 35, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 23, Meta: Meta{Title: "systemManagement", Icon: "user"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 36, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "35", Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "roleManagement", Icon: "avatar"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 37, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "35", Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 2, Meta: Meta{Title: "userManagement", Icon: "coordinate"}},
		{ZC_MODEL: global.ZC_MODEL{ID: 38, CreatedAt: time.Now(), UpdatedAt: time.Now()}, MenuLevel: 0, Hidden: false, ParentId: "35", Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 1, Meta: Meta{Title: "operationHistory", Icon: "pie-chart"}},
	}
	if err = db.Save(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"表数据初始化失败!")
	}

	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "autoPkg").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
