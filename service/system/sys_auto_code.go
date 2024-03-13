package system

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"oplian/resource/autocode_template/subcontract"

	"oplian/global"
	"oplian/model/system"
	"oplian/utils"

	"gorm.io/gorm"
)

const (
	autoPath           = "autocode_template/"
	autocodePath       = "resource/autocode_template"
	plugPath           = "resource/plug_template"
	packageService     = "service/%s/enter.go"
	packageServiceName = "service"
	packageRouter      = "router/%s/enter.go"
	packageRouterName  = "router"
	packageAPI         = "api/v1/%s/enter.go"
	packageAPIName     = "api/v1"
)

type autoPackage struct {
	path string
	temp string
	name string
}

var (
	packageInjectionMap map[string]astInjectionMeta
	injectionPaths      []injectionMeta
	caser               = cases.Title(language.English)
)

func Init(Package string) {
	injectionPaths = []injectionMeta{
		{
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.SInitialize, "gorm.go"),
			funcName:    "MysqlTables",
			structNameF: Package + ".%s{},",
		},
		{
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.SInitialize, "router.go"),
			funcName:    "Routers",
			structNameF: Package + "Router.Init%sRouter(PrivateGroup)",
		},
		{
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SApi, Package), "enter.go"),
			funcName:    "ApiGroup",
			structNameF: "%sApi",
		},
		{
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SRouter, Package), "enter.go"),
			funcName:    "RouterGroup",
			structNameF: "%sRouter",
		},
		{
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SService, Package), "enter.go"),
			funcName:    "ServiceGroup",
			structNameF: "%sService",
		},
	}

	packageInjectionMap = map[string]astInjectionMeta{
		packageServiceName: {
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, "service", "enter.go"),
			importCodeF:  "oplian/%s/%s",
			packageNameF: "%s",
			groupName:    "ServiceGroup",
			structNameF:  "%sServiceGroup",
		},
		packageRouterName: {
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, "router", "enter.go"),
			importCodeF:  "oplian/%s/%s",
			packageNameF: "%s",
			groupName:    "RouterGroup",
			structNameF:  "%s",
		},
		packageAPIName: {
			path: filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, "api/v1", "enter.go"),
			importCodeF:  "oplian/%s/%s",
			packageNameF: "%s",
			groupName:    "ApiGroup",
			structNameF:  "%sApiGroup",
		},
	}
}

type injectionMeta struct {
	path        string
	funcName    string
	structNameF string
}

type astInjectionMeta struct {
	path         string
	importCodeF  string
	structNameF  string
	packageNameF string
	groupName    string
}

type tplData struct {
	template         *template.Template
	autoPackage      string
	locationPath     string
	autoCodePath     string
	autoMoveFilePath string
}

type AutoCodeService struct{}

var AutoCodeServiceApp = new(AutoCodeService)

func (autoCodeService *AutoCodeService) PreviewTemp(autoCode system.AutoCodeStruct) (map[string]string, error) {
	makeDictTypes(&autoCode)
	for i := range autoCode.Fields {
		if autoCode.Fields[i].FieldType == "time.Time" {
			autoCode.HasTimer = true
			break
		}
		if autoCode.Fields[i].Require {
			autoCode.NeedValid = true
			break
		}
	}
	dataList, _, needMkdir, err := autoCodeService.getNeedList(&autoCode)
	if err != nil {
		return nil, err
	}

	// Create a folder before writing files
	if err = utils.CreateDir(needMkdir...); err != nil {
		return nil, err
	}

	// Create a map
	ret := make(map[string]string)

	// Generate map
	for _, value := range dataList {
		ext := ""
		if ext = filepath.Ext(value.autoCodePath); ext == ".txt" {
			continue
		}
		f, err := os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_WRONLY, 0o755)
		if err != nil {
			return nil, err
		}
		if err = value.template.Execute(f, autoCode); err != nil {
			return nil, err
		}
		_ = f.Close()
		f, err = os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_RDONLY, 0o755)
		if err != nil {
			return nil, err
		}
		builder := strings.Builder{}
		builder.WriteString("```")

		if ext != "" && strings.Contains(ext, ".") {
			builder.WriteString(strings.Replace(ext, ".", "", -1))
		}
		builder.WriteString("\n\n")
		data, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		builder.Write(data)
		builder.WriteString("\n\n```")

		pathArr := strings.Split(value.autoCodePath, string(os.PathSeparator))
		ret[pathArr[1]+"-"+pathArr[3]] = builder.String()
		_ = f.Close()

	}
	defer func() {
		if err := os.RemoveAll(autoPath); err != nil {
			return
		}
	}()
	return ret, nil
}

func makeDictTypes(autoCode *system.AutoCodeStruct) {
	DictTypeM := make(map[string]string)
	for _, v := range autoCode.Fields {
		if v.DictType != "" {
			DictTypeM[v.DictType] = ""
		}
	}

	for k := range DictTypeM {
		autoCode.DictTypes = append(autoCode.DictTypes, k)
	}
}

func (autoCodeService *AutoCodeService) CreateTemp(autoCode system.AutoCodeStruct, ids ...uint) (err error) {
	makeDictTypes(&autoCode)
	for i := range autoCode.Fields {
		if autoCode.Fields[i].FieldType == "time.Time" {
			autoCode.HasTimer = true
			break
		}
		if autoCode.Fields[i].Require {
			autoCode.NeedValid = true
			break
		}
	}

	if autoCode.AutoMoveFile && AutoCodeHistoryServiceApp.Repeat(autoCode.BusinessDB, autoCode.StructName, autoCode.Package) {
		return RepeatErr
	}
	dataList, fileList, needMkdir, err := autoCodeService.getNeedList(&autoCode)
	if err != nil {
		return err
	}
	meta, _ := json.Marshal(autoCode)

	if err = utils.CreateDir(needMkdir...); err != nil {
		return err
	}

	for _, value := range dataList {
		f, err := os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_WRONLY, 0o755)
		if err != nil {
			return err
		}
		if err = value.template.Execute(f, autoCode); err != nil {
			return err
		}
		_ = f.Close()
	}

	defer func() {
		if err := os.RemoveAll(autoPath); err != nil {
			return
		}
	}()
	bf := strings.Builder{}
	idBf := strings.Builder{}
	injectionCodeMeta := strings.Builder{}
	for _, id := range ids {
		idBf.WriteString(strconv.Itoa(int(id)))
		idBf.WriteString(";")
	}
	if autoCode.AutoMoveFile {
		Init(autoCode.Package)
		for index := range dataList {
			autoCodeService.addAutoMoveFile(&dataList[index])
		}

		for _, value := range dataList {
			if utils.FileExist(value.autoMoveFilePath) {
				return errors.New(fmt.Sprintf("the target file already exists:%s\n", value.autoMoveFilePath))
			}
		}
		for _, value := range dataList {
			if err := utils.FileMove(value.autoCodePath, value.autoMoveFilePath); err != nil {
				return err
			}
		}
		err = injectionCode(autoCode.StructName, &injectionCodeMeta)
		if err != nil {
			return
		}
		// 保存生成信息
		for _, data := range dataList {
			if len(data.autoMoveFilePath) != 0 {
				bf.WriteString(data.autoMoveFilePath)
				bf.WriteString(";")
			}
		}

		var gormPath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
			global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.SInitialize, "gorm.go")
		var routePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
			global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.SInitialize, "router.go")
		var imporStr = fmt.Sprintf("oplian/model/%s", autoCode.Package)
		_ = ImportReference(routePath, "", "", autoCode.Package, "")
		_ = ImportReference(gormPath, imporStr, "", "", "")

	} else {
		if err = utils.ZipFiles("./ginvueadmin.zip", fileList, ".", "."); err != nil {
			return err
		}
	}
	if autoCode.AutoMoveFile || autoCode.AutoCreateApiToSql {
		if autoCode.TableName != "" {
			err = AutoCodeHistoryServiceApp.CreateAutoCodeHistory(
				string(meta),
				autoCode.StructName,
				autoCode.Description,
				bf.String(),
				injectionCodeMeta.String(),
				autoCode.TableName,
				idBf.String(),
				autoCode.Package,
			)
		} else {
			err = AutoCodeHistoryServiceApp.CreateAutoCodeHistory(
				string(meta),
				autoCode.StructName,
				autoCode.Description,
				bf.String(),
				injectionCodeMeta.String(),
				autoCode.StructName,
				idBf.String(),
				autoCode.Package,
			)
		}
	}
	if err != nil {
		return err
	}
	if autoCode.AutoMoveFile {
		return system.ErrAutoMove
	}
	return nil
}

func (autoCodeService *AutoCodeService) GetAllTplFile(pathName string, fileList []string) ([]string, error) {
	files, err := os.ReadDir(pathName)
	for _, fi := range files {
		if fi.IsDir() {
			fileList, err = autoCodeService.GetAllTplFile(pathName+"/"+fi.Name(), fileList)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.HasSuffix(fi.Name(), ".tpl") {
				fileList = append(fileList, pathName+"/"+fi.Name())
			}
		}
	}
	return fileList, err
}

func (autoCodeService *AutoCodeService) DropTable(BusinessDb, tableName string) error {
	if BusinessDb != "" {
		return global.ZC_DB.Exec("DROP TABLE " + tableName).Error
	} else {
		return global.MustGetGlobalDBByDBName(BusinessDb).Exec("DROP TABLE " + tableName).Error
	}
}

func (autoCodeService *AutoCodeService) addAutoMoveFile(data *tplData) {
	base := filepath.Base(data.autoCodePath)
	fileSlice := strings.Split(data.autoCodePath, string(os.PathSeparator))
	n := len(fileSlice)
	if n <= 2 {
		return
	}
	if strings.Contains(fileSlice[1], "server") {
		if strings.Contains(fileSlice[n-2], "router") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root, global.ZC_CONFIG.AutoCode.Server,
				fmt.Sprintf(global.ZC_CONFIG.AutoCode.SRouter, data.autoPackage), base)
		} else if strings.Contains(fileSlice[n-2], "api") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SApi, data.autoPackage), base)
		} else if strings.Contains(fileSlice[n-2], "service") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SService, data.autoPackage), base)
		} else if strings.Contains(fileSlice[n-2], "model") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SModel, data.autoPackage), base)
		} else if strings.Contains(fileSlice[n-2], "request") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SRequest, data.autoPackage), base)
		}
	} else if strings.Contains(fileSlice[1], "web") {
		if strings.Contains(fileSlice[n-1], "js") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Web, global.ZC_CONFIG.AutoCode.WApi, base)
		} else if strings.Contains(fileSlice[n-2], "form") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Web, global.ZC_CONFIG.AutoCode.WForm, filepath.Base(filepath.Dir(filepath.Dir(data.autoCodePath))), strings.TrimSuffix(base, filepath.Ext(base))+"Form.vue")
		} else if strings.Contains(fileSlice[n-2], "table") {
			data.autoMoveFilePath = filepath.Join(global.ZC_CONFIG.AutoCode.Root,
				global.ZC_CONFIG.AutoCode.Web, global.ZC_CONFIG.AutoCode.WTable, filepath.Base(filepath.Dir(filepath.Dir(data.autoCodePath))), base)
		}
	}
}

func (autoCodeService *AutoCodeService) AutoCreateApi(a *system.AutoCodeStruct) (ids []uint, err error) {
	apiList := []system.SysApi{
		{
			Path:        "/" + a.Abbreviation + "/" + "create" + a.StructName,
			Description: "新增" + a.Description,
			ApiGroup:    a.Abbreviation,
			Method:      "POST",
		},
		{
			Path:        "/" + a.Abbreviation + "/" + "delete" + a.StructName,
			Description: "删除" + a.Description,
			ApiGroup:    a.Abbreviation,
			Method:      "DELETE",
		},
		{
			Path:        "/" + a.Abbreviation + "/" + "delete" + a.StructName + "ByIds",
			Description: "批量删除" + a.Description,
			ApiGroup:    a.Abbreviation,
			Method:      "DELETE",
		},
		{
			Path:        "/" + a.Abbreviation + "/" + "update" + a.StructName,
			Description: "更新" + a.Description,
			ApiGroup:    a.Abbreviation,
			Method:      "PUT",
		},
		{
			Path:        "/" + a.Abbreviation + "/" + "find" + a.StructName,
			Description: "根据ID获取" + a.Description,
			ApiGroup:    a.Abbreviation,
			Method:      "GET",
		},
		{
			Path:        "/" + a.Abbreviation + "/" + "get" + a.StructName + "List",
			Description: "获取" + a.Description + "列表",
			ApiGroup:    a.Abbreviation,
			Method:      "GET",
		},
	}
	err = global.ZC_DB.Transaction(func(tx *gorm.DB) error {
		for _, v := range apiList {
			var api system.SysApi
			if errors.Is(tx.Where("path = ? AND method = ?", v.Path, v.Method).First(&api).Error, gorm.ErrRecordNotFound) {
				if err = tx.Create(&v).Error; err != nil { // 遇到错误时回滚事务
					return err
				} else {
					ids = append(ids, v.ID)
				}
			}
		}
		return nil
	})
	return ids, err
}

func (autoCodeService *AutoCodeService) getNeedList(autoCode *system.AutoCodeStruct) (dataList []tplData, fileList []string, needMkdir []string, err error) {

	utils.TrimSpace(autoCode)
	for _, field := range autoCode.Fields {
		utils.TrimSpace(field)
	}

	tplFileList, err := autoCodeService.GetAllTplFile(autocodePath, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	dataList = make([]tplData, 0, len(tplFileList))
	fileList = make([]string, 0, len(tplFileList))
	needMkdir = make([]string, 0, len(tplFileList))

	for _, value := range tplFileList {
		dataList = append(dataList, tplData{locationPath: value, autoPackage: autoCode.Package})
	}

	for index, value := range dataList {
		dataList[index].template, err = template.ParseFiles(value.locationPath)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	for index, value := range dataList {
		trimBase := strings.TrimPrefix(value.locationPath, autocodePath+"/")
		if trimBase == "readme.txt.tpl" {
			dataList[index].autoCodePath = autoPath + "readme.txt"
			continue
		}

		if lastSeparator := strings.LastIndex(trimBase, "/"); lastSeparator != -1 {
			origFileName := strings.TrimSuffix(trimBase[lastSeparator+1:], ".tpl")
			firstDot := strings.Index(origFileName, ".")
			if firstDot != -1 {
				var fileName string
				if origFileName[firstDot:] != ".go" {
					fileName = autoCode.PackageName + origFileName[firstDot:]
				} else {
					fileName = autoCode.HumpPackageName + origFileName[firstDot:]
				}

				dataList[index].autoCodePath = filepath.Join(autoPath, trimBase[:lastSeparator], autoCode.PackageName,
					origFileName[:firstDot], fileName)
			}
		}

		if lastSeparator := strings.LastIndex(dataList[index].autoCodePath, string(os.PathSeparator)); lastSeparator != -1 {
			needMkdir = append(needMkdir, dataList[index].autoCodePath[:lastSeparator])
		}
	}
	for _, value := range dataList {
		fileList = append(fileList, value.autoCodePath)
	}
	return dataList, fileList, needMkdir, err
}

func injectionCode(structName string, bf *strings.Builder) error {
	for _, meta := range injectionPaths {
		code := fmt.Sprintf(meta.structNameF, structName)
		if err := utils.AutoInjectionCode(meta.path, meta.funcName, code); err != nil {
			return err
		}
		bf.WriteString(fmt.Sprintf("%s@%s@%s;", meta.path, meta.funcName, code))
	}
	return nil
}

func (autoCodeService *AutoCodeService) CreateAutoCode(s *system.SysAutoCode) error {
	if s.PackageName == "autocode" || s.PackageName == "system" || s.PackageName == "example" || s.PackageName == "" {
		return errors.New("cannot use reserved package name")
	}
	if !errors.Is(global.ZC_DB.Where("package_name = ?", s.PackageName).First(&system.SysAutoCode{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("same PackageName exists")
	}
	if e := autoCodeService.CreatePackageTemp(s.PackageName); e != nil {
		return e
	}
	return global.ZC_DB.Create(&s).Error
}

func (autoCodeService *AutoCodeService) GetPackage() (pkgList []system.SysAutoCode, err error) {
	err = global.ZC_DB.Find(&pkgList).Error
	return pkgList, err
}

func (autoCodeService *AutoCodeService) DelPackage(a system.SysAutoCode) error {
	return global.ZC_DB.Delete(&a).Error
}

func (autoCodeService *AutoCodeService) CreatePackageTemp(packageName string) error {
	Init(packageName)
	pendingTemp := []autoPackage{{
		path: packageService,
		name: packageServiceName,
		temp: string(subcontract.Server),
	}, {
		path: packageRouter,
		name: packageRouterName,
		temp: string(subcontract.Router),
	}, {
		path: packageAPI,
		name: packageAPIName,
		temp: string(subcontract.API),
	}}
	for i, s := range pendingTemp {
		pendingTemp[i].path = filepath.Join(global.ZC_CONFIG.AutoCode.Root, global.ZC_CONFIG.AutoCode.Server, filepath.Clean(fmt.Sprintf(s.path, packageName)))
	}
	// 选择模板
	for _, s := range pendingTemp {
		err := os.MkdirAll(filepath.Dir(s.path), 0755)
		if err != nil {
			return err
		}

		f, err := os.Create(s.path)
		if err != nil {
			return err
		}

		defer f.Close()

		temp, err := template.New("").Parse(s.temp)
		if err != nil {
			return err
		}
		err = temp.Execute(f, struct {
			PackageName string `json:"package_name"`
		}{packageName})
		if err != nil {
			return err
		}
	}

	for _, v := range pendingTemp {
		meta := packageInjectionMap[v.name]
		if err := ImportReference(meta.path, fmt.Sprintf(meta.importCodeF, v.name, packageName), fmt.Sprintf(meta.structNameF, caser.String(packageName)), fmt.Sprintf(meta.packageNameF, packageName), meta.groupName); err != nil {
			return err
		}
	}
	return nil
}

type Visitor struct {
	ImportCode  string
	StructName  string
	PackageName string
	GroupName   string
}

func (vi *Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:

		if n.Tok == token.IMPORT && vi.ImportCode != "" {
			vi.addImport(n)
			return nil
		}
		if n.Tok == token.TYPE && vi.StructName != "" && vi.PackageName != "" && vi.GroupName != "" {
			vi.addStruct(n)
			return nil
		}
	case *ast.FuncDecl:
		if n.Name.Name == "Routers" {
			vi.addFuncBodyVar(n)
			return nil
		}

	}
	return vi
}

func (vi *Visitor) addStruct(genDecl *ast.GenDecl) ast.Visitor {
	for i := range genDecl.Specs {
		switch n := genDecl.Specs[i].(type) {
		case *ast.TypeSpec:
			if strings.Index(n.Name.Name, "Group") > -1 {
				switch t := n.Type.(type) {
				case *ast.StructType:
					f := &ast.Field{
						Names: []*ast.Ident{
							{
								Name: vi.StructName,
								Obj: &ast.Object{
									Kind: ast.Var,
									Name: vi.StructName,
								},
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: vi.PackageName,
							},
							Sel: &ast.Ident{
								Name: vi.GroupName,
							},
						},
					}
					t.Fields.List = append(t.Fields.List, f)
				}
			}
		}
	}
	return vi
}

func (vi *Visitor) addImport(genDecl *ast.GenDecl) ast.Visitor {
	hasImported := false
	for _, v := range genDecl.Specs {
		importSpec := v.(*ast.ImportSpec)
		if importSpec.Path.Value == strconv.Quote(vi.ImportCode) {
			hasImported = true
		}
	}
	if !hasImported {
		genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(vi.ImportCode),
			},
		})
	}
	return vi
}

func (vi *Visitor) addFuncBodyVar(funDecl *ast.FuncDecl) ast.Visitor {
	hasVar := false
	for _, v := range funDecl.Body.List {
		switch varSpec := v.(type) {
		case *ast.AssignStmt:
			for i := range varSpec.Lhs {
				switch nn := varSpec.Lhs[i].(type) {
				case *ast.Ident:
					if nn.Name == vi.PackageName+"Router" {
						hasVar = true
					}
				}
			}
		}
	}
	if !hasVar {
		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: vi.PackageName + "Router",
					Obj: &ast.Object{
						Kind: ast.Var,
						Name: vi.PackageName + "Router",
					},
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "router",
						},
						Sel: &ast.Ident{
							Name: "RouterGroupApp",
						},
					},
					Sel: &ast.Ident{
						Name: caser.String(vi.PackageName),
					},
				},
			},
		}
		funDecl.Body.List = append(funDecl.Body.List, funDecl.Body.List[1])
		index := 1
		copy(funDecl.Body.List[index+1:], funDecl.Body.List[index:])
		funDecl.Body.List[index] = assignStmt
	}
	return vi
}

func ImportReference(filepath, importCode, structName, packageName, groupName string) error {
	fSet := token.NewFileSet()
	fParser, err := parser.ParseFile(fSet, filepath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	importCode = strings.TrimSpace(importCode)
	v := &Visitor{
		ImportCode:  importCode,
		StructName:  structName,
		PackageName: packageName,
		GroupName:   groupName,
	}
	if importCode == "" {
		ast.Print(fSet, fParser)
	}

	ast.Walk(v, fParser)

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fSet, fParser)
	if err != nil {
		log.Fatal(err)
	}
	return os.WriteFile(filepath, buffer.Bytes(), 0o600)
}

func (autoCodeService *AutoCodeService) CreatePlug(plug system.AutoPlugReq) error {
	plug.CheckList()
	tplFileList, _ := autoCodeService.GetAllTplFile(plugPath, nil)
	for _, tpl := range tplFileList {
		temp, err := template.ParseFiles(tpl)
		if err != nil {
			zap.L().Error("parse err", zap.String("tpl", tpl), zap.Error(err))
			return err
		}
		pathArr := strings.SplitAfter(tpl, "/")
		if strings.Index(pathArr[2], "tpl") < 0 {
			dirPath := filepath.Join(global.ZC_CONFIG.AutoCode.Root, global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SPlug, plug.Snake+"/"+pathArr[2]))
			os.MkdirAll(dirPath, 0755)
		}
		file := filepath.Join(global.ZC_CONFIG.AutoCode.Root, global.ZC_CONFIG.AutoCode.Server, fmt.Sprintf(global.ZC_CONFIG.AutoCode.SPlug, plug.Snake+"/"+tpl[len(plugPath):len(tpl)-4]))
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			zap.L().Error("open file", zap.String("tpl", tpl), zap.Error(err), zap.Any("plug", plug))
			return err
		}
		defer f.Close()

		err = temp.Execute(f, plug)
		if err != nil {
			zap.L().Error("exec err", zap.String("tpl", tpl), zap.Error(err), zap.Any("plug", plug))
			return err
		}
	}
	return nil
}

func (autoCodeService *AutoCodeService) InstallPlugin(file *multipart.FileHeader) (web, server int, err error) {
	const ZCPLUGPINATH = "./gva-plug-temp/"
	defer os.RemoveAll(ZCPLUGPINATH)
	_, err = os.Stat(ZCPLUGPINATH)
	if os.IsNotExist(err) {
		os.Mkdir(ZCPLUGPINATH, os.ModePerm)
	}

	src, err := file.Open()
	if err != nil {
		return -1, -1, err
	}
	defer src.Close()

	out, err := os.Create(ZCPLUGPINATH + file.Filename)
	if err != nil {
		return -1, -1, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	paths, err := utils.Unzip(ZCPLUGPINATH+file.Filename, ZCPLUGPINATH)
	paths = filterFile(paths)
	var webIndex = -1
	var serverIndex = -1
	for i := range paths {
		paths[i] = filepath.ToSlash(paths[i])
		pathArr := strings.Split(paths[i], "/")
		ln := len(pathArr)
		if ln < 2 {
			continue
		}
		if pathArr[ln-2] == "server" && pathArr[ln-1] == "plugin" {
			serverIndex = i
		}
		if pathArr[ln-2] == "web" && pathArr[ln-1] == "plugin" {
			webIndex = i
		}
	}
	if webIndex == -1 && serverIndex == -1 {
		zap.L().Error("non standard plugins, please automatically migrate and use according to the documentation")
		return webIndex, serverIndex, errors.New("non standard plugins, please automatically migrate and use according to the documentation")
	}

	if webIndex != -1 {
		err = installation(paths[webIndex], global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.Web)
		if err != nil {
			return webIndex, serverIndex, err
		}
	}

	if serverIndex != -1 {
		err = installation(paths[serverIndex], global.ZC_CONFIG.AutoCode.Server, global.ZC_CONFIG.AutoCode.Server)
	}
	return webIndex, serverIndex, err
}

func installation(path string, formPath string, toPath string) error {
	return nil
}

func filterFile(paths []string) []string {
	np := make([]string, 0, len(paths))
	for _, path := range paths {
		if ok, _ := skipMacSpecialDocument(path); ok {
			continue
		}
		np = append(np, path)
	}
	return np
}

func skipMacSpecialDocument(src string) (bool, error) {
	if strings.Contains(src, ".DS_Store") || strings.Contains(src, "__MACOSX") {
		return true, nil
	}
	return false, nil
}
