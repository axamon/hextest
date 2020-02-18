package main

import (
	"fmt"
	"net/http"
	"os"
)

// ExitCodes to convey status.
const (
	OK int = iota
	Warning
	Error
)

func main() {

	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()

	// idMicroservice := os.Args[1]
	port := os.Args[2]
	// fmt.Println(idMicroservice, port)
	resp, err := http.Get("http://127.0.0.1" + port + "/tickets/status")
	if err != nil {
		os.Exit(2)
	}
	fmt.Println(resp.StatusCode)
	if resp.StatusCode < 399 {
		os.Exit(0)
	}

	os.Exit(2)
}
