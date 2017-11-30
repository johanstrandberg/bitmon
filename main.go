package main

import (
	"fmt"
	"time"
	"net/http"
	"log"
	"os"
	"strconv"
	"io/ioutil"
	"gopkg.in/mailgun/mailgun-go.v1"
	"encoding/json"
)

type Configuration struct {
    MailList    []string
	MailSecret  string
	MailKey		string
	MailDomain	string
	Threshold	[]float64
	URL			string
}

var configuration Configuration

func main() {
	f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	err = nil
	log.SetOutput(f)
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err = decoder.Decode(&configuration)
	alarm := Alarm{ thresholds: configuration.Threshold }
	if err != nil {
	  fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", configuration)
	client := &http.Client{}
	for true {
		req, err := http.NewRequest("GET", configuration.URL, nil)
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return
		}
		var f float64
		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			f, _ = strconv.ParseFloat(bodyString, 64)
			f = f / 1000000000.0
			t := time.Now()
			log.Println(fmt.Sprintf("%+v", f) + " - " + t.UTC().Format("15:04:05"))
			alert := alarm.Check(f)
			if alert {
				for _, e := range configuration.MailList {
					sendMail(e, "bitmon@samuraiclick.com", fmt.Sprintf("Alert! Bitcoin hashrate at %v!", f), t.UTC().String())
				}
			}
			
			resp.Body.Close()
		} else {
			log.Println("problem " + resp.Status)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func sendMail(to string, from string, subject string, body string) {
	if len(configuration.MailDomain) < 1 || len(configuration.MailKey) < 1 || len(configuration.MailSecret) < 1 {
		log.Println("Missing env vars for sending mail.")
		return
	}
	log.Println("Sending new mail ...")
	mg := mailgun.NewMailgun(configuration.MailDomain, configuration.MailSecret, configuration.MailKey)
	message := mailgun.NewMessage(
		from,
		subject,
		body,
		to)
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)
	return
}