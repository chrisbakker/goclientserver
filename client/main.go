// Example client application that uses the message client package
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"bakkers.us/message"
	"bakkers.us/message/client"
)

func main() {

	client.Start("localhost", 9999, handleMessage)
	fmt.Println("Type message to server.")
	clientReader := bufio.NewReader(os.Stdin)

	for {
		// Will read a string from standard in,
		// Wrap it in a message and send it to the server
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			input := strings.TrimSpace(clientRequest)
			if strings.ToLower(input) == "quit" {
				fmt.Println("Shutting down")
				return
			}
			message := message.NewMessage()
			message.Status = "ok"
			message.Body = input
			client.Send(message)
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

	}
}

func handleMessage(message message.Message) {
	fmt.Println("handleMessage: " + message.Body)
}
