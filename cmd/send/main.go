package main

import (
	"fmt"
	"time"

	"github.com/luopengift/email"
)

func main() {
	smtp := email.NewSMTP("smtp.exmail.qq.com", "465", "xxx@qq.com", "xxx")
	smtp.SSL = true
	for i := 0; i < 1; i++ {
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
		// attach1, err := email.NewAttachment("/Users/xxx/Desktop/list.go", false)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// attach2, _ := email.NewAttachment("/Users/xxx/Desktop/main.go", true)
		msg := email.NewMessage().From("xxx@qq.com").To("xxx@qq.com").
			Bcc("870148195@qq.com").HTML("hello") //.Attachment(attach1).Attachment(attach2)
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