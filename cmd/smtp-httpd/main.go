package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/luopengift/log"
	"github.com/luopengift/types"

	"github.com/luopengift/email"
	"github.com/luopengift/gohttp"
)

var mail = &email.SMTP{}

type Mail struct {
	*email.SMTP
	gohttp.APIHandler
}

func (m *Mail) Initialize() {
	m.SMTP = mail
}

func (m *Mail) GET() {
	log.Info("%#v", m.SMTP)
	m.Output(m.SMTP)
}

func (m *Mail) POST() {
	msg := &email.Message{}
	if m.Err = json.Unmarshal(m.GetBodyArgs(), msg); m.Err != nil {
		m.Set(101, "unmarshal post body error")
		return
	}
	if m.Err = m.SMTP.Init(); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "init error")
		return
	}
	if m.Err = m.SMTP.Auth(); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "auth error")
		return
	}
	//log.Display("dd", msg)
	if m.Err = m.SMTP.Send(msg); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "send error")
	}
}

func main() {
	c := flag.String("c", "conf.yml", "(conf)配置文件")
	p := flag.String("p", "8888", "(port)端口")
	flag.Parse()
	if err := types.ParseConfigFile(mail, *c); err != nil {
		log.Error("%v", err)
		return
	}
	app := gohttp.Init()
	app.Route("/api/v1/email", &Mail{})
	app.Run(fmt.Sprintf(":%v", *p))
}
