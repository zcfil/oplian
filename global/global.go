package global

import (
	"sync"

	"oplian/utils/timer"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"oplian/config"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	ZC_DB     *gorm.DB
	ZC_DBList map[string]*gorm.DB
	ZC_CONFIG config.Server
	ZC_VP     *viper.Viper
	// ZC_LOG    *oplogging.Logger
	ROOM_CONFIG config.ServerRoom

	ZC_LOG                 *zap.Logger
	ZC_Timer               timer.Timer = timer.NewTimerTask()
	ZC_Concurrency_Control             = &singleflight.Group{}

	BlackCache map[string]int64
	BlackLock  sync.RWMutex

	lock      sync.RWMutex
	OpUUID    uuid.UUID
	OpC2UUID  uuid.UUID
	GateWayID uuid.UUID
	LocalIP   string
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return ZC_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := ZC_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
