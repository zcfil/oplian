package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"oplian/global"
	sysModel "oplian/model/system"
	"oplian/service/system"
	"time"
)

const initOrderPatrolConfig = initOrderAuthority + 1

type initPatrolConfig struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPatrolConfig, &initPatrolConfig{})
}

func (i *initPatrolConfig) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysPatrolConfig{})
}

func (i *initPatrolConfig) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysPatrolConfig{})
}

func (i initPatrolConfig) InitializerName() string {
	return sysModel.SysPatrolConfig{}.TableName()
}

func (i *initPatrolConfig) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.SysPatrolConfig{
		{ZC_MODEL: global.ZC_MODEL{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, PatrolType: 1, IntervalHours: 1, IntervalTime: 3600},
		{ZC_MODEL: global.ZC_MODEL{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, PatrolType: 2, IntervalHours: 1, IntervalTime: 3600},
		{ZC_MODEL: global.ZC_MODEL{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, PatrolType: 3, IntervalHours: 1, IntervalTime: 3600},
	}
	if err = db.Create(&entities).Error; err != nil {
		log.Println("create data error：", err)
		return ctx, errors.Wrap(err, sysModel.SysPatrolConfig{}.TableName()+"表数据初始化失败!")
	}
	log.Println("create data success：", entities)
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initPatrolConfig) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("patrol_type = ? ", "1").
		First(&sysModel.SysPatrolConfig{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
