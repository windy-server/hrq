package hrq

import (
	"io/ioutil"
	"net/http"
)

// Response wraps http.Response.
type Response struct {
	Res  *http.Response
	Body []byte
}

// Content returns request body by byte.
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
