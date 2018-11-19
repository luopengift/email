package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"path/filepath"
)

// SendExt is used some special scenario,
// must require MAIL command before RCPT command.
func (s *SMTP) SendExt(msg *Message) error {
	return s.sendExt(msg)
}
func (s *SMTP) sendExt(msg *Message) error {
	var buf bytes.Buffer
	for _, key := range []string{"Date", "Subject", "MIME-Version", "Reply-To"} {
		fmt.Fprintf(&buf, "%v: %s\r\n", key, msg.Header.Get(key))
	}
	from := msg.Header.Get("From")
	if from == "" {
		msg.From(s.Username)
	}
	froms, err := msg.Header.AddressList(from)
	if err != nil {
		return err
	}
	for _, from := range froms {
		if err = s.client.Mail(from.Address); err != nil {
			return err
		}
		fmt.Fprintf(&buf, "From: %s\r\n", from.String())
	}
	for _, key := range []string{"To", "Cc", "Bcc"} {
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
	_, err = s.Write(buf.Bytes())
	return err
}
