package hrq

import (
	"net/http"
	"net/http/cookiejar"
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
