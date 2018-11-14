package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/luopengift/email"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/luopengift/types"
)

var config = &email.Config{}

// Mail mail
type Mail struct {
	*email.Config
	gohttp.APIHandler
}

// Initialize init
func (m *Mail) Initialize() {
	m.Config = config
}

// GET method
func (m *Mail) GET() {
	log.Info("%#v", m.Config)
	m.Output(m.Config)
}

// POST method
func (m *Mail) POST() {
	msg := email.NewMessage()
	if m.Err = json.Unmarshal(m.GetBodyArgs(), &msg); m.Err != nil {
		m.Set(101, "unmarshal post body error")
		return
	}
	smtp := email.New(m.Config)
	if m.Err = smtp.Init(); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "init error")
		return
	}
	if m.Err = smtp.Auth(); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "auth error")
		return
	}
	//log.Display("dd", msg)
	if m.Err = smtp.Send(msg); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "send error")
	}
}

func main() {
	c := flag.String("c", "conf.yml", "(conf)配置文件")
	addr := flag.String("http", ":8888", "(ip:port)IP:端口")
	flag.Parse()
	if err := types.ParseConfigFile(config, *c); err != nil {
		log.Error("%v", err)
		return
	}
	fmt.Println(config)
	app := gohttp.Init()
	app.Route("/api/v1/email", &Mail{})
	app.Run(*addr)
}
