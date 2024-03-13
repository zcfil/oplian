package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	AttoFIL = 18
	NanoFIL = 9
)

func NanoOrAttoToFIL(fil string, filtype int) (res float64, err error) {
	//大于18or9位
	if len(fil) > filtype {
		str := fil[0:len(fil)-filtype] + "." + fil[len(fil)-filtype:]
		res, err = strconv.ParseFloat(str, 64)
		return
	}
	//小于18or9位
	str := "0."
	for i := 0; i < filtype-len(fil); i++ {
		str += "0"
	}
	str = str + fil
	res, err = strconv.ParseFloat(str, 64)
	return
}

func BlockHeight() uint64 {
	//dataStr := "2022-11-02 02:13:00" //校准网重置时间
	dataStr := "2020-08-25 06:00:00" //正式网
	t, _ := time.ParseInLocation("2006-01-02 15:04:05 ", dataStr, time.Local)
	t1 := time.Now().UnixNano() / 1e6
	t2 := t.UnixNano() / 1e6
	num := (t1 - t2) / 30 / 1000

	return uint64(num)
}

func BlockHeightToTime(num int64) time.Time {

	num = num * 30 * 1e3
	dataStr := "2020-08-25 06:00:00"
	t, _ := time.ParseInLocation("2006-01-02 15:04:05 ", dataStr, time.Local)
	t2 := t.UnixNano() / 1e6
	t1 := (t2 + num) / 1e3

	return time.Unix(t1, 0)
}

func SectorNumString(miner string, number uint64) string {
	miner = strings.Replace(miner, "f0", "", 1)
	miner = strings.Replace(miner, "t0", "", 1)
	return fmt.Sprintf("s-t0%s-%d", miner, number)
}

func CheckSectorNum(buf string) bool {
	reg := regexp.MustCompile(`s-t0\d*-\d*`)
	result := reg.FindString(buf)
	if result == "" || result != buf {
		return false
	}
	return true
}

var byteSizeUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB"}

func SizeStr(bi string) string {
	b, _ := new(big.Int).SetString(bi, 10)
	r := new(big.Rat).SetInt(b)
	den := big.NewRat(1, 1024)

	var i int
	for f, _ := r.Float64(); f >= 1024 && i+1 < len(byteSizeUnits); f, _ = r.Float64() {
		i++
		r = r.Mul(r, den)
	}

	f, _ := r.Float64()
	return fmt.Sprintf("%.4g %s", f, byteSizeUnits[i])
}

func FMinerID(miner string) string {
	miner = strings.Replace(miner, "f0", "", 1)
	miner = strings.Replace(miner, "t0", "", 1)
	return fmt.Sprintf("f0%s", miner)
}
func StringToSectorID(sn string) (actor uint64, number uint64) {
	sn = strings.ReplaceAll(sn, "s-t0", "")
	strs := strings.Split(sn, "-")
	if len(strs) < 2 {
		return 0, 0
	}

	actor, _ = strconv.ParseUint(strs[0], 10, 64)
	number, _ = strconv.ParseUint(strs[1], 10, 64)
	return actor, number
}
func MinerActorID(miner string) uint64 {
	miner = strings.Replace(miner, "f0", "", 1)
	miner = strings.Replace(miner, "t0", "", 1)
	actor, _ := strconv.ParseUint(miner, 10, 64)
	return actor
}
