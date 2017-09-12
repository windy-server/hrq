package hrq

import (
	"net/http"
	"testing"
)

func TestHeaderValue(t *testing.T) {
	res := &Response{
		Response: &http.Response{
			Header: http.Header{},
		},
	}
	res.Header = http.Header(map[string][]string{
		"foo": []string{"bar"},
	})
	value := res.HeaderValue("foo")
	if value != "bar" {
		t.Fatalf("Response header value is wrong. value is %#v", value)
	}
}
