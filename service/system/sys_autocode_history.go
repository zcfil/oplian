package system

import (
	"errors"
	"fmt"
	systemReq "oplian/model/system/request"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"oplian/model/system/response"

	"oplian/global"
	"oplian/model/common/request"
	"oplian/model/system"
	"oplian/utils"

	"go.uber.org/zap"
)

var RepeatErr = errors.New("重复创建")

type AutoCodeHistoryService struct{}

var AutoCodeHistoryServiceApp = new(AutoCodeHistoryService)

// CreateAutoCodeHistory Create a code generator history
// RouterPath : RouterPath@RouterString;RouterPath2@RouterString2
func (autoCodeHistoryService *AutoCodeHistoryService) CreateAutoCodeHistory(meta, structName, structCNName, autoCodePath string, injectionMeta string, tableName string, apiIds string, Package string) error {
	return global.ZC_DB.Create(&system.SysAutoCodeHistory{
		Package:       Package,
		RequestMeta:   meta,
		AutoCodePath:  autoCodePath,
		InjectionMeta: injectionMeta,
		StructName:    structName,
		StructCNName:  structCNName,
		TableName:     tableName,
		ApiIDs:        apiIds,
	}).Error
}

// First Gets data for code generator history based on id

func (autoCodeHistoryService *AutoCodeHistoryService) First(info *request.GetById) (string, error) {
	var meta string
	return meta, global.ZC_DB.Model(system.SysAutoCodeHistory{}).Select("request_meta").Where("id = ?", info.Uint()).First(&meta).Error
}

// Repeat Detection duplication

func (autoCodeHistoryService *AutoCodeHistoryService) Repeat(businessDB, structName, Package string) bool {
	var count int64
	global.ZC_DB.Model(&system.SysAutoCodeHistory{}).Where("business_db = ? and struct_name = ? and package = ? and flag = 0", businessDB, structName, Package).Count(&count)
	return count > 0
}

// RollBack rollback

func (autoCodeHistoryService *AutoCodeHistoryService) RollBack(info *systemReq.RollBack) error {
	md := system.SysAutoCodeHistory{}
	if err := global.ZC_DB.Where("id = ?", info.ID).First(&md).Error; err != nil {
		return err
	}

	ids := request.IdsReq{}
	idsStr := strings.Split(md.ApiIDs, ";")
	for i := range idsStr[0 : len(idsStr)-1] {
		id, err := strconv.Atoi(idsStr[i])
		if err != nil {
			return err
		}
		ids.Ids = append(ids.Ids, id)
	}
	err := ApiServiceApp.DeleteApisByIds(ids)
	if err != nil {
		global.ZC_LOG.Error("ClearTag DeleteApiByIds:", zap.Error(err))
	}

	if info.DeleteTable {
		if err = AutoCodeServiceApp.DropTable(md.BusinessDB, md.TableName); err != nil {
			global.ZC_LOG.Error("ClearTag DropTable:", zap.Error(err))
		}
	}

	for _, path := range strings.Split(md.AutoCodePath, ";") {

		_path, err := filepath.Abs(path)
		if err != nil || _path != path {
			continue
		}

		nPath := filepath.Join(global.ZC_CONFIG.AutoCode.Root,
			"rm_file", time.Now().Format("20060102"), filepath.Base(filepath.Dir(filepath.Dir(path))), filepath.Base(filepath.Dir(path)), filepath.Base(path))

		for utils.FileExist(nPath) {

			nPath += fmt.Sprintf("_%d", time.Now().Nanosecond())
		}
		err = utils.FileMove(path, nPath)
		if err != nil {
			global.ZC_LOG.Error("file move err ", zap.Error(err))
		}
		//_ = utils.DeLFile(path)
	}

	for _, v := range strings.Split(md.InjectionMeta, ";") {
		// RouterPath@functionName@RouterString
		meta := strings.Split(v, "@")
		if len(meta) == 3 {
			_ = utils.AutoClearCode(meta[0], meta[2])
		}
	}
	md.Flag = 1
	return global.ZC_DB.Save(&md).Error
}

// Delete Delete historical data

func (autoCodeHistoryService *AutoCodeHistoryService) Delete(info *request.GetById) error {
	return global.ZC_DB.Where("id = ?", info.Uint()).Delete(&system.SysAutoCodeHistory{}).Error
}

// GetList Obtain historical system data

func (autoCodeHistoryService *AutoCodeHistoryService) GetList(info request.PageInfo) (list []response.AutoCodeHistory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.ZC_DB.Model(&system.SysAutoCodeHistory{})
	var entities []response.AutoCodeHistory
	err = db.Count(&total).Error
	if err != nil {
		return nil, total, err
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&entities).Error
	return entities, total, err
}
