package example

import "oplian/service"

type ApiGroup struct {
	FileUploadAndDownloadApi
}

var (
	FileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
