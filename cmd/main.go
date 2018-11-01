package main

import (
	"fmt"
	"time"

	"github.com/luopengift/email"
)

func main() {
	smtp := email.NewSMTP("smtp.exmail.qq.com:25", "xx", "xx")
	for i := 0; i < 20; i++ {
		now := time.Now()
		if err := smtp.Init(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("init success.")
		if err := smtp.Auth(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("auth success.")

		msg := email.NewMessage().SetFrom("xx").SetTo("xx")
		if err := smtp.Send(msg); err != nil {
			fmt.Println(err)
			return
		}

		if err := smtp.Close(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(i, "send success.", time.Since(now))
	}
}
