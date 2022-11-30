package httpClient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type response struct {
	Status     string
	StatusCode int
	Error      error
	TotalTime  time.Duration
	Body       []byte
}

type request struct {
	Method  string
	Url     string
	Body    []byte
	Headers []map[string]string
}

var httpTimeout = 20 * time.Second

func newRequest(request request) (*http.Request, error) {
	req, err := http.NewRequest(request.Method, request.Url, bytes.NewReader(request.Body))

	// 设置header
	for _, header := range request.Headers {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	return req, err
}

func httpResponse(req *http.Request, responseErr error) response {
	var res response

	if responseErr != nil {
		res.Error = responseErr
		return res
	}

	startT := time.Now()
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)

	res.TotalTime = time.Now().Sub(startT)
	if err != nil { // 这里报错说明没有响应内容
		res.Error = err
		return res
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	res.Status = resp.Status
	res.StatusCode = resp.StatusCode
	res.Body, res.Error = ioutil.ReadAll(resp.Body)
	return res
}

// SetTimeout 	设置请求超时时间.
// second		请求超时时间.
func SetTimeout(second time.Duration) bool {
	httpTimeout = second * time.Second
	return true
}

// GetTimeout 	查看请求超时时间.
func GetTimeout() time.Duration {
	return httpTimeout
}

// GET 发送GET请求.
// url：		请求地址.
// headers：	设置请求头.
func GET(url string, headers ...map[string]string) response {
	return httpResponse(newRequest(request{
		Method:  "GET",
		Url:     url,
		Body:    nil,
		Headers: headers,
	}))
}

// POST 发送POST请求.
// url：		请求地址.
// data：		POST请求提交的数据.
// headers：	设置请求头，默认Content-Type：application/json;charset=UTF-8.
func POST(url string, data []byte, headers ...map[string]string) response {
	if len(headers) == 0 {
		headers = append(headers, map[string]string{"Content-Type": "application/json;charset=UTF-8"})
	}

	return httpResponse(newRequest(request{
		Method:  "POST",
		Url:     url,
		Body:    data,
		Headers: headers,
	}))
}
