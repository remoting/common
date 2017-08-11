package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpResponse struct {
	Data   []byte
	Status int
	Header map[string]string
}

var HTTP_TIMEOUT int = 30

func getHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(HTTP_TIMEOUT-10))
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(HTTP_TIMEOUT)))
				return conn, nil
			},
		},
	}
	return client
}
func HttpRequest(method, url string, header *map[string]string, _body io.Reader) (*HttpResponse, error) {
	client := getHttpClient()
	req, err := http.NewRequest(method, url, _body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http error:url=%s,error=%s\n", url, err.Error())
		return nil, err //handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http body error:url=%s,error=%s\n", url, err.Error())
		return nil, err //handle error
	}
	defer resp.Body.Close()
	_resp := new(HttpResponse)
	_resp.Status = resp.StatusCode
	_resp.Data = body
	_resp.Header = make(map[string]string, 0)
	for k, _ := range resp.Header {
		_resp.Header[k] = resp.Header.Get(k)
	}
	return _resp, nil
}
func getParamString(params *map[string]string) string {
	q, _ := url.ParseQuery("")
	if params != nil {
		for name, value := range *params {
			q.Add(name, value)
		}
	}
	return q.Encode()
}
func getUrlParams(_url string, params *map[string]string) string {
	if strings.Index(_url, "?") >= 0 {
		return _url + "&" + getParamString(params)
	} else {
		return _url + "?" + getParamString(params)
	}
}
func initHeader(header *map[string]string) map[string]string {
	if header == nil {
		return make(map[string]string, 0)
	}
	return *header
}
func GetJsonHeader() map[string]string {
	header := make(map[string]string, 0)
	header["Content-Type"] = "application/json"
	return header
}
func HttpDelete(url string, params, header *map[string]string) (*HttpResponse, error) {
	_header := initHeader(header)
	_, ok := _header["Content-Type"]
	if !ok {
		_header["Content-Type"] = "application/json"
	}
	return HttpRequest("DELETE", getUrlParams(url, params), &_header, nil)
}
func HttpGet(url string, params, header *map[string]string) (*HttpResponse, error) {
	_header := initHeader(header)
	_, ok := _header["Content-Type"]
	if !ok {
		_header["Content-Type"] = "application/json"
	}
	return HttpRequest("GET", getUrlParams(url, params), &_header, nil)
}
func HttpPost(url string, params, header *map[string]string) (*HttpResponse, error) {
	_header := initHeader(header)
	_, ok := _header["Content-Type"]
	if !ok {
		_header["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return HttpRequest("POST", url, &_header, strings.NewReader(getParamString(params)))
}
func HttpPut(url string, params, header *map[string]string) (*HttpResponse, error) {
	_header := initHeader(header)
	_, ok := _header["Content-Type"]
	if !ok {
		_header["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return HttpRequest("POST", url, &_header, strings.NewReader(getParamString(params)))
}
