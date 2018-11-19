package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/luopengift/email"
	"github.com/luopengift/log"
	"github.com/luopengift/types"
	"github.com/luopengift/version"
)

func main() {
	c := flag.String("conf", "conf.yml", "(conf)配置文件")
	v := flag.Bool("version", false, "(version)版本")
	flag.Parse()
	if *v {
		log.ConsoleWithMagenta("%v", version.String())
		return
	}
	config := &email.Config{}
	if err := types.ParseConfigFile(config, *c); err != nil {
		log.Error("%v", err)
		return
	}
	smtp, err := email.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer smtp.Close()
	now := time.Now()
	msg := email.NewMessage().From("xx@xx.com").To("xx@xx.com", "xx@xx.com").
		Bcc("xx@xx.com").HTML("hello") //.Attachment(attach1).Attachment(attach2)
	if err := smtp.Send(msg); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("send success.", time.Since(now))
}
