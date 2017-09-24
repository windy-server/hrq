package hrq

import (
	"net/http"
	"net/http/cookiejar"
	Url "net/url"
	"strings"
)

// Session is a session.
type Session struct {
	*http.Client
}

// NewSession return a session.
func NewSession() (s *Session, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	cli := &http.Client{
		Jar: jar,
	}
	s = &Session{
		Client: cli,
	}
	return
}

// Send send a request.
func (s *Session) Send(r *Request) (res *Response, err error) {
	s.Timeout = r.Timeout
	return send(s, r)
}

// CookieValue returns a cookie value.
func (s *Session) CookieValue(url, name string) string {
	u, _ := Url.Parse(url)
	cookieList := s.Jar.Cookies(u)
	lowerName := strings.ToLower(name)
	for _, c := range cookieList {
		if strings.ToLower(c.Name) == lowerName {
			return c.Value
		}
	}
	return ""
}
