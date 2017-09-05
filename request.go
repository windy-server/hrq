package hrq

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultTimeout is seconds of timeout.
var DefaultTimeout = 15

// DefaultContentType is a default content-type of request.
var DefaultContentType = "application/x-www-form-urlencoded"

// Request inherits http.Request.
type Request struct {
	*http.Request
	Timeout time.Duration
	Data    interface{}
}

func (r *Request) setBody(values *strings.Reader) {
	body := ioutil.NopCloser(values)
	r.Body = body
	r.ContentLength = int64(values.Len())
	s := *values
	r.GetBody = func() (io.ReadCloser, error) {
		r := s
		return ioutil.NopCloser(&r), nil
	}
}

// SetTimeout sets timeout.
func (r *Request) SetTimeout(timeout int) *Request {
	r.Timeout = time.Duration(timeout) * time.Second
	return r
}

// Send sends request.
// If method is POST and content-type is application/x-www-form-urlencoded,
// the request data is urlencoded.
// If method is POST and content-type is application/json,
// the request data is converted to json string.
func (r *Request) Send() (res *Response, err error) {
	if r.Method == "POST" && r.Data != nil {
		if r.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
			data, ok := r.Data.(map[string][]string)
			if !ok {
				err := errors.New("data is not a map[string][]string at Request.Send()")
				return nil, err
			}
			values := strings.NewReader(url.Values(data).Encode())
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
	// else if r.Method == "POST" && r.GetHeader("Content-Type") == "multipart/form-data" {
	// 	var buffer bytes.Buffer
	// 	writer := multipart.NewWriter(&buffer)
	// 	data, ok := r.Data.(map[string]string)
	// 	if !ok {
	// 		err := errors.New("data is not a map[string]string at Request.Send()")
	// 		return nil, err
	// 	}
	// 	for k, v := range data {
	// 		writer.WriteField(k, v)
	// 	}
	// }

	cli := &http.Client{
		Timeout: r.Timeout,
	}
	response, err := cli.Do(r.Request)
	if err != nil {
		return
	}
	res = &Response{
		Response: response,
	}
	return
}

// SetHeader sets a value of request header.
func (r *Request) SetHeader(key, value string) *Request {
	r.Header.Set(key, value)
	return r
}

// GetHeader returns a value of request header.
func (r *Request) GetHeader(key string) string {
	return r.Header.Get(key)
}

// DelHeader delete a value of request header by key.
func (r *Request) DelHeader(key string) *Request {
	r.Header.Del(key)
	return r
}

// PutCookie makes a cookie which is setted name and value.
// It adds a cookie to a request.
func (r *Request) PutCookie(name, value string) *Request {
	r.AddCookie(&http.Cookie{Name: name, Value: value})
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
		Request: request,
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
func Post(url string, data interface{}) (req *Request, err error) {
	req, err = NewRequest("POST", url, nil, DefaultTimeout)
	if err != nil {
		return
	}
	req.SetHeader("Content-Type", DefaultContentType)
	req.Data = data
	return
}
