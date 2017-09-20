package hrq

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"time"
)

const applicationFormUrlencoded = "application/x-www-form-urlencoded"
const applicationJSON = "application/json"
const multipartFormData = "multipart/form-data"

// DefaultTimeout is seconds of timeout.
var DefaultTimeout = 15

// DefaultContentType is a default content-type of request.
var DefaultContentType = applicationFormUrlencoded

// File is file for multipart/form.
type File struct {
	ContentType string
	FieldName   string
	Name        string
	File        *os.File
}

// Request inherits http.Request.
type Request struct {
	*http.Request
	Timeout time.Duration
	Data    interface{}
	Files   []*File
	// Gzip is a flag to decide whether to compress by gzip or not.
	// (defaut false)
	Gzip bool
}

func (r *Request) contentType() string {
	return r.HeaderValue("Content-Type")
}

func (r *Request) isPostOrPut() bool {
	return r.Method == "POST" || r.Method == "PUT"
}

func (r *Request) setBody(b []byte) {
	if r.Gzip {
		var buffer bytes.Buffer
		writer := gzip.NewWriter(&buffer)
		writer.Write(b)
		writer.Close()
		b = buffer.Bytes()
	}
	r.GetBody = func() (io.ReadCloser, error) {
		reader := bytes.NewReader(b)
		return ioutil.NopCloser(reader), nil
	}
	reader := bytes.NewReader(b)
	r.Body = ioutil.NopCloser(reader)
}

// UseGzip makes Request.Gzip to true.
func (r *Request) UseGzip() *Request {
	r.Gzip = true
	return r
}

// AddFile sets file for multipart/form.
func (r *Request) AddFile(contentType, fieldName, fileName string, file *os.File) *Request {
	f := &File{
		ContentType: contentType,
		FieldName:   fieldName,
		Name:        fileName,
		File:        file,
	}
	r.Files = append(r.Files, f)
	return r
}

// SetTimeout sets timeout.
func (r *Request) SetTimeout(timeout int) *Request {
	r.Timeout = time.Duration(timeout) * time.Second
	return r
}

// SetApplicationFormUrlencoded is an alias of req.SetHeader("Content-Type", "application/x-www-form-urlencoded").
func (r *Request) SetApplicationFormUrlencoded() *Request {
	return r.SetHeader("Content-Type", applicationFormUrlencoded)
}

// SetApplicationJSON is an alias of req.SetHeader("Content-Type", "application/json").
func (r *Request) SetApplicationJSON() *Request {
	return r.SetHeader("Content-Type", applicationJSON)
}

// SetMultipartFormData is an alias of req.SetHeader("Content-Type", "multipart/form-data").
func (r *Request) SetMultipartFormData() *Request {
	return r.SetHeader("Content-Type", multipartFormData)
}

// Send sends request.
// If method is POST and content-type is application/x-www-form-urlencoded,
// the request data is urlencoded.
// If method is POST and content-type is application/json,
// the request data is converted to json string.
func (r *Request) Send() (res *Response, err error) {
	if r.isPostOrPut() && r.Data != nil && r.HeaderValue("Content-Type") != multipartFormData {
		if r.contentType() == applicationFormUrlencoded {
			data, ok := r.Data.(map[string]string)
			if !ok {
				err := errors.New("data is not a map[string]string at Request.Send()")
				return nil, err
			}
			mapStringList := mapStringList(data)
			values := []byte(url.Values(mapStringList).Encode())
			r.setBody(values)
		} else if r.contentType() == applicationJSON {
			jsonBytes, err := json.Marshal(r.Data)
			if err != nil {
				return nil, err
			}
			r.setBody(jsonBytes)
		}
	} else if r.isPostOrPut() && r.contentType() == multipartFormData {
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		data, ok := r.Data.(map[string]string)
		if !ok {
			err := errors.New("data is not a map[string]string at Request.Send()")
			return nil, err
		}
		for k, v := range data {
			writer.WriteField(k, v)
		}
		for _, file := range r.Files {
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", file.ContentType)
			desc := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, file.FieldName, file.Name)
			part.Set("Content-Disposition", desc)
			fileWriter, err := writer.CreatePart(part)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fileWriter, file.File)
			if err != nil {
				return nil, err
			}
			defer file.File.Close()
		}
		writer.Close()
		r.SetHeader("Content-Type", writer.FormDataContentType())
		r.ContentLength = int64(buffer.Len())
		b := buffer.Bytes()
		r.setBody(b)
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	cli := &http.Client{
		Timeout: r.Timeout,
		Jar:     jar,
	}
	requestHistory := []*http.Request{}
	cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		requestHistory = via
		if len(via) >= 10 {
			return errors.New("there are 10 redirects")
		}
		return nil
	}
	if r.Gzip {
		r.SetHeader("Content-Encoding", "gzip")
	}
	response, err := cli.Do(r.Request)
	if err != nil {
		return
	}
	res = &Response{
		Response: response,
		History:  requestHistory,
	}
	return
}

// SetHeader sets a value of request header.
func (r *Request) SetHeader(key, value string) *Request {
	r.Header.Set(key, value)
	return r
}

// HeaderValue returns a value of request header.
func (r *Request) HeaderValue(key string) string {
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
		Gzip:    false,
	}
	return
}

// Get make a request whose method is GET.
func Get(url string) (req *Request, err error) {
	req, err = NewRequest("GET", url, nil, DefaultTimeout)
	return
}

// Delete make a request whose method is DELETE.
func Delete(url string) (req *Request, err error) {
	req, err = NewRequest("DELETE", url, nil, DefaultTimeout)
	return
}

// Head make a request whose method is HEAD.
func Head(url string) (req *Request, err error) {
	req, err = NewRequest("HEAD", url, nil, DefaultTimeout)
	return
}

// Options make a request whose method is OPTIONS.
func Options(url string) (req *Request, err error) {
	req, err = NewRequest("OPTIONS", url, nil, DefaultTimeout)
	return
}

func postOrPut(method, url string, data interface{}) (req *Request, err error) {
	req, err = NewRequest(method, url, nil, DefaultTimeout)
	if err != nil {
		return
	}
	req.SetHeader("Content-Type", DefaultContentType)
	req.Files = []*File{}
	req.Data = data
	return
}

// Post make a request whose method is POST.
func Post(url string, data interface{}) (req *Request, err error) {
	req, err = postOrPut("POST", url, data)
	return
}

// Put make a request whose method is PUT.
func Put(url string, data interface{}) (req *Request, err error) {
	req, err = postOrPut("PUT", url, data)
	return
}
