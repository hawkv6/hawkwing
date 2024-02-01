package main

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var (
	port       string
	host       string
	serverName string
)

var rootCmd = &cobra.Command{
	Use: "demo",
}

var clientCmd = &cobra.Command{
	Use: "client",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("tcp6", "["+host+"]"+":"+port)
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		defer conn.Close()

		for {
			_, err = conn.Write([]byte("hello\n"))
			if err != nil {
				fmt.Println("Error sending message:", err.Error())
				break
			}

			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error reading response:", err.Error())
				break
			}

			fmt.Print("Received: ", message)
			time.Sleep(2 * time.Second) // Wait a bit before sending the next message
		}
	},
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		fmt.Print("Received: ", string(message))
		conn.Write([]byte("connected to " + serverName + " on port " + port + "\n"))
	}
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server...")

		ln, err := net.Listen("tcp6", ":"+port)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			return
		}
		defer ln.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				return
			}
			go handleConnection(conn)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&host, "host", "H", "::1", "Host to connect to")
	clientCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to connect to")
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	serverCmd.Flags().StringVarP(&serverName, "name", "n", "server", "Name of the server")
}

func main() {
	rootCmd.Execute()
}
