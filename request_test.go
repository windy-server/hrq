package hrq

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// test for header
		h1 := r.Header.Get("h1")
		if h1 != "h2" {
			t.Fatalf("h1 is wrong. h1 is %#v", h1)
		}
		h3 := r.Header.Get("h3")
		if h3 != "h4" {
			t.Fatalf("h3 is wrong. h3 is %#v", h3)
		}
		// test for query
		v1 := r.URL.Query().Get("abc")
		if v1 != "def" {
			t.Fatalf("v1 is wrong. v1 is %#v", v1)
		}
		v2 := r.URL.Query().Get("hij")
		if v2 != "klm" {
			t.Fatalf("v2 is wrong. v2 is %#v", v2)
		}
		// test for cookie
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Fatalf("c1 is wrong. c1 is %#v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Fatalf("c2 is wrong. c2 is %#v", c2)
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
		t.Fatalf("text is wrong. text is %#v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Fatalf("CookiesMap() is wrong. cm is %#v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Fatalf("CookieValue() is wrong. a is %#v", a)
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		foo := r.PostForm["foo"][0]
		if foo != "123" {
			t.Fatalf("foo is wrong. foo is %#v", foo)
		}
		bar := r.PostForm["bar"][0]
		if bar != "&456" {
			t.Fatalf("bar is wrong. bar is %#v", bar)
		}
		c1, _ := r.Cookie("c1")
		if c1.Value != "v1" {
			t.Fatalf("c1 is wrong. c1 is %#v", c1)
		}
		c2, _ := r.Cookie("c2")
		if c2.Value != "v2" {
			t.Fatalf("c2 is wrong. c2 is %#v", c2)
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
	data := map[string]string{
		"foo": "123",
		"bar": "&456",
	}
	req, _ := Post(url, data)
	req.SetApplicationFormUrlencoded()
	req.PutCookie("c1", "v1")
	req.PutCookie("c2", "v2")
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "FooBar" {
		t.Fatalf("text is wrong). text is %#v", text)
	}
	cookies := map[string]string{
		"a": "b",
		"c": "d",
	}
	cm := res.CookiesMap()
	if cookies["a"] != cm["a"] || cookies["b"] != cm["b"] {
		t.Fatalf("CookiesMap() is wrong. cm is %#v", cm)
	}
	a := res.CookieValue("a")
	if a != "b" {
		t.Fatalf("CookieValue() is wrong. a is %#v", a)
	}
}

func TestMultipartFormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		foo := r.FormValue("foo")
		if foo != "123" {
			t.Fatalf("foo is wrong. foo is %#v", foo)
		}
		bar := r.FormValue("bar")
		if bar != "&456" {
			t.Fatalf("bar is wrong. bar is %#v", bar)
		}
		file, header, _ := r.FormFile("foo")
		b, _ := ioutil.ReadAll(file)
		s := string(b)
		if s != "foobar\n" {
			t.Fatalf("file is wrong. %#v", s)
		}
		if header.Filename != "foo.txt" {
			t.Fatalf("filename is wrong. %#v", header.Filename)
		}
	}))
	url := server.URL
	data := map[string]string{
		"foo": "123",
		"bar": "&456",
	}
	file, _ := os.Open("test/foo.txt")
	req, _ := Post(url, data)
	req.AddFile("text/plain", "foo", "foo.txt", file)
	req.SetMultipartFormData()
	req.Send()
}

func TestGzipRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		reader, _ := gzip.NewReader(r.Body)
		contentEncoding := r.Header["Content-Encoding"][0]
		if contentEncoding != "gzip" {
			t.Fatalf("contentEncoding is wrong. %#v", contentEncoding)
		}
		defer reader.Close()
		body, _ := ioutil.ReadAll(reader)
		s := string(body)
		if s != "foo=123" {
			t.Fatalf("Request.UseGzip() is wrong. %#v", s)
		}
	}))
	url := server.URL
	data := map[string]string{
		"foo": "123",
	}
	req, _ := Post(url, data)
	req.SetApplicationFormUrlencoded().UseGzip()
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
			t.Fatalf("Request data is wrong")
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
	res.JSON(&d)
	v1 := d[0]
	v2 := d[1]
	if v1 != "abc" && v2 != "efg" {
		t.Fatalf("list is wrong. d is %#v", d)
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

func TestDelete(t *testing.T) {
	req, _ := Delete("http://example.com")
	if req.Method != "DELETE" {
		t.Fatalf("req.Method is wrong by Delete(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Delete(). req.Timeout is %#v", req.Timeout)
	}
}

func TestHead(t *testing.T) {
	req, _ := Head("http://example.com")
	if req.Method != "HEAD" {
		t.Fatalf("req.Method is wrong by Head(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Head(). req.Timeout is %#v", req.Timeout)
	}
}

func TestOptions(t *testing.T) {
	req, _ := Options("http://example.com")
	if req.Method != "OPTIONS" {
		t.Fatalf("req.Method is wrong by Options(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Options(). req.Timeout is %#v", req.Timeout)
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

func TestPut(t *testing.T) {
	req, _ := Put("http://example.com", nil)
	if req.Method != "PUT" {
		t.Fatalf("req.Method is wrong by Put(). req.Method is %#v", req.Method)
	}
	if req.Timeout != time.Duration(DefaultTimeout)*time.Second {
		t.Fatalf("req.Timeout is wrong by Put(). req.Timeout is %#v", req.Timeout)
	}
	ct := req.HeaderValue("Content-Type")
	if ct != DefaultContentType {
		t.Fatalf("Content-Type is wrong by Put(). Content-Type is %#v", ct)
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
