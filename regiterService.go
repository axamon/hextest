package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// ServiceData contains the information to register the service.
type ServiceData struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Tags    []string `json:"Tags"`
	Address string   `json:"Address"`
	Port    int      `json:"Port"`
	Meta    struct {
		ServiceVersion string `json:"service_version"`
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

func registerService(id, name, version, address, port, deregister, interval, timeout string) {

	// Cleans up the port variable to make it an int.
	portclean := strings.Replace(port, ":", "", 1)
	portNum, err := strconv.Atoi(portclean)
	if err != nil {
		log.Print(err)
	}

	// Creates new ServiceData variable.
	var s ServiceData

	s.ID = id
	s.Name = name
	s.Port = portNum
	s.Address = address
	// adds in tags the version and the uuid of the microservice to retrieve from DNS.
	s.Tags = []string{version, id}
	s.Meta.ServiceVersion = version
	s.EnableTagOverride = false
	s.Check.DeregisterCriticalServiceAfter = deregister
	s.Check.Args = []string{"checkticket", id, port}
	s.Check.Interval = interval
	s.Check.Timeout = timeout
	s.Weights.Passing = 10
	s.Weights.Warning = 1

	// marshal the struct s in the json fields.
	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	// Creates the payload to send.
	payload := bytes.NewBuffer(jsonData)

	// Creates an http client.
	c := &http.Client{}

	// Creates the http request.
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

	// Executes the http request.
	_, err = c.Do(req)
	if err != nil {
		log.Println(err)
	}
}
