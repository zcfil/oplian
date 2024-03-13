package system

import (
	"errors"
	"go/token"
	"strings"

	"oplian/global"
)

// AutoCodeStruct 初始版本自动化代码工具
type AutoCodeStruct struct {
	StructName         string   `json:"structName"`
	TableName          string   `json:"tableName"`
	PackageName        string   `json:"packageName"`
	HumpPackageName    string   `json:"humpPackageName"`
	Abbreviation       string   `json:"abbreviation"`
	Description        string   `json:"description"`
	AutoCreateApiToSql bool     `json:"autoCreateApiToSql"`
	AutoCreateResource bool     `json:"autoCreateResource"`
	AutoMoveFile       bool     `json:"autoMoveFile"`
	BusinessDB         string   `json:"businessDB"`
	Fields             []*Field `json:"fields,omitempty"`
	HasTimer           bool
	DictTypes          []string `json:"-"`
	Package            string   `json:"package"`
	PackageT           string   `json:"-"`
	NeedValid          bool     `json:"-"`
}

func (a *AutoCodeStruct) Pretreatment() {
	a.KeyWord()
	a.SuffixTest()
}

func (a *AutoCodeStruct) KeyWord() {
	if token.IsKeyword(a.Abbreviation) {
		a.Abbreviation = a.Abbreviation + "_"
	}
}

func (a *AutoCodeStruct) SuffixTest() {
	if strings.HasSuffix(a.HumpPackageName, "test") {
		a.HumpPackageName = a.HumpPackageName + "_"
	}
}

type Field struct {
	FieldName       string `json:"fieldName"`       // Field名
	FieldDesc       string `json:"fieldDesc"`       // 中文名
	FieldType       string `json:"fieldType"`       // Field数据类型
	FieldJson       string `json:"fieldJson"`       // FieldJson
	DataTypeLong    string `json:"dataTypeLong"`    // 数据库字段长度
	Comment         string `json:"comment"`         // 数据库字段描述
	ColumnName      string `json:"columnName"`      // 数据库字段
	FieldSearchType string `json:"fieldSearchType"` // 搜索条件
	DictType        string `json:"dictType"`        // 字典
	Require         bool   `json:"require"`         // 是否必填
	ErrorText       string `json:"errorText"`       // 校验失败文字
	Clearable       bool   `json:"clearable"`       // 是否可清空
}

var ErrAutoMove error = errors.New("code created successfully and file moved successfully")

type SysAutoCode struct {
	global.ZC_MODEL
	PackageName string `json:"packageName" gorm:"comment:包名"`
	Label       string `json:"label" gorm:"comment:展示名"`
	Desc        string `json:"desc" gorm:"comment:描述"`
}

type AutoPlugReq struct {
	PlugName    string         `json:"plugName"`
	Snake       string         `json:"snake"`
	RouterGroup string         `json:"routerGroup"`
	HasGlobal   bool           `json:"hasGlobal"`
	HasRequest  bool           `json:"hasRequest"`
	HasResponse bool           `json:"hasResponse"`
	NeedModel   bool           `json:"needModel"`
	Global      []AutoPlugInfo `json:"global,omitempty"`
	Request     []AutoPlugInfo `json:"request,omitempty"`
	Response    []AutoPlugInfo `json:"response,omitempty"`
}

func (a *AutoPlugReq) CheckList() {
	a.Global = bind(a.Global)
	a.Request = bind(a.Request)
	a.Response = bind(a.Response)

}
func bind(req []AutoPlugInfo) []AutoPlugInfo {
	var r []AutoPlugInfo
	for _, info := range req {
		if info.Effective() {
			r = append(r, info)
		}
	}
	return r
}

type AutoPlugInfo struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	Desc string `json:"desc"`
}

func (a AutoPlugInfo) Effective() bool {
	return a.Key != "" && a.Type != "" && a.Desc != ""
}
