package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RequestDo Encapsulate HTTP requests
func RequestDo(url, router, token string, request interface{}, timeout time.Duration) ([]byte, error) {

	client := &http.Client{
		Timeout: timeout,
	}
	var buf []byte
	if request != nil {
		var err error
		buf, err = json.Marshal(request)
		if err != nil {
			log.Println("Serialization failure：", err.Error())
			return nil, err
		}
	}

	path := "http://" + url + router
	req, err := http.NewRequest("POST", path, bytes.NewReader(buf))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error：", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Request error result：", *resp)
		return nil, fmt.Errorf("请求错误代码：%d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func LotusApiInfoDecompose(apiInfo string) (token, ip, port string) {
	strs := strings.Split(apiInfo, ":")
	if len(strs) <= 1 {
		return "", "", ""
	}
	token = strs[0]
	strs = strings.Split(strs[1], "/")
	if len(strs) <= 5 {
		return "", "", ""
	}
	ip = strs[2]
	port = strs[4]
	return
}

func LotusApiInfoMerge(token, ip, port string) (apiInfo string) {
	return fmt.Sprintf("%s:/ip4/%s/tcp/%s/http", token, ip, port)
}
