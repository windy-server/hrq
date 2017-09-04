package hrq

import (
	"testing"
	"time"
)

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

func TestGet(t *testing.T) {
	req, _ := Get("http://example.com")
	if req.Method != "GET" {
		t.Errorf("req.Method is wrong by Get(). req.Method is %v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Errorf("req.Timeout is wrong by Get(). req.Timeout is %v", req.Timeout)
	}
}

func TestPost(t *testing.T) {
	req, _ := Post("http://example.com", nil)
	if req.Method != "POST" {
		t.Errorf("req.Method is wrong by Post(). req.Method is %v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Errorf("req.Timeout is wrong by Post(). req.Timeout is %v", req.Timeout)
	}
	ct := req.GetHeader("Content-Type")
	if ct != DefaultContentType {
		t.Errorf("Content-Type is wrong by Post(). Content-Type is %v", ct)
	}
}

func TestSetTimeout(t *testing.T) {
	req, _ := Get("http://example.com")
	req.SetTimeout(100)
	timeout := time.Duration(100) * time.Second
	if req.Timeout != timeout {
		t.Errorf("req.SetTimeout() is wrong.")
	}
}
