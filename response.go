package hrq

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

// Response wraps http.Response.
type Response struct {
	Res  *http.Response
	Body []byte
}

// Content returns response body by byte.
func (r *Response) Content() ([]byte, error) {
	if r.Body != nil {
		return r.Body, nil
	}
	defer r.Res.Body.Close()
	bs, err := ioutil.ReadAll(r.Res.Body)
	if err != nil {
		return nil, err
	}
	r.Body = bs
	return bs, err
}

// ContentType returns content-type in response header..
func (r *Response) ContentType() string {
	for k, v := range r.Res.Header {
		if strings.ToLower(k) == "content-type" {
			return v[0]
		}
	}
	return ""
}

// Encode returns encode of response body.
func (r *Response) Encode() (encode string, err error) {
	contentType := r.ContentType()
	body, err := r.Content()
	if err != nil {
		return
	}
	_, encode, _ = charset.DetermineEncoding(body, contentType)
	return
}

// Text returns response body by string.
func (r *Response) Text() (text string, err error) {
	encode, err := r.Encode()
	if err != nil {
		return
	}
	content, err := r.Content()
	if err != nil {
		return
	}
	br := bytes.NewReader(content)
	rl, err := charset.NewReaderLabel(encode, br)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(rl)
	if err != nil {
		return
	}
	text = string(bs)
	return
}
