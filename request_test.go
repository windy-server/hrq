package hrq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			t.Fatalf("h1 is wrong in TestGetRequest(). h1 is %#v", h1)
		}
		h3 := r.Header.Get("h3")
		if h3 != "h4" {
			t.Fatalf("h3 is wrong in TestGetRequest(). h3 is %#v", h3)
		}
		// test for query
		v1 := r.URL.Query().Get("abc")
		if v1 != "def" {
			t.Fatalf("v1 is wrong in TestGetRequest(). v1 is %#v", v1)
		}
		v2 := r.URL.Query().Get("hij")
		if v2 != "klm" {
			t.Fatalf("v2 is wrong in TestGetRequest(). v2 is %#v", v2)
		}
		// test for cookie
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Fatalf("c1 is wrong in TestGetRequest(). c1 is %#v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Fatalf("c2 is wrong in TestGetRequest(). c2 is %#v", c2)
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
	params := map[string]string{
		"abc": "def",
		"hij": "klm",
	}
	url := MakeURL(server.URL, params)
	req, _ := Get(url)
	req.SetHeader("h1", "h2")
	req.SetHeader("h3", "h4")
	req.PutCookie("c1", "v1")
	req.PutCookie("c2", "v2")
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "FooBar" {
		t.Fatalf("text is wrong in TestGetRequest(). text is %#v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Fatalf("CookiesMap() is wrong in TestGetRequest(). cm is %#v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Fatalf("CookieValue() is wrong in TestGetRequest(). a is %#v", a)
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		foo := r.PostForm["foo"][0]
		if foo != "123" {
			t.Fatalf("foo is wrong in TestGetRequest(). foo is %#v", foo)
		}
		bar := r.PostForm["bar"][0]
		if bar != "&456" {
			t.Fatalf("bar is wrong in TestGetRequest(). bar is %#v", bar)
		}
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Fatalf("c1 is wrong in TestGetRequest(). c1 is %#v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Fatalf("c2 is wrong in TestGetRequest(). c2 is %#v", c2)
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
		"bar": []string{"&456"},
	}
	req, _ := Post(url, data)
	req.SetApplicationFormUrlencoded()
	req.PutCookie("c1", "v1")
	req.PutCookie("c2", "v2")
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "FooBar" {
		t.Fatalf("text is wrong in TestGetRequest(). text is %#v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Fatalf("CookiesMap() is wrong in TestGetRequest(). cm is %#v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Fatalf("CookieValue() is wrong in TestGetRequest(). a is %#v", a)
	}
}

func TestMultipartFormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)
		foo := r.FormValue("foo")
		if foo != "123" {
			t.Fatalf("foo is wrong in TestGetRequest(). foo is %#v", foo)
		}
		bar := r.FormValue("bar")
		if bar != "&456" {
			t.Fatalf("bar is wrong in TestGetRequest(). bar is %#v", bar)
		}
	}))
	url := server.URL
	data := map[string]string{
		"foo": "123",
		"bar": "&456",
	}
	req, _ := Post(url, data)
	req.SetMultipartFormData()
	req.Send()
}

func TestApplicationJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf(err.Error())
		}
		var list []string
		err = json.Unmarshal(body, &list)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if list[0] != "foo" || list[1] != "bar" {
			t.Fatalf("Request data is wrong in TestJSON()")
		}
		fmt.Fprintf(w, `["abc", "efg"]`)
	}))
	url := server.URL
	data := []string{
		"foo",
		"bar",
	}
	req, _ := Post(url, data)
	req.SetApplicationJSON()
	res, _ := req.Send()
	var d []string
	j, _ := res.JSON(d)
	switch list := j.(type) {
	case []interface{}:
		v1, _ := list[0].(string)
		v2, _ := list[0].(string)
		if v1 != "abc" && v2 != "efg" {
			t.Fatalf("list is wrong in TestJSON(). list is %#v", list)
		}
	default:
		t.Fatalf("list is wrong in TestJSON(). list is %#v", list)
	}
}

func TestHeader(t *testing.T) {
	r, _ := Get("http://example.com")
	r.SetHeader("foo", "bar")
	v := r.HeaderValue("foo")
	if v != "bar" {
		t.Fatalf("SetHeader is wrong. v is %#v", v)
	}
	r.DelHeader("foo")
	v = r.HeaderValue("foo")
	if v != "" {
		t.Fatalf("DelHeader is wrong. v is %#v", v)
	}
}

func TestGet(t *testing.T) {
	req, _ := Get("http://example.com")
	if req.Method != "GET" {
		t.Fatalf("req.Method is wrong by Get(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Get(). req.Timeout is %#v", req.Timeout)
	}
}

func TestPost(t *testing.T) {
	req, _ := Post("http://example.com", nil)
	if req.Method != "POST" {
		t.Fatalf("req.Method is wrong by Post(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Post(). req.Timeout is %#v", req.Timeout)
	}
	ct := req.HeaderValue("Content-Type")
	if ct != DefaultContentType {
		t.Fatalf("Content-Type is wrong by Post(). Content-Type is %#v", ct)
	}
}

func TestSetTimeout(t *testing.T) {
	req, _ := Get("http://example.com")
	req.SetTimeout(100)
	timeout := time.Duration(100) * time.Second
	if req.Timeout != timeout {
		t.Fatalf("req.SetTimeout() is wrong.")
	}
}
