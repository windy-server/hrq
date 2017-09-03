package hrq

import (
	"net/http"
	"testing"
)

func TestGetHeader(t *testing.T) {
	res := &Response{
		Response: &http.Response{
			Header: http.Header{},
		},
	}
	res.Header = http.Header(map[string][]string{
		"foo": []string{"bar"},
	})
	value := res.GetHeader("foo")
	if value != "bar" {
		t.Errorf("Response header value is wrong. value is %v", value)
	}
}
