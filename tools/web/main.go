package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/mattes/go-asciibot"
)

var (
	host string
	port string
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: web [options]")
		flag.PrintDefaults()
	}
	flag.StringVar(&host, "host", "host-b", "host name which should be displayed")
	flag.StringVar(&port, "port", "8080", "port number to listen on")
}

func generateMessage(host string) string {
	image := asciibot.Random()
	return fmt.Sprintf("Welcome from %s\n%s", host, image)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request from", r.RemoteAddr)
		fmt.Fprintln(w, generateMessage(host))
	})

	fmt.Printf("Starting server on host %s at port %s\n", host, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
