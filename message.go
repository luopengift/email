package email

import (
	"encoding/json"
	"net/mail"
	"time"
)

// Message email message
type Message struct {
	mail.Header
	Body        string
	Attachments []*Attachment
}

// SetHeader set header
func (msg *Message) SetHeader(key string, value ...string) *Message {
	msg.Header[key] = value
	return msg
}

// Subject set subject
func (msg *Message) Subject(value ...string) *Message {
	msg.SetHeader("Subject", value...)
	return msg
}

// From set sender
func (msg *Message) From(value ...string) *Message {
	msg.SetHeader("From", value...)
	return msg
}

// To set receivers
func (msg *Message) To(value ...string) *Message {
	msg.SetHeader("To", value...)
	return msg
}

// Cc set cc
func (msg *Message) Cc(value ...string) *Message {
	msg.SetHeader("Cc", value...)
	return msg
}

// Bcc set bcc
func (msg *Message) Bcc(value ...string) *Message {
	msg.SetHeader("Bcc", value...)
	return msg
}

//ReplyTo reply to
func (msg *Message) ReplyTo(values ...string) *Message {
	msg.SetHeader("Reply-To", values...)
	return msg
}

// Now set time now
func (msg *Message) Now() *Message {
	msg.SetHeader("Date", time.Now().Format(time.RFC1123Z))
	return msg
}

// Version set default version
func (msg *Message) Version() *Message {
	msg.SetHeader("MIME-Version", "1.0")
	return msg
}

// HTML html
func (msg *Message) HTML(body string) *Message {
	msg.SetHeader("Content-Type", "text/html")
	msg.Body = body
	return msg
}

// Text text
func (msg *Message) Text(body string) *Message {
	msg.SetHeader("Content-Type", "text/plain")
	msg.Body = body
	return msg
}

func (msg *Message) Attachment(attach *Attachment) *Message {
	msg.Attachments = append(msg.Attachments, attach)
	return msg
}

func (msg *Message) UnmarshalJSON(data []byte) error {
	m := map[string]string{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	if body, ok := m["Body"]; ok {
		msg.Text(body)
		delete(m, "Body")
	}
	for k, v := range m {
		msg.SetHeader(k, v)
	}
	return nil
}

// NewMessage new message
func NewMessage() *Message {
	msg := &Message{
		Header: make(mail.Header),
	}
	msg.Now().Version()
	return msg
}
