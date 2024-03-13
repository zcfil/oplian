package lotus

import "oplian/service"

type ApiGroup struct {
	DispatchApi
	DeployApi
	WorkerClusterApi
	SectorsRecoverApi
}

var (
	DeployService   = service.ServiceGroupApp.LotusServiceGroup.DeployService
	dispatchService = service.ServiceGroupApp.LotusServiceGroup.DispatchService
)
