package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"
)

func main() {
	http.HandleFunc("/tickets", GetAllTickets)
	http.ListenAndServe(":80", nil)
}

// GetAllTickets retrieìves all tickets
func GetAllTickets(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	// find ticket microsrvice
	addr := findMicroservice(ctx, "ticket")

	// send ticket to microservice
	c := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://"+addr+"/tickets/getall", nil)
	if err != nil {
		log.Println("Errore nella creazione request: ", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Println("Errore nella response: ", err)
	}

	payload, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	_, _ = w.Write(payload)
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
