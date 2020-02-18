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
	"strconv"
	"strings"
	"time"

	"github.com/axamon/hextest/ticket"
)

func main() {
	ctx := context.Background()

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
	// fmt.Println(jsonData)

	// Creates the payload to send.
	payload := bytes.NewBuffer(jsonData)

	// find ticket microsrvice
	addr := findMicroservice(ctx, "ticket")

	fmt.Println(addr)

	// aggiornamento contesto.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// send ticket to microservice
	c := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://"+addr+"/tickets/new", payload)
	if err != nil {
		log.Println("Errore nella creazione request: ", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("Errore nella response: ", err)
	}

	if resp.StatusCode < 299 {
		return
	}
	log.Fatal("Something went wrong", resp.StatusCode)
}

// findMicroservice finds the most suitable microservice available.
func findMicroservice(ctx context.Context, service string) string {
	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			// connects to consul
			return d.DialContext(ctx, "udp", net.JoinHostPort("127.0.0.1", "8600"))
		},
	}

	// retrievs cname and corresponding adresses.
	cname, addr, err := r.LookupSRV(ctx, service, "", "service.consul")
	if err != nil {
		log.Println(err)
	}
	ipSlice, err := r.LookupIPAddr(ctx, cname)
	if err != nil {
		log.Println(err)
	}
	ip := ipSlice[0].IP.String()

	// sorts microservices from better to worst.
	sort.Slice(addr, func(i, j int) bool { return addr[i].Weight > addr[j].Weight })

	port := strconv.Itoa(int(addr[0].Port))

	return ip + ":" + port
}

func getInfo(s string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter " + s + ": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	return strings.ReplaceAll(text, "\n", "")
}
