package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
)

const registrationData = `{
	"ID": "ticket1",
	"Name": "ticket",
	"Tags": [
	  "primary",
	  "v1"
	],
	"Address": "127.0.0.1",
	"Port": 3000,
	"Meta": {
	  "ticket_version": "4.0"
	},
	"EnableTagOverride": false,
	"Check": {
	  "DeregisterCriticalServiceAfter": "90m",
	  "Args": ["/usr/local/bin/checkticket"],
	  "Interval": "10s",
	  "Timeout": "5s"
	},
	"Weights": {
	  "Passing": 10,
	  "Warning": 1
	}
  }
  `

  type ServiceData struct {
		ID      string   `json:"ID"`
		Name    string   `json:"Name"`
		Tags    []string `json:"Tags"`
		Address string   `json:"Address"`
		Port    int      `json:"Port"`
		Meta    struct {
			TicketVersion string `json:"ticket_version"`
		} `json:"Meta"`
		EnableTagOverride bool `json:"EnableTagOverride"`
		Check             struct {
			DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"`
			Args                           []string `json:"Args"`
			Interval                       string   `json:"Interval"`
			Timeout                        string   `json:"Timeout"`
		} `json:"Check"`
		Weights struct {
			Passing int `json:"Passing"`
			Warning int `json:"Warning"`
		} `json:"Weights"`
	}
  }

func registerService() {
	var s SerServiceData

	s.ID = uiid.New()
	s.Port = 3000
	
	c := &http.Client{}

	payload := bytes.NewBuffer([]byte(registrationData))

	registrationURI, err := url.ParseRequestURI("http://127.0.0.1:8500/v1/agent/service/register?replace-existing-checks=true")
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("PUT", registrationURI.String(), payload)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Valid", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = c.Do(req)
	if err != nil {
		log.Println(err)
	}
}
