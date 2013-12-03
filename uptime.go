package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup
var config Config

const delay = 1e9 * 100

type SMTPConfig struct {
	Username       string
	Password       string
	Host           string
	OutgoingServer string `json:"outgoing_server"`
	From           string
}

type Config struct {
	ToEmail    []string `json:"to_email"`
	DomainList []string `json:"domain_list"`
	SMTP       SMTPConfig
}

func main() {

	runtime.GOMAXPROCS(4)

	configFile, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(configFile, &config)

	for {
		go ping(config.DomainList)
		time.Sleep(delay)
	}

}

func ping(domainList []string) {
	for i := range domainList {

		fmt.Println("Pinging", domainList[i])

		wg.Add(1)

		go func(domain string) {

			url := fmt.Sprintf("http://%s", domain)

			res, err := http.Head(url)

			if err == nil {
				if res.StatusCode == 200 {
					fmt.Println("Site is up")
				} else {
					fmt.Println("Site is down:", res.Status)
					go alert(domain, res.Status)
				}
			} else {
				fmt.Println("Site didn't respond:", err)
				go alert(domain, err.Error())
			}

			wg.Done()

		}(domainList[i])

		wg.Wait()

	}
}

func alert(domain string, error_message string) {

	fmt.Println("Sending mail")

	subject := fmt.Sprintf("Subject: %s is down\r\n\r\n", domain)

	auth := smtp.PlainAuth(
		"",
		config.SMTP.Username,
		config.SMTP.Password,
		config.SMTP.Host,
	)

	err := smtp.SendMail(
		config.SMTP.OutgoingServer,
		auth,
		config.SMTP.From,
		config.ToEmail,
		[]byte(subject+error_message),
	)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mail sent")
	}

}
