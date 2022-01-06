package main

import (
	"fmt"
)

const version = "0.1.0"

func main() {
	fmt.Printf("s@w test app v: %s \r\n", version)
	if err := loadConfig(); err != nil {
		fmt.Println("error read config.yaml: ", err.Error())
		fmt.Println("Press Enter to exit")
		fmt.Scanln()
		return
	}
	if err := smtpSend(&tEmail{
		From:    config.From,
		To:      config.To,
		Subject: config.Subject,
		Body:    config.Body,
	}); err != nil {
		fmt.Println("errSMTP: ", err.Error())
	} else {
		fmt.Println("E-mail has been successfully sent")
	}
	fmt.Println("Press Enter to exit")
	fmt.Scanln()
}
