// Simple tcp client package. This is not a runnable package.
// You have to import and provide a listener.
package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"bakkers.us/message"
)

var con net.Conn
var hostPort string
var handler func(message.Message)

// Start the tcp client and listen for server messages
func Start(host string, port int32, messageHandler func(message.Message)) {
	hostPort = host + ":" + fmt.Sprintf("%v", port)
	fmt.Println("host and port: " + hostPort)
	handler = messageHandler

	connectToServer()
}

// Stops the server connection
func Stop() {
	con.Close()
}

// Tries to connect to server and if fails, keeps trying every 5 seconds
func connectToServer() {
	var err error
	con, err = net.Dial("tcp", hostPort)
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		connectToServer()
	}
	log.Println("Connection established " + con.LocalAddr().String())

	serverReader := bufio.NewReader(con)
	go receiveLoop(*serverReader, handler)
}

// receiveLoop handles all client communication after a connection 
// has been established.
func receiveLoop(reader bufio.Reader, handler func(message.Message)) {
	for {

		// Waiting for the server response
		serverResponse, err := reader.ReadString('\n')

		switch err {
		case nil:
			messageStr := strings.TrimSpace(serverResponse)
			message := message.NewMessage()
			json.Unmarshal([]byte(messageStr), &message)
			handler(message)
		case io.EOF:
			log.Println("server closed the connection")
			connectToServer()
			return
		default:
			log.Printf("server error: %v\n", err)
			Stop()
			return
		}
	}

}

// Send will send messages to server
func Send(message message.Message) {
	messageJSON, _ := json.Marshal(message)
	messageJSON = append(messageJSON, "\n"...)
	if _, err := con.Write([]byte(messageJSON)); err != nil {
		log.Printf("failed to send the client request: %v\n", err)
	}
}
