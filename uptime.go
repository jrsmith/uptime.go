package main

import (
  "sync"
  "fmt"
  "net/http"
  "net/smtp"
  "time"
)

var wg sync.WaitGroup
const delay = 1e9 * 100

func main() {

  var domain_list = []string{
    "http://google.com",
    "http://facebook.com",
    "http://reddit.com"
  }

  for {
    go ping(domain_list)
    time.Sleep(delay)
  }

}

func ping(domain_list []string) {
  for i := range domain_list {

    fmt.Println("Pinging", domain_list[i])

    wg.Add(1)

    go func(domain string) {

      res, err := http.Head(domain)

      if err == nil {
        if res.StatusCode == 200 {
          fmt.Println("Site is up")
        } else {
          fmt.Println("Site is down:", res.Status)
          go alert(res.Status)
        }
      } else {
        fmt.Println("Site didn't respond:", err)
        go alert(err.Error())
      }

      wg.Done()

    }(domain_list[i])

    wg.Wait()

  }
}

func alert(error_message string) {

  fmt.Println("Sending mail")

  subject := "Subject: Test\r\n\r\n"

  auth := smtp.PlainAuth(
    "",
    "noreply@example.com",
    "password",
    "smtp.example.com",
  )

  err := smtp.SendMail(
    "smtp.example.com:587",
    auth,
    "noreply@example.com",
    []string{"recipient@example.com"},
    []byte(subject + error_message),
  )

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Mail sent")
  }

}
