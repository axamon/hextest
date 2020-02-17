package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// ExitCodes to convey status.
const (
	OK int = iota
	Warning
	Error
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	idMicroservice := os.Args[1]
	port := os.Args[2]
	fmt.Println(idMicroservice, port)
	resp, err := http.Get("http://127.0.0.1:" + port + "/tickets")
	if err != nil {
		os.Exit(2)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode < 399 {
		os.Exit(0)
	}

	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			// connect to consul
			return d.DialContext(ctx, "udp", net.JoinHostPort("127.0.0.1", "8600"))
		},
	}

	cname, addr, err := r.LookupSRV(ctx, idMicroservice+".ticket", "", "service.consul")
	if err != nil {
		log.Println(err)
		os.Exit(Warning)
	}

	fmt.Println(cname, addr[0].Target, addr[0].Port, addr[0].Priority, addr[0].Weight)
	sort.Slice(addr, func(i, j int) bool { return addr[i].Weight > addr[j].Weight })
	for i := range addr {
		if strings.HasPrefix(addr[i].Target, idMicroservice) {
			os.Exit(0)
		}
		fmt.Println(addr[i].Port, addr[i].Target)
	}

	os.Exit(1)
}
