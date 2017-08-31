package hrq

import (
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
	bytes, err := ioutil.ReadAll(r.Res.Body)
	if err != nil {
		return nil, err
	}
	r.Body = bytes
	return bytes, err
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
