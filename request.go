package hrq

import (
	"io"
	"net/http"
	"time"
)

// DefaultTimeout is seconds of timeout.
var DefaultTimeout = 15

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
