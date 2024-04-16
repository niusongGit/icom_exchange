package http_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetRequest 发送HTTP GET请求并返回响应体
func GetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type PostResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func PostReques(url, rpcUser, rpcPwd string, params map[string]interface{}) (interface{}, error) {
	header := http.Header{
		"user":     []string{rpcUser},
		"password": []string{rpcPwd},
	}
	client := &http.Client{Timeout: time.Second * 6}
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bs)))
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error() + "Offline")
	}
	defer resp.Body.Close()

	var resultBs []byte

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	resultBs, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(PostResult)
	buf := bytes.NewBuffer(resultBs)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}

	if result.Code != 2000 {
		return nil, errors.New(result.Message)
	}
	return result.Result, nil
}
