package util

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
)

const TIMEOUT = 220

func Get(url string, headerData map[string]string) (string, error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if response == nil {
		return "", nil
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Post(url string, params []byte, headerData map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func PostV2(url string) (string, error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(body), nil
}

func PostSSE(url string, params []byte, headerData map[string]string, channel chan string) (res string, err error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	if err != nil {
		log.Println(err)
		return
	}
	if len(headerData) > 0 {
		for k, v := range headerData {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	if resp == nil {
		log.Println("Received empty response")
		err = fmt.Errorf("received empty response")
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK response: %d\n", resp.StatusCode)
		err = fmt.Errorf("received non-OK response: %d", resp.StatusCode)
		return
	}
	var completeMessage strings.Builder
	// 使用 bufio.Scanner 逐行读取 SSE 数据
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// 忽略注释行（以 ":" 开头的行）
		if strings.HasPrefix(line, ":") {
			continue
		}
		// 解析事件类型和数据
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			//fmt.Println("Received data:", data)
			channel <- data
			message := gjson.Get(data, "message").String()
			completeMessage.WriteString(message)
		}
	}
	res = completeMessage.String()
	// 检查是否发生读取错误
	if err1 := scanner.Err(); err1 != nil {
		fmt.Printf("Error reading SSE stream: %v\n", err1)
		err = fmt.Errorf("error reading SSE stream: %v", err1)
		return
	}
	return
}

func GetV2(url string) (string, error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logx.Info(err)
		return "", err
	}
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		logx.Info(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Info(err)
		return "", err
	}
	return string(body), nil
}

func SetHeader(w http.ResponseWriter, filename string) {
	// 获取文件扩展名
	ext := filepath.Ext(filename)
	// 根据扩展名设置ContentType
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}

func GetV3(url string, params []byte, headerData map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: TIMEOUT * time.Second}

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	if len(headerData) > 0 {
		for k, v := range headerData {
			request.Header.Set(k, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}

	defer response.Body.Close()
	// 检查HTTP状态码
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// PostV3 发送HTTP请求的通用方法
func PostV3(ctx context.Context, url string, method string, requestBody interface{}, responseBody interface{}) error {
	logger := logx.WithContext(ctx)
	// 将请求体序列化为JSON
	var jsonData []byte
	var err error

	if requestBody != nil {
		jsonData, err = json.Marshal(requestBody)
		if err != nil {
			logger.Errorf("failed to marshal request body: %+v", err)
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Errorf("failed to create request: %w", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 设置超时时间
	client := &http.Client{
		Timeout: 120 * time.Second, // 设置为30秒或更长
	}
	resp, err := client.Do(req)
	if err != nil {
		// 记录详细的错误信息
		logger.Errorf("failed to send request to %s: %+v", url, err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("failed to read response body: %w", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}
	logger.Infof("url: %s,request: %s,response:%+v, body: %s", url, string(jsonData), resp, string(body))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("request failed with status code: %d, response: %s", resp.StatusCode, string(body))
		return fmt.Errorf("request failed with status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	if responseBody != nil {
		if err := json.Unmarshal(body, responseBody); err != nil {
			logger.Errorf("failed to unmarshal response body: %w", err)
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}
