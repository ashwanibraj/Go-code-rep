//This is the server side of my client server chat application.
//To run this, place this file on a system. Place the client file on another system. 
//Start the server first. Command is "go run Chat_server.go"
//Then start the client from another system. 
//This application uses TCP to establish connection. 
//My effort is to implement simple message passing in a network as per the concepts taught in class.

package main

import (
	"fmt"
	"bufio"
	"net"
	"os"
)

var flag = make(chan bool)

func main() {
	listener, err := net.Listen("tcp", ":8080") // Listening for new connection at port 80 using TCP
	if err != nil{
		listener.Close()
		fmt.Println("Error while connection:", err.Error())
	} else {
		conn, err := listener.Accept() //Accept the connection from client, if no error
		if err != nil {
			fmt.Println("Error from client:", err.Error())
			conn.Close()
		} else {
			fmt.Println("Connected to client.")
			go send2Client(&conn) //Calling goroutine to send messages to the client, using pointer to the connection
			go receiveFromClient(&conn) //Calling another goroutine to receive messages from the client
			<- flag //Waiting to hear on the channel. This makes sure the server doesnt 
					//close down till either an error happens or client exits
		}
	}		
}

func receiveFromClient(conn *net.Conn) {
	for {
		clMsg, err := bufio.NewReader(*conn).ReadString('\n') // Reading from io buffer of open connection
		if err != nil {
			fmt.Println("Error while reading text:", err.Error())			
			break //breaking from infinite loop in case of error
		} else {
			fmt.Println("Client>", string(clMsg)) // Printing client messages
		}
	}
	flag <- true // Pass true on to the channel
}

func send2Client(conn *net.Conn) {
	for {
		serMsg, err := bufio.NewReader(os.Stdin).ReadString('\n') //Reading from io buffer of stdin for server messages
		if err != nil {
			fmt.Println("Error while reading text from server:", err.Error())
			break
		} else {
			fmt.Fprintf(*conn, serMsg)
		}
	}
	flag <- true
}