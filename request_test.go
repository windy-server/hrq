package hrq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// test for header
		h1 := r.Header.Get("h1")
		if h1 != "h2" {
			t.Errorf("h1 is wrong in TestGetRequest(). h1 is %v", h1)
		}
		h3 := r.Header.Get("h3")
		if h3 != "h4" {
			t.Errorf("h3 is wrong in TestGetRequest(). h3 is %v", h3)
		}
		// test for query
		v1 := r.URL.Query().Get("abc")
		if v1 != "def" {
			t.Errorf("v1 is wrong in TestGetRequest(). v1 is %v", v1)
		}
		v2 := r.URL.Query().Get("hij")
		if v2 != "klm" {
			t.Errorf("v2 is wrong in TestGetRequest(). v2 is %v", v2)
		}
		// test for cookie
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Errorf("c1 is wrong in TestGetRequest(). c1 is %v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Errorf("c2 is wrong in TestGetRequest(). c2 is %v", c2)
		}
		values := [][]string{
			[]string{"a", "b"},
			[]string{"c", "d"},
		}
		for _, v := range values {
			cookie := &http.Cookie{
				Name:  v[0],
				Value: v[1],
			}
			http.SetCookie(w, cookie)
		}
		fmt.Fprintf(w, "FooBar")
	}))
	url := server.URL + "?abc=def&hij=klm"
	req, _ := Get(url)
	req.SetHeader("h1", "h2")
	req.SetHeader("h3", "h4")
	req.PutCookie("c1", "v1")
	req.PutCookie("c2", "v2")
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "FooBar" {
		t.Errorf("text is wrong in TestGetRequest(). text is %v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Errorf("CookiesMap() is wrong in TestGetRequest(). cm is %v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Errorf("CookieValue() is wrong in TestGetRequest(). a is %v", a)
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		foo := r.PostForm["foo"][0]
		if foo != "123" {
			t.Errorf("foo is wrong in TestGetRequest(). foo is %v", foo)
		}
		bar := r.PostForm["bar"][0]
		if bar != "456" {
			t.Errorf("bar is wrong in TestGetRequest(). bar is %v", bar)
		}
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Errorf("c1 is wrong in TestGetRequest(). c1 is %v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Errorf("c2 is wrong in TestGetRequest(). c2 is %v", c2)
		}
		values := [][]string{
			[]string{"a", "b"},
			[]string{"c", "d"},
		}
		for _, v := range values {
			cookie := &http.Cookie{
				Name:  v[0],
				Value: v[1],
			}
			http.SetCookie(w, cookie)
		}
		fmt.Fprintf(w, "FooBar")
	}))
	url := server.URL
	data := map[string][]string{
		"foo": []string{"123"},
		"bar": []string{"456"},
	}
	req, _ := Post(url, data)
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	req.PutCookie("c1", "v1")
	req.PutCookie("c2", "v2")
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "FooBar" {
		t.Errorf("text is wrong in TestGetRequest(). text is %v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Errorf("CookiesMap() is wrong in TestGetRequest(). cm is %v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Errorf("CookieValue() is wrong in TestGetRequest(). a is %v", a)
	}
}

func TestHeader(t *testing.T) {
	r, _ := Get("http://example.com")
	r.SetHeader("foo", "bar")
	v := r.HeaderValue("foo")
	if v != "bar" {
		t.Errorf("SetHeader is wrong. v is %v", v)
	}
	r.DelHeader("foo")
	v = r.HeaderValue("foo")
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
	ct := req.HeaderValue("Content-Type")
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
