package hrq

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html/charset"
)

// Response inherits http.Response.
type Response struct {
	*http.Response
	// History is the redirect history.
	History []*http.Request
	rawBody []byte
}

// URL returns a request url.
func (r *Response) URL() *url.URL {
	return r.Response.Request.URL
}

// CookieValue returns a cookie value.
func (r *Response) CookieValue(name string) string {
	lowerName := strings.ToLower(name)
	for _, c := range r.Response.Cookies() {
		if strings.ToLower(c.Name) == lowerName {
			return c.Value
		}
	}
	return ""
}

// CookiesMap returns the response cookies by map.
func (r *Response) CookiesMap() map[string]string {
	cookies := map[string]string{}
	for _, c := range r.Response.Cookies() {
		cookies[c.Name] = c.Value
	}
	return cookies
}

// HeaderValue returns header value.
func (r *Response) HeaderValue(name string) string {
	lowerName := strings.ToLower(name)
	for k, v := range r.Header {
		if strings.ToLower(k) == lowerName {
			return v[0]
		}
	}
	return ""
}

// Content returns response body by byte.
func (r *Response) Content() (bs []byte, err error) {
	if r.rawBody != nil {
		return r.rawBody, nil
	}
	defer r.Body.Close()
	encoding := r.HeaderValue("Content-Encoding")
	body := r.Body
	if encoding == "gzip" {
		body, err = gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		defer body.Close()
	}
	bs, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	r.rawBody = bs
	return bs, err
}

// ContentType returns content-type in response header..
func (r *Response) ContentType() string {
	return r.HeaderValue("Content-Type")
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

// JSON returns unmarshal response body.
func (r *Response) JSON(t interface{}) error {
	rawBody, err := r.Content()
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawBody, t)
	return err
}
