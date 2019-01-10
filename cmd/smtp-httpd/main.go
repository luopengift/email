package main

import (
	"encoding/json"

	"github.com/luopengift/email"
	"github.com/luopengift/framework"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
)

var config = map[string]*email.Config{}

// Mail mail
type Mail struct {
	gohttp.APIHandler
}

// GET method
func (m *Mail) GET() {
	log.Info("%#v", config)
	m.Output(config)
}

// POST method
func (m *Mail) POST() {
	msg := email.NewMessage()
	if m.Err = json.Unmarshal(m.GetBodyArgs(), &msg); m.Err != nil {
		m.Set(101, "unmarshal post body error")
		return
	}
	conf, ok := config[msg.Get("From")]
	if !ok {
		m.Err = log.Errorf("can not find email config by Username=%v", msg.Get("From"))
		log.Error("%v", m.Err)
		m.Set(101, m.Err.Error())
		return
	}
	var smtp *email.SMTP
	if smtp, m.Err = email.New(conf); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "new error")
		return
	}
	defer smtp.Close()

	var txt string
	if len(msg.Body) > 1000 {
		txt = msg.Body[:1000]
	} else {
		txt = msg.Body
	}
	log.Info("send From: %v, To: %v, Cc: %v Subject: %v \n=> %v", msg.Get("From"), msg.Get("To"), msg.Get("Cc"), msg.Get("Subject"), txt)
	if m.Err = smtp.Send(msg); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "send error")
	}
}

func main() {
	framework.BindConfig(&config)
	framework.HttpdRoute("^/api/v1/email$", &Mail{})
	framework.Run()
}
