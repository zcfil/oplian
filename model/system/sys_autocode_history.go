package system

import (
	"strconv"
	"strings"

	"oplian/global"
	"oplian/model/common/request"
)

type SysAutoCodeHistory struct {
	global.ZC_MODEL
	Package       string `json:"package"`
	BusinessDB    string `json:"businessDB"`
	TableName     string `json:"tableName"`
	RequestMeta   string `gorm:"type:text" json:"requestMeta,omitempty"`
	AutoCodePath  string `gorm:"type:text" json:"autoCodePath,omitempty"`
	InjectionMeta string `gorm:"type:text" json:"injectionMeta,omitempty"`
	StructName    string `json:"structName"`
	StructCNName  string `json:"structCNName"`
	ApiIDs        string `json:"apiIDs,omitempty"`
	Flag          int    `json:"flag"`
}

func (m *SysAutoCodeHistory) ToRequestIds() request.IdsReq {
	if m.ApiIDs == "" {
		return request.IdsReq{}
	}
	slice := strings.Split(m.ApiIDs, ";")
	ids := make([]int, 0, len(slice))
	length := len(slice)
	for i := 0; i < length; i++ {
		id, _ := strconv.ParseInt(slice[i], 10, 32)
		ids = append(ids, int(id))
	}
	return request.IdsReq{Ids: ids}
}
