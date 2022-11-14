// Example server application that uses the message server package
package main

import (
	"fmt"
	"time"

	"bakkers.us/message"
	"bakkers.us/message/server"
)

func main() {
	port := server.Start(messageHandler)
	fmt.Println("server port: " + fmt.Sprint(port))

	for {
		time.Sleep(1 * time.Second)
	}
}

func messageHandler(message message.Message) message.Message {
	body := message.Body
	body = "Echo: " + body
	message.Body = body
	return message
}
