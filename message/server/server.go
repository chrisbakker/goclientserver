// Simple tcp server package. This is not a runnable package.
// You have to import and provide a listener.
package server

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"fmt"
	"log"
	"net"

	"bakkers.us/message"
)

var selectedPort int32 = 9999
var handler func(message.Message) message.Message

// Start the tcp client and listen for tcp connections.
func Start(messageHandler func(message.Message) message.Message) int32 {
	handler = messageHandler
	listener, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(selectedPort))
	if err != nil {
		log.Println("Could not listen on port "+fmt.Sprint(selectedPort), err)
		for i := 1; i < 5; i++ {
			selectedPort++
			listener, err = net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(selectedPort))
			if err != nil {
				log.Println("Could not listen on port "+fmt.Sprint(selectedPort), err)
			} else {
				break
			}
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	go handleListenerLoop(listener)

	//	defer listener.Close()
	return selectedPort
}

// Forever loops looking for new tcp connections
func handleListenerLoop(listener net.Listener) {
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClientRequest(con)
	}
}

// Once a client connection is established, this function handles all messaging
func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			messageStr := strings.TrimSpace(clientRequest)
			log.Println(messageStr)
			message := message.NewMessage()
			json.Unmarshal([]byte(messageStr), &message)

			response := handler(message)
			messageJSON, _ := json.Marshal(response)
			messageJSON = append(messageJSON, "\n"...)
			if _, err = con.Write([]byte(messageJSON)); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
