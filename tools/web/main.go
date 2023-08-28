package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <host-b|host-c> <port>")
		return
	}

	host := os.Args[1]
	port := os.Args[2]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch host {
		case "host-b":
			fmt.Fprintln(w, "Welcome from Host-B")
		case "host-c":
			fmt.Fprintln(w, "Welcome from Host-C")
		default:
			http.Error(w, "Invalid host", http.StatusBadRequest)
		}
	})

	fmt.Printf("Starting server on host %s at port %s\n", host, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
