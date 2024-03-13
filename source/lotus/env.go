package lotus

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	sysModel "oplian/model/lotus"
	"oplian/service/system"
)

const InitZc = 100
const initOrderEnv = InitZc + 1

type initEnv struct{}

// auto run
func init() {
	system.RegisterInit(initOrderEnv, &initEnv{})
}

func (i *initEnv) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.LotusEnv{})
}

func (i *initEnv) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.LotusEnv{})
}

func (i initEnv) InitializerName() string {
	return sysModel.LotusEnv{}.TableName()
}

func (i *initEnv) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	entities := []sysModel.LotusEnv{
		{
			Ekey:    "RUST_BACKTRACE",
			Evalue:  "full",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "set rust log",
		},
		{
			Ekey:    "RUSTFLAGS",
			Evalue:  "-C target-cpu=native -g",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "rust instruction set",
		},
		{
			Ekey:    "FFI_BUILD_FROM_SOURCE",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
		},
		{
			Ekey:    "RUST_LOG",
			Evalue:  "info",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "rust log format",
		},
		{
			Ekey:    "LOTUS_PATH",
			Evalue:  "/mnt/md0/ipfs/data/lotus",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "lotus data path",
		},
		{
			Ekey:    "LOTUS_MINER_PATH",
			Evalue:  "/mnt/md0/ipfs/data/lotusminer",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "miner data path",
		},
		{
			Ekey:    "IPFS_GATEWAY",
			Evalue:  "https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "ipfs gateway",
		},
		{
			Ekey:    "FIL_PROOFS_PARAMETER_CACHE",
			Evalue:  "/mnt/md0/filecoin-proof-parameters",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "proof parameters data path",
		},
		{
			Ekey:    "WORKER_PATH",
			Evalue:  "/mnt/md0/filecoin-proof-parameters",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "worker data path",
		},
		{
			Ekey:    "TRUST_PARAMS",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
		},
		{
			Ekey:    "TRUST_PARAMS_FORCE",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
		},
		{
			Ekey:    "FIL_PROOFS_USE_GPU_TREE_BUILDER",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "use gpu",
		},
		{
			Ekey:    "FIL_PROOFS_USE_MULTICORE_SDR",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "cup use multicore",
		},
		{
			Ekey:    "P1_CORES",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "add worker thread (1 recommended)",
		},
		{
			Ekey:    "FIL_PROOFS_MAXIMIZE_CACHING",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
		},
		{
			Ekey:    "FIL_PROOFS_MULTICORE_SDR_PRODUCERS",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "cup use multicore",
		},
		{
			Ekey:    "MANAGE_C2",
			Evalue:  "1",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "cup use multicore",
		},
		{
			Ekey:    "URL_C2",
			Evalue:  "192.168.1.1:4567,192.168.1.2:4567",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "C2 remote cluster URLs separated by commas",
		},
		{
			Ekey:    "TOKEN_C2",
			Evalue:  "",
			MinerId: "default",
			Etype:   sysModel.Default,
			Remark:  "C2 token",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.LotusEnv{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)

	return next, err
}

func (i *initEnv) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var record sysModel.LotusEnv
	if errors.Is(db.Where("miner_id = ?", "default").First(&record).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
