package hrq

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultTimeout is seconds of timeout.
var DefaultTimeout = 15

// DefaultContentType is a content-type of request.
var DefaultContentType = "application/x-www-form-urlencoded"

// Request wraps http.Request.
type Request struct {
	Req     *http.Request
	Timeout time.Duration
	Data    map[string][]string
}

func (r *Request) setBody(values *strings.Reader) {
	body := ioutil.NopCloser(values)
	r.Req.Body = body
	r.Req.ContentLength = int64(values.Len())
	s := *values
	r.Req.GetBody = func() (io.ReadCloser, error) {
		r := s
		return ioutil.NopCloser(&r), nil
	}
}

// Send sends request.
func (r *Request) Send() (res *Response, err error) {
	if r.Req.Method == "POST" && r.Data != nil {
		if r.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
			values := strings.NewReader(url.Values(r.Data).Encode())
			r.setBody(values)
		} else if r.GetHeader("Content-Type") == "application/json" {
			jsonBytes, err := json.Marshal(r.Data)
			if err != nil {
				return nil, err
			}
			values := strings.NewReader(string(jsonBytes))
			r.setBody(values)
		}
	}
	cli := &http.Client{
		Timeout: r.Timeout,
	}
	response, err := cli.Do(r.Req)
	if err != nil {
		return
	}
	res = &Response{
		Res: response,
	}
	return
}

// SetHeader sets a value of request header.
func (r *Request) SetHeader(key, value string) *Request {
	r.Req.Header.Set(key, value)
	return r
}

// GetHeader returns a value of request header.
func (r *Request) GetHeader(key string) string {
	return r.Req.Header.Get(key)
}

// DelHeader delete a value of request header by key.
func (r *Request) DelHeader(key string) *Request {
	r.Req.Header.Del(key)
	return r
}

// NewRequest make a Request.
func NewRequest(method, url string, body io.Reader, timeoutSecond int) (req *Request, err error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	timeout := time.Duration(timeoutSecond) * time.Second
	req = &Request{
		Req:     request,
		Timeout: timeout,
	}
	return
}

// Get make a request whose method is GET.
func Get(url string) (req *Request, err error) {
	req, err = NewRequest("GET", url, nil, DefaultTimeout)
	return
}

// Post make a request whose method is GET.
func Post(url string, data map[string][]string) (req *Request, err error) {
	req, err = NewRequest("Post", url, nil, DefaultTimeout)
	if err != nil {
		return
	}
	req.SetHeader("Content-Type", DefaultContentType)
	req.Data = data
	return
}
