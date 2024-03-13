package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//@function: MD5V
//@description: md5 encryption
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// GetFileMD5 Get file Md5
func GetFileMD5(pathName string) string {
	f, err := os.Open(pathName)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}
	has := md5hash.Sum(nil)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
