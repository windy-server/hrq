package hrq

import (
	"io"
	"net/http"
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
}

// Send sends request.
func (r *Request) Send() (res *Response, err error) {
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
func Post(url string) (req *Request, err error) {
	req, err = NewRequest("Post", url, nil, DefaultTimeout)
	req.SetHeader("Content-Type", DefaultContentType)
	return
}
