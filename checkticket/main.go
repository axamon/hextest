package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
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

	r := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			// connect to consul
			return d.DialContext(ctx, "udp", net.JoinHostPort("127.0.0.1", "8600"))
		},
	}

	cname, addr, err := r.LookupSRV(ctx, "ticket", "", "service.consul")
	if err != nil {
		log.Println(err)
		os.Exit(Error)
	}

	fmt.Println(cname, addr[0].Target, addr[0].Port, addr[0].Priority, addr[0].Weight)
	sort.Slice(addr, func(i, j int) bool { return addr[i].Weight > addr[j].Weight })
	for i := range addr {
		fmt.Println(addr[i].Port, addr[i].Target)
	}
	// minute := time.Now().Minute()

	// var exitCode int
	// switch {
	// case minute%3 == 0:
	// 	exitCode = Error
	// case minute%2 == 0:
	// 	exitCode = Warning
	// }
	// os.Exit(exitCode)
}
