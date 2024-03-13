package example

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/common/response"
	"oplian/model/example"
	exampleRes "oplian/model/example/response"
)

type FileUploadAndDownloadApi struct{}

// UploadFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Example of uploading files
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file
// @Success   200   {object}  response.Response{data=exampleRes.ExaFileResponse,msg=string}
// @Router    /fileUploadAndDownload/upload [post]
func (b *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.ZC_LOG.Error("Received file failed!", zap.Error(err))
		response.FailWithMessage("Received file failed", c)
		return
	}
	file, err = FileUploadAndDownloadService.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		global.ZC_LOG.Error("Failed to modify database link!", zap.Error(err))
		response.FailWithMessage("Failed to modify database link", c)
		return
	}
	response.OkWithDetailed(exampleRes.ExaFileResponse{File: file}, "Upload successful", c)
}

// EditFileName Edit file name or comment
func (b *FileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = FileUploadAndDownloadService.EditFileName(file)
	if err != nil {
		global.ZC_LOG.Error("Editing failed!", zap.Error(err))
		response.FailWithMessage("Editing failed", c)
		return
	}
	response.OkWithMessage("Edit successful", c)
}

// DeleteFile
// @Tags      ExaFileUploadAndDownload
// @Summary   delete file
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{msg=string}
// @Router    /fileUploadAndDownload/deleteFile [post]
func (b *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := FileUploadAndDownloadService.DeleteFile(file); err != nil {
		global.ZC_LOG.Error("Delete failed!", zap.Error(err))
		response.FailWithMessage("Delete failed", c)
		return
	}
	response.OkWithMessage("Delete successful", c)
}

// GetFileList
// @Tags      ExaFileUploadAndDownload
// @Summary   Page file list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}
// @Router    /fileUploadAndDownload/getFileList [post]
func (b *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := FileUploadAndDownloadService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.ZC_LOG.Error("Acquisition failed!", zap.Error(err))
		response.FailWithMessage("Acquisition failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Successfully obtained", c)
}
