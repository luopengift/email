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

var config []*email.Config

// Mail mail
type Mail struct {
	config []*email.Config
	gohttp.APIHandler
}

// Initialize init
func (m *Mail) Initialize() {
	m.config = config
}

// GET method
func (m *Mail) GET() {
	log.Info("%#v", m.config)
	m.Output(m.config)
}

// POST method
func (m *Mail) POST() {
	msg := email.NewMessage()
	if m.Err = json.Unmarshal(m.GetBodyArgs(), &msg); m.Err != nil {
		m.Set(101, "unmarshal post body error")
		return
	}
	var getConfig = func() (*email.Config, error) {
		for _, conf := range m.config {
			if conf.Username == msg.Get("From") {
				return conf, nil
			}
		}
		return nil, log.Errorf("can not find email config by Username=%v", msg.Get("From"))
	}
	config, err := getConfig()
	if err != nil {
		log.Error("%v", err)
		m.Set(101, err.Error())
		return
	}
	var smtp *email.SMTP
	if smtp, m.Err = email.New(config); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "new error")
		return
	}
	defer smtp.Close()
	//log.Display("dd", msg)
	if m.Err = smtp.Send(msg); m.Err != nil {
		log.Error("%v", m.Err)
		m.Set(101, "send error")
	}
}

func main() {
	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	addr := flag.String("http", ":8888", "(http)IP:端口")
	flag.Parse()
	if err := types.ParseConfigFile(&config, *c); err != nil {
		log.Error("%v", err)
		return
	}
	fmt.Println(config)
	app := gohttp.Init()
	app.Route("^/api/v1/email$", &Mail{})
	app.RouteFunCtx("^/-/reload$", func(ctx *gohttp.Context) {
		if err := types.ParseConfigFile(&config, *c); err != nil {
			log.Error("%v", err)
			ctx.Output(err, 400)
			return
		}
		ctx.Output("ok")
	})
	app.Run(*addr)
}
