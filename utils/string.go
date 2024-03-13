package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io/ioutil"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ZERO  = 0
	ONE   = 1
	TWO   = 2
	THREE = 3
	FOUR  = 4
	Five  = 5
)

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(e, 10, 64)
}

func IntToString(e int) string {
	return strconv.Itoa(e)
}

func Float64ToString(e float64) string {
	return strconv.FormatFloat(e, 'f', -1, 64)
}

func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}

func IdsStrToIdsInt64Group(key string, c *gin.Context) []int64 {
	IDS := make([]int64, 0)
	ids := strings.Split(c.Request.FormValue(key), ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := strconv.ParseInt(ids[i], 10, 64)
		IDS = append(IDS, ID)
	}
	return IDS
}

func StructToJsonStr(e interface{}) (string, error) {
	if b, err := json.Marshal(e); err == nil {
		return string(b), err
	} else {
		return "", err
	}
}

func GetBodyString(c *gin.Context) (string, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return string(body), nil
	} else {
		return "", err
	}
}

func JsonStrToMap(e string) (map[string]interface{}, error) {
	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(e), &dict); err == nil {
		return dict, err
	} else {
		return nil, err
	}
}

func LimitAndOrderBy(param map[string]string) string {
	str := ""
	//排序
	if param["sort"] != "" {
		str += ` order by ` + param["sort"]
		if param["order"] != "" {
			str += " " + param["order"]
		}
	}
	if param["isexp"] == "" {
		param["isexp"] = "0"
	}
	if param["isexp"] != "1" {

		if param["page"] != "" && param["pageSize"] != "" {
			pageNum, _ := strconv.Atoi(param["page"])
			pageSize, _ := strconv.Atoi(param["pageSize"])
			if pageNum != 0 && pageSize != 0 {
				str += ` limit ` + strconv.Itoa((pageNum-1)*pageSize) + `,` + param["pageSize"]
			}
		}

	}

	return str
}

func LimitAndOrder(sort, order string, page, pageSize int) string {
	str := ""

	if sort != "" {
		str += ` order by ` + sort
		if order != "" {
			str += " " + order
		}
	}
	if pageSize > 0 {
		str += ` limit ` + strconv.Itoa((page-1)*pageSize) + `,` + strconv.Itoa(pageSize)
	}
	return str
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func GetRandString(num int) string {
	s, _ := rand.Prime(rand.Reader, 32)
	str := strconv.FormatInt(time.Now().UnixNano(), 10) + s.String()
	sh := sha1.New()
	sh.Write([]byte(str))
	res := sh.Sum(nil)
	return hex.EncodeToString(res)[0:num]
}

func BeginToEndTimestampstr(param map[string]interface{}) (con string) {
	field := ""
	if param["field"] == nil {
		return
	} else {
		field = param["field"].(string)
	}

	if param["beginTime"] != nil {
		beginTime := param["beginTime"].(string)
		if beginTime != "" {
			t1, _ := DateToTimeStamp(beginTime+" 00:00:00", "2006-01-02 15:04:05")
			con += " and " + field + ">= " + strconv.FormatInt(t1, 10)
		}
	}
	if param["endTime"] != nil {
		endTime := param["endTime"].(string)
		if endTime != "" {
			t2, _ := DateToTimeStamp(endTime+" 23:59:59", "2006-01-02 15:04:05")
			con += " and " + field + "<= " + strconv.FormatInt(t2, 10)
		}
	}
	return con
}

func BeginToEndDatestr(param map[string]interface{}) (con string) {
	field := ""
	if param["field"] == nil {
		return
	} else {
		field = param["field"].(string)
	}

	if param["beginTime"] != nil {
		beginTime := param["beginTime"].(string)
		if beginTime != "" {
			con += " and " + field + ">= '" + beginTime + "'"
		}
	}
	if param["endTime"] != nil {
		endTime := param["endTime"].(string)
		if endTime != "" {
			con += " and " + field + "<= '" + endTime + " 23:59:59'"
		}
	}
	return con
}

func SqlReplaceParames(sql string, param map[string]string) string {
	fa := false
	start := 0
	sqlstr := sql
	fl := true
	for i, v := range sql {
		if v == ':' {
			start = i + 1
			fa = true
		}
		if (v == '\n' || v == '\t' || v == ' ' || v == ',' || v == ')' || v == '%' || v == '"' || v == '=' || len(sql)-1 == i) && fa && fl {
			field := sql[start:i]
			//最后一个
			if len(sql)-1 == i && v != ' ' && v != '\n' && v != '\t' && v != ')' {
				field = sql[start : i+1]
			}
			if param[field] != "" {
				if sql[start-3] == '%' {
					sqlstr = strings.Replace(sqlstr, "%%:"+field+"%%", `'%%`+param[field]+`%%'`, 1)
				} else if sql[start-2] == '%' {
					sqlstr = strings.Replace(sqlstr, "%:"+field+"%", `'%`+param[field]+`%'`, 1)
				} else {
					flen := len(field)
					//避免包含于字段 如 :bank  :banknum
					if len(sql) > start+flen {
						v1 := sql[start+flen]
						//v2 := ','
						if v1 == '\n' || v1 == '\t' || v1 == ' ' || v1 == ',' || v1 == ')' || v1 == '%' || v1 == '"' || v1 == '=' {
							//fmt.Println(field,v1)
							sqlstr = strings.Replace(sqlstr, ":"+field, `'`+param[field]+`'`, 1)
							fa = false
						}
						continue
					}
					sqlstr = strings.Replace(sqlstr, ":"+field, `'`+param[field]+`'`, -1)
				}
				fa = false
			} else {
				if _, ok := param[field]; ok {
					sqlstr = strings.Replace(sqlstr, ":"+field, `'`+param[field]+`'`, 1)
					fa = false
					continue
				}
				if sql[i-1] == '\'' || sql[i-1] == '"' {
					fa = false
					continue
				}
				sqlstr = field + " Parameter does not exist!"
				return sqlstr
			}
		}

	}
	return sqlstr
}

func ExecScript(script string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", script)
	b, err := cmd.CombinedOutput()
	return b, err
}

func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}

func GrpcConnect(ip, port string) (co *grpc.ClientConn, err error) {
	conn, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	if err != nil {
		return conn, errors.New(fmt.Sprintf("grpc.Dial Connection failed:%s:%s ", ip, port))
	}
	return conn, nil
}

func IsStringsUnique(data []string) bool {
	tempMap := make(map[string]bool)
	for _, val := range data {
		_, ok := tempMap[val]
		if !ok {
			tempMap[val] = true
		}
	}
	if len(tempMap) == len(data) {
		return true
	}
	return false
}

func IsNull(str string) bool {
	return str == "" || str == "undefined" || str == " "
}

func FileCoinStrToUint64(str string) (uint64, error) {

	addr, err := address.NewFromString(str)
	if err != nil {
		return 0, err
	}
	miner, err := address.IDFromAddress(addr)
	if err != nil {
		return 0, err
	}

	return miner, nil
}

func Replace(str string) string {

	if str == "" {
		return ""
	}

	strVal := strings.Replace(str, "\n", "", -1)
	strVal = strings.Replace(strVal, " ", "", -1)
	return strVal
}

func TrimBlankSpace(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func FormDataStrToArray(str string) []string {
	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	str = strings.Replace(str, `"`, "", -1)
	return strings.Split(str, ",")
}

func GetNumFromStr(str string) string {
	if str == "" {
		return ""
	}
	str = strings.TrimSpace(str)
	reg := regexp.MustCompile(`[\s|_|a-z|A-Z]{1,}`)
	return reg.ReplaceAllString(str, "")
}

func GetLetterFromStr(str string) string {
	if str == "" {
		return ""
	}
	str = strings.TrimSpace(str)
	reg := regexp.MustCompile(`[\W|_|0-9]{1,}`)
	return reg.ReplaceAllString(str, "")
}

func FilterIP(str string) string {
	reg := regexp.MustCompile(`([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(.([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])){3}`)
	return reg.FindString(str)
}

func IsNumber(str string) bool {
	b, err := regexp.MatchString("^[0-9]+$", str)
	if err != nil {
		return false
	}
	if !b {
		return false
	}
	return true
}

func IsInStrList(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	//index的取值：[0,len(strArray)]
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

func KeepNumAndPoint(data string) string {
	if len(data) == 0 {
		return ""
	}
	r := ""
	for i := 0; i < len(data); i++ {
		if IsAlnum(data[i]) {
			r = r + string(data[i])
		}
	}
	return r
}

func IsAlnum(b byte) bool {
	return (b == '.') || (b >= '0' && b <= '9')
}
