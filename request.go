package php2go

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	httpClient *http.Client
	once       sync.Once
)

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	DisableKeepAlives:     false,
}

// GetHttpClient 建议使用连接池复用 或 使用 github.com/go-resty/resty/v2
func GetHttpClient() *http.Client {
	once.Do(func() {
		httpClient = GetClient()
	})
	return httpClient
}

// GetClient 获取新的客户端
func GetClient() *http.Client {
	client := &http.Client{
		Transport: DefaultTransport,
		Timeout:   30 * time.Second,
	}
	return client
}

func HttpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
			log.Println("404,403:", url, resp.Status)
			return []byte(""), nil
		}
		return nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
func HttpDo(url string, param string, method string) ([]byte, error) {
	if method == "" {
		method = "POST"
	}
	req, err := http.NewRequest(method,
		url,
		strings.NewReader(param))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	defer req.Body.Close()
	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("http status error")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func HttpJson(url string, data interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// GetClientIp 获取客户端IP
func GetClientIp(r *http.Request) (ip string) {
	headers := []string{"X-Real-IP", "X-Forwarded-For"}
	for _, header := range headers {
		ip = r.Header.Get(header)
		if ip != "" {
			ip = strings.Split(ip, ",")[0]
			break
		}
	}
	if ip == "" {
		ip = r.RemoteAddr
		if strings.ContainsRune(r.RemoteAddr, ':') {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
	}
	if ip == "::1" || ip == "" {
		ip = "127.0.0.1"
	}
	return ip
}
