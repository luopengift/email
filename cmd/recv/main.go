package main

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"

	"github.com/luopengift/email/pop3"
	"github.com/luopengift/log"
)

func main() {
	log.Info("Testing recv pop3...")
	client, err := pop3.Dial("imap.exmail.qq.com:110")
	if err != nil {
		log.Error("client err: %v", err)
		return
	}
	defer client.Quit()
	err = client.Auth("xx@qq.com", "xx")
	if err != nil {
		log.Error("auth err: %v", err)
		return
	}
	msgs, err := client.ListAll()
	if err != nil {
		log.Error("list err: %v", err)
		return
	}
	for index, msg := range msgs {
		//log.Info("%v, %#v", index, msg)
		txt, err := client.Retr(msg.Seq)
		if err != nil {
			log.Error("%v", err)
			return
		}
		msg, err := Parse(txt)
		log.Info("%v, %v", msg, err)
		s, err := base64.StdEncoding.DecodeString(msg)
		fmt.Println(string(s), err)
		if index == 0 {
			return
		}
	}
}

func Parse(txt string) (string, error) {
	msg, err := mail.ReadMessage(strings.NewReader(txt))
	if err != nil {
		log.Error("read message err: %v", err)
		return "", err
	}
	for k, v := range msg.Header {
		log.Info("%v: %v", k, v)
	}
	buf := make([]byte, 10000)
	n, err := msg.Body.Read(buf)
	if err != nil {
		log.Error("%v", err)
		return "", err
	}
	return string(buf[:n]), nil
}
