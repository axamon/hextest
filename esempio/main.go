package main

import (
	"fmt"

	"github.com/lixiangzhong/dnsutil"
)

func main() {

	var dig dnsutil.Dig
	dig.SetDNS("127.0.0.1")                  //or ns.xxx.com
	a, err := dig.A("ticket.service.consul") // dig google.com @8.8.8.8
	fmt.Println(a, err)
}
