package email

import (
	"net/mail"
	"time"
)

type Message struct {
	mail.Header
}

func (msg *Message) SetHeader(key string, value ...string) *Message {
	msg.Header[key] = value
	return msg
}

func (msg *Message) SetSubject(value ...string) *Message {
	msg.SetHeader("Subject", value...)
	return msg
}

func (msg *Message) SetFrom(value ...string) *Message {
	msg.SetHeader("From", value...)
	return msg
}
func (msg *Message) SetTo(value ...string) *Message {
	msg.SetHeader("To", value...)
	return msg
}
func (msg *Message) SetCc(value ...string) *Message {
	msg.SetHeader("Cc", value...)
	return msg
}
func (msg *Message) SetNow() *Message {
	msg.SetHeader("Date", time.Now().Format(time.RFC1123Z))
	return msg
}

func (msg *Message) SetVersion() *Message {
	msg.SetHeader("MIME-Version", "1.0")
	return msg
}

func NewMessage() *Message {
	msg := &Message{
		Header: make(mail.Header),
	}
	msg.SetNow().SetVersion()
	return msg
}
