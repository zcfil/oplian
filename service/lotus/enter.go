package lotus

import (
	"oplian/service/lotus/deploy"
	"oplian/service/lotus/dispatch"
	"oplian/service/lotus/oplocal"
)

type ServiceGroup struct {
	dispatch.DispatchService
	oplocal.WorkerRunService
	oplocal.MinerService
	oplocal.OpLotusService
	oplocal.StorageService
	deploy.DeployService
	deploy.WorkerClusterService
	deploy.SectorsRecoverService
}
