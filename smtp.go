package email

import (
	"bytes"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/luopengift/types"
)

// SMTP client
type SMTP struct {
	Host     string
	Username string
	Password string
	client   *smtp.Client
}

func NewSMTP(host, username, password string) *SMTP {
	return &SMTP{
		Host:     host,
		Username: username,
		Password: password,
	}
}

func (s *SMTP) auth(mechs string) (smtp.Auth, error) {
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "LOGIN":
			return LoginAuth(s.Username, s.Password), nil
		case "CRAM-MD5":
			return smtp.CRAMMD5Auth(s.Username, s.Password), nil
		case "PLAIN":
			host, _, err := net.SplitHostPort(s.Host)
			if err != nil {
				return nil, fmt.Errorf("host error:%v", s.Host)
			}
			return smtp.PlainAuth("", s.Username, s.Password, host), nil
		}
	}
	return nil, nil
}

func (s *SMTP) Parse(v interface{}) error {
	return types.Format(v, s)
}

// Init init smtp config and client
func (s *SMTP) Init() (err error) {
	s.client, err = smtp.Dial(s.Host)
	return err
}

func (s *SMTP) Auth() error {
	if ok, mech := s.client.Extension("AUTH"); ok {
		auth, err := s.auth(mech)
		if err != nil {
			return err
		}
		if auth != nil {
			if err := s.client.Auth(auth); err != nil {
				return fmt.Errorf("%T failed: %s", auth, err)
			}
		}
	}
	return nil
}
func (s *SMTP) Noop() error {
	return s.client.Noop()
}

// Write write
func (s *SMTP) Write(b []byte) (int, error) {
	w, err := s.client.Data()
	if err != nil {
		return 0, err
	}
	defer w.Close()
	return w.Write(b)
}

func (s *SMTP) Close() error {
	return s.client.Quit()
}

func (s *SMTP) Send(msg *Message) error {
	var buf bytes.Buffer
	for key := range msg.Header {
		switch key {
		case "From":
			froms, err := msg.Header.AddressList(key)
			if err != nil {
				return err
			}
			for _, from := range froms {
				if err = s.client.Mail(from.Address); err != nil {
					return err
				}
				fmt.Fprintf(&buf, "%v: %s\r\n", key, from.Address)
			}
		case "To", "Cc", "Bcc":
			recvs, err := msg.Header.AddressList(key)
			if err != nil {
				return err
			}
			for _, recv := range recvs {
				if err = s.client.Rcpt(recv.Address); err != nil {
					return err
				}
				fmt.Fprintf(&buf, "%v: %s\r\n", key, recv.Address)
			}
		case "Date", "Subject", "Context-Type", "MIME-Version":
			fmt.Fprintf(&buf, "%v: %s\r\n", key, msg.Header.Get(key))
		}
	}
	buf.WriteString("\r\n")
	_, err := s.Write(buf.Bytes())
	return err
}
