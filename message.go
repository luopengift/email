package email

import (
	"bytes"
	"fmt"
	"time"
)

type Message struct {
	Date    time.Time `json:"Date"`
	From    string    `json:"From"`
	To      string    `json:"To"`
	Cc      string    `json:"Cc"`
	Bcc     string    `json:"Bcc"`
	Subject string    `json:"Subject"`
}

func (msg *Message) Bytes() []byte {
	var content bytes.Buffer
	fmt.Fprintf(&content, "From: %s\r\n", msg.From)
	fmt.Fprintf(&content, "To: %s\r\n", msg.To)
	fmt.Fprintf(&content, "Subject: %s\r\n", msg.Subject)
	//fmt.Fprintf(&content, "Content-Type: %s\r\n", msg.ContentType)
	fmt.Fprintf(&content, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&content, "Date: %s\r\n", msg.Date.Format(time.RFC1123Z))
	fmt.Fprintf(&content, "\r\n")

	//fmt.Fprintf(&content, "%v", ctn.Body)

	return content.Bytes()
}
