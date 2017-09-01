package hrq

import "testing"

func TestHeader(t *testing.T) {
	r, _ := Get("http://example.com")
	r.SetHeader("foo", "bar")
	v := r.GetHeader("foo")
	if v != "bar" {
		t.Errorf("SetHeader is wrong. v is %v", v)
	}
	r.DelHeader("foo")
	v = r.GetHeader("foo")
	if v != "" {
		t.Errorf("DelHeader is wrong. v is %v", v)
	}
}
