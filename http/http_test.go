package http

import "testing"
import "fmt"
import "strings"

func TestHttp_1(t *testing.T) {
	//错误代码示范
	var header *map[string]string
	if *header == nil {
		fmt.Printf("%v", header)
	}
}
func TestHttp001(t *testing.T) {
	header := make(map[string]string, 1)
	params := make(map[string]string, 1)
	params["a"] = "a"
	resp, err := HttpPost("http://www.baidu.com", &params, &header)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Printf("%s\n", string(resp.Data)[0:10])
	fmt.Printf("%s\n", header["Content-Type"])
}

func TestHttp002(t *testing.T) {
	params := make(map[string]string, 1)
	params["中文"] = "中文"
	params["e"] = "e"
	_url1 := getUrlParams("http://www.baidu.com?a=b", &params)
	fmt.Printf("%s\n", _url1)
	_url2 := getUrlParams("http://www.baidu.com", &params)
	fmt.Printf("%s\n", _url2)
}

func TestHttp003(t *testing.T) {
	header := make(map[string]string, 1)
	header["Content-Type"] = "application/json"
	params := make(map[string]string, 1)
	resp, err := HttpDelete("http://www.baidu.com", &params, &header)
	if err != nil {
		t.Errorf("%+v", err)
	} else if strings.Trim(string(resp.Data), "")[0:10] != "<!DOCTYPE " {
		t.Errorf("error")
	}
}
