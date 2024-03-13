package example

import (
	"errors"

	"gorm.io/gorm"
	"oplian/global"
	"oplian/model/example"
)

type FileUploadAndDownloadService struct{}

//@function: FindOrCreateFile
//@description: When uploading a file, it detects the current file properties. If no file is created, the current slice of the file is returned
//@param: fileMd5 string, fileName string, chunkTotal int
//@return: file model.ExaFile, err error

func (e *FileUploadAndDownloadService) FindOrCreateFile(fileMd5 string, fileName string, chunkTotal int) (file example.ExaFile, err error) {
	var cfile example.ExaFile
	cfile.FileMd5 = fileMd5
	cfile.FileName = fileName
	cfile.ChunkTotal = chunkTotal

	if errors.Is(global.ZC_DB.Where("file_md5 = ? AND is_finish = ?", fileMd5, true).First(&file).Error, gorm.ErrRecordNotFound) {
		err = global.ZC_DB.Where("file_md5 = ? AND file_name = ?", fileMd5, fileName).Preload("ExaFileChunk").FirstOrCreate(&file, cfile).Error
		return file, err
	}
	cfile.IsFinish = true
	cfile.FilePath = file.FilePath
	err = global.ZC_DB.Create(&cfile).Error
	return cfile, err
}

//@function: CreateFileChunk
//@description: Create a file slicing record
//@param: id uint, fileChunkPath string, fileChunkNumber int
//@return: error

func (e *FileUploadAndDownloadService) CreateFileChunk(id uint, fileChunkPath string, fileChunkNumber int) error {
	var chunk example.ExaFileChunk
	chunk.FileChunkPath = fileChunkPath
	chunk.ExaFileID = id
	chunk.FileChunkNumber = fileChunkNumber
	err := global.ZC_DB.Create(&chunk).Error
	return err
}

//@function: DeleteFileChunk
//@description: Example Delete a file slice record
//@param: fileMd5 string, fileName string, filePath string
//@return: error

func (e *FileUploadAndDownloadService) DeleteFileChunk(fileMd5 string, filePath string) error {
	var chunks []example.ExaFileChunk
	var file example.ExaFile
	err := global.ZC_DB.Where("file_md5 = ? ", fileMd5).First(&file).
		Updates(map[string]interface{}{
			"IsFinish":  true,
			"file_path": filePath,
		}).Error
	if err != nil {
		return err
	}
	err = global.ZC_DB.Where("exa_file_id = ?", file.ID).Delete(&chunks).Unscoped().Error
	return err
}
