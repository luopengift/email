package main

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/luopengift/email"
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
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
	body := map[string]string{}
	if m.Err = json.Unmarshal(m.GetBodyArgs(), &body); m.Err != nil {
		m.Set(101, "unmarshal post body error")
		return
	}
	msg := email.NewMessage().From(body["From"]).To(strings.Split(body["To"], ",")...).Subject(body["Subject"]).Text(body["body"])
	if body["Cc"] != "" {
		msg.Cc(strings.Split(body["Cc"], ",")...)
	}
	//log.Display("msg", msg)
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
	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	addr := flag.String("http", ":8888", "(http)IP:端口")
	v := flag.Bool("version", false, "(version)版本")
	flag.Parse()
	if *v {
		log.ConsoleWithMagenta("%v", version.String())
		return
	}
	file := log.NewFile("%Y-%M-%D.log")
	file.SetMaxBytes(200 * 1024 * 1024) // 200MB
	log.SetOutput(file)

	if err := types.ParseConfigFile(&config, *c); err != nil {
		log.Error("%v", err)
		return
	}
	log.Debug("config: %v", config)
	app := gohttp.Init()
	app.Log.SetOutput(file)
	app.Route("^/api/v1/email$", &Mail{})
	app.RouteFunCtx("^/-/reload$", func(ctx *gohttp.Context) {
		if err := types.ParseConfigFile(&config, *c); err != nil {
			log.Error("%v", err)
			ctx.Output(err, 400)
			return
		}
		ctx.Output("ok")
	})
	log.Info("init ok!")
	app.Run(*addr)
}
