package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	"github.com/luopengift/log"
	"github.com/luopengift/types"
)

// SMTP client
type SMTP struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
	SSL      bool   `json:"ssl" yaml:"ssl"` //not use
	client   *smtp.Client
}

// NewSMTP new smtp
func NewSMTP(host, port, username, password string) *SMTP {
	return &SMTP{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Timeout:  1,
		SSL:      false,
	}
}

// SetTimeout set timeout
func (s *SMTP) SetTimeout(timeout int) {
	s.Timeout = timeout
}

func (s *SMTP) auth(mechs string) (smtp.Auth, error) {
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "LOGIN":
			return LoginAuth(s.Username, s.Password), nil
		case "CRAM-MD5":
			return smtp.CRAMMD5Auth(s.Username, s.Password), nil
		case "PLAIN":
			return smtp.PlainAuth("", s.Username, s.Password, s.Host), nil
		}
	}
	return nil, nil
}

// Parse smtp from v
func (s *SMTP) Parse(v interface{}) error {
	return types.Format(v, s)
}

// Init init smtp config and client
func (s *SMTP) Init() (err error) {
	server := fmt.Sprintf("%s:%s", s.Host, s.Port)
	//s.client, err = smtp.Dial(server)
	conn, err := net.DialTimeout("tcp4", server, time.Duration(s.Timeout)*time.Second)
	if err != nil {
		return nil
	}
	if s.Port == "465" {
		conn = tls.Client(conn, s.tlsConfig())
	}
	s.client, err = smtp.NewClient(conn, s.Host)
	return err
}

func (s *SMTP) tlsConfig() *tls.Config {
	return &tls.Config{
		ServerName:         s.Host,
		InsecureSkipVerify: true,
	}

}

// Auth auth
func (s *SMTP) Auth() error {
	//Check if TLS is required
	if s.Port == "465" || s.Port == "587" {
		if ok, _ := s.client.Extension("STARTTLS"); ok {
			if err := s.client.StartTLS(s.tlsConfig()); err != nil {
				return err
			}
		}
	}

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

// Noop send noop cmd
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

// Close connection
func (s *SMTP) Close() error {
	return s.client.Quit()
}

// Send message
func (s *SMTP) Send(msg *Message) error {

	var buf bytes.Buffer
	for key := range msg.Header {
		switch key {
		case "From":
			froms, err := msg.Header.AddressList(key)
			if err != nil {
				return err
			}
			log.Debug("%v, %v", key, froms[0].String())
			for _, from := range froms {
				if err = s.client.Mail(from.Address); err != nil {
					return err
				}
				fmt.Fprintf(&buf, "%v: %s\r\n", key, from.String())
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
				if key != "Bcc" {
					fmt.Fprintf(&buf, "%v: %s\r\n", key, recv.String())
				}
			}
		case "Date", "Subject", "MIME-Version", "Reply-To":
			fmt.Fprintf(&buf, "%v: %s\r\n", key, msg.Header.Get(key))
		}
	}

	boundary := "f46d043c813270fc6b04c2d223da"
	if len(msg.Attachments) > 0 {
		fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
		fmt.Fprintf(&buf, "\r\n--%s\r\n", boundary)
	}

	fmt.Fprintf(&buf, "Content-Type: %s; charset=utf-8\r\n\r\n", msg.Header.Get("Content-Type"))
	fmt.Fprintf(&buf, msg.Body)
	fmt.Fprintf(&buf, "\r\n")

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			fmt.Fprintf(&buf, "\r\n\r\n--%s\r\n", boundary)

			if attachment.Inline {
				fmt.Fprintf(&buf, "Content-Type: message/rfc822\r\n")
				fmt.Fprintf(&buf, `Content-Disposition: inline; filename="%s"`, attachment.Name)
				buf.WriteString("\r\n\r\n")
				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Name)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					fmt.Fprintf(&buf, "Content-Type: %s\r\n", mimetype)
				} else {
					fmt.Fprintf(&buf, "Content-Type: application/octet-stream\r\n")
				}
				fmt.Fprintf(&buf, "Content-Transfer-Encoding: base64\r\n")

				fmt.Fprintf(&buf,
					`Content-Disposition: attachment; filename="=?UTF-8?B?%s?="`,
					base64.StdEncoding.EncodeToString([]byte(attachment.Name)),
				)
				buf.WriteString("\r\n\r\n")
				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			fmt.Fprintf(&buf, "\r\n--%s", boundary)
		}
		buf.WriteString("--")
	}

	_, err := s.Write(buf.Bytes())
	return err
}
