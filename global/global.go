package global

import (
	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"oplian/config"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	ZC_DB     *gorm.DB
	ZC_CONFIG config.Server
	ZC_VP     *viper.Viper
	// ZC_LOG    *oplogging.Logger
	ROOM_CONFIG config.ServerRoom

	ZC_LOG                 *zap.Logger
	ZC_Concurrency_Control = &singleflight.Group{}

	OpUUID    uuid.UUID
	OpC2UUID  uuid.UUID
	GateWayID uuid.UUID
	LocalIP   string
)
