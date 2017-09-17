package hrq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestHistory(t *testing.T) {
	time := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if time <= 3 {
			http.Redirect(w, r, "/"+strconv.Itoa(time), http.StatusFound)
			time++
			return
		}
		fmt.Fprintf(w, "foobar")
	}))
	url := server.URL
	req, _ := Get(url)
	res, _ := req.Send()
	text, _ := res.Text()
	if text != "foobar" {
		t.Fatalf("content is wrong. %#v", text)
	}
	if len(res.History) != 3 {
		t.Fatalf("History is wrong. %#v", res.History)
	}
	if res.History[0].URL.Path != "" {
		t.Fatalf("History is wrong. %#v", res.History[0].URL.String())
	}
	if res.History[1].URL.Path != "/1" {
		t.Fatalf("History is wrong. %#v", res.History[1].URL.String())
	}
	if res.History[2].URL.Path != "/2" {
		t.Fatalf("History is wrong. %#v", res.History[2].URL.String())
	}
	if res.Request.URL.Path != "/3" {
		t.Fatalf("History is wrong. %#v", res.History[3].URL.String())
	}
}
