package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/axamon/hextest/ticket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// create tt to send.
	tt := new(ticket.Ticket)

	tt.Creator = getInfo("Creator")
	tt.Description = getInfo("Description")
	tt.Title = getInfo("Title")

	// marshal the struct s in the json fields.
	jsonData, err := json.Marshal(tt)
	if err != nil {
		log.Println(err)
	}

	// Creates the payload to send.
	payload := bytes.NewBuffer(jsonData)

	// find ticket microsrvice
	addr := findMicroservice(ctx, "ticket")

	// send ticket to microservice
	c := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://"+addr.Target+":"+string(addr.Port)+"/tickets", payload)
	if err != nil {
		log.Println(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode < 299 {
		return
	}
	log.Fatal("Something went wrong")
}

// findMicroservice finds the most suitable microservice available.
func findMicroservice(ctx context.Context, service string) *net.SRV {
	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			// connects to consul
			return d.DialContext(ctx, "udp", net.JoinHostPort("127.0.0.1", "8600"))
		},
	}

	_, addr, err := r.LookupSRV(ctx, service, "", "service.consul")
	if err != nil {
		log.Println(err)
	}

	// sorts microservices from better to worst.
	sort.Slice(addr, func(i, j int) bool { return addr[i].Weight > addr[j].Weight })

	info := addr[0]
	return info
}

func getInfo(s string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter " + s + ":")
	text, _ := reader.ReadString('\n')
	return text
}
