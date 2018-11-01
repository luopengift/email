package email

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// SMTP client
type SMTP struct {
	Host     string
	Username string
	Password string
	client   *smtp.Client
}

func (s *SMTP) auth(mechs string) (smtp.Auth, error) {
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "CRAM-MD5":
			return smtp.CRAMMD5Auth(s.Username, s.Password), nil
		case "PLAIN":
			host, _, err := net.SplitHostPort(s.Host)
			if err != nil {
				return nil, fmt.Errorf("host error:%v", s.Host)
			}
			return smtp.PlainAuth("", s.Username, s.Password, host), nil
		case "LOGIN":
			return LoginAuth(s.Username, s.Password), nil
		}
	}
	return nil, nil
}

// Init init smtp config and client
func (s *SMTP) Init(v ...interface{}) error {
	var err error
	if s.client, err = smtp.Dial(s.Host); err != nil {
		return err
	}
	return nil

}

// Send Email
func (s *SMTP) Write(b []byte) (int, error) {
	if ok, mech := s.client.Extension("AUTH"); ok {
		auth, err := s.auth(mech)
		if err != nil {
			return 0, err
		}
		if auth != nil {
			if err := s.client.Auth(auth); err != nil {
				return 0, fmt.Errorf("%T failed: %s", auth, err)
			}
		}
	}
	return 0, nil
}
