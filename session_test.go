package hrq

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSession(t *testing.T) {
	time := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var values [][]string
		// test for cookie
		if time == 1 {
			values = [][]string{
				{"a1", "b1"},
				{"c1", "d1"},
			}
		} else {
			values = [][]string{
				{"a2", "b2"},
				{"c2", "d2"},
			}
			a1, _ := r.Cookie("a1")
			if a1.Value != "b1" {
				t.Fatalf("a1 is wrong. a1 is %#v", a1)
			}
			c1, _ := r.Cookie("c1")
			if c1.Value != "d1" {
				t.Fatalf("c1 is wrong. c1 is %#v", c1)
			}
		}
		for _, v := range values {
			cookie := &http.Cookie{
				Name:  v[0],
				Value: v[1],
			}
			http.SetCookie(w, cookie)
		}
		time++
	}))
	session, _ := NewSession()
	req, _ := Get(server.URL)
	res, _ := session.Send(req)
	cookies := map[string]string{
		"a1": "b1",
		"c1": "d1",
	}
	cm := res.CookiesMap()
	if cookies["a1"] != cm["a1"] || cookies["b1"] != cm["b1"] {
		t.Fatalf("CookiesMap() is wrong. cm is %#v", cm)
	}
	a1 := res.CookieValue("a1")
	if a1 != "b1" {
		t.Fatalf("CookieValue() is wrong. a1 is %#v", a1)
	}
	req, _ = Get(server.URL)
	res, _ = session.Send(req)
	u, _ := url.Parse(server.URL)
	cookieList := session.Jar.Cookies(u)
	if len(cookieList) != 4 {
		t.Fatalf("cookieList is wrong. a1 is %#v", a1)
	}
	if session.CookieValue(server.URL, "a1") != "b1" {
		t.Fatalf("session.CookieValue() is wrong.")
	}
	if session.CookieValue(server.URL, "C1") != "d1" {
		t.Fatalf("session.CookieValue() is wrong.")
	}
	if session.CookieValue(server.URL, "A2") != "b2" {
		t.Fatalf("session.CookieValue() is wrong.")
	}
	if session.CookieValue(server.URL, "c2") != "d2" {
		t.Fatalf("session.CookieValue() is wrong.")
	}
}
