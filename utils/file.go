package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"time"
)

type FileInfo struct {
	Path     string
	FileName string
	FileData []byte
}

func FileUpLoad(c *gin.Context, path string) (m map[string]int64, e error) {

	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	if form == nil {
		return nil, errors.New("files is nil")
	}
	files := form.File["files"]
	if files == nil {
		return nil, errors.New("files is nil")
	}

	err = AddDir(path)
	if err != nil {
		return nil, err
	}

	fileMap := make(map[string]int64)
	for _, file := range files {

		dst := path + "/" + file.Filename
		err := SaveUploadedFile(file, dst)
		if err != nil {
			return nil, err
		}
		fileMap[file.Filename] = file.Size
	}

	return fileMap, nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return err
	}
	return err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func AddDir(filePath string) error {

	if _, err := os.Stat(filePath); err != nil {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportExcelFile(c *gin.Context, titleList, colList []string, m map[int]map[string]string, fileName string) {

	file := xlsx.NewFile()

	sheet, _ := file.AddSheet("Sheet1")

	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
	}

	var ar []int
	for k, _ := range m {
		ar = append(ar, k)
	}
	sort.Ints(ar)

	for i := range ar {
		row := sheet.AddRow()
		for _, v := range colList {
			cell := row.AddCell()
			cell.Value = m[i][v]
		}
	}

	var buffer bytes.Buffer
	_ = file.Write(&buffer)
	content := bytes.NewReader(buffer.Bytes())

	fileName = fmt.Sprintf("%s%s%s.xlsx", fileName, `-`, GetNowStr())
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), content)
}

func CreateFile(file FileInfo) error {

	if ok, _ := PathExists(file.Path); !ok {
		err := os.MkdirAll(file.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fileUrl := file.Path + "/" + file.FileName
	f, err := os.OpenFile(fileUrl, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}

	_, err = f.Write(file.FileData)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func DelFile(fileName string) error {
	return os.Remove(fileName)
}

func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}
