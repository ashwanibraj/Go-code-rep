//This is the server side of vector clock implementation example.
//With every message, the time vector is sent from host.
//To run this, place this file on a system. Place the client file on another system. 
//Start the server first using command "go VectorClocks_Server.go"

package main

import (
	"fmt"
	"sync"
	"bufio"
	"net"
	"os"
	"strings"
	"strconv"
)

var vectorClocks = struct{
	sync.RWMutex
	clock [2]int
}{clock: [2]int{0,0}}

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
			<- flag 
		}
	}		
}

func receiveFromClient(conn *net.Conn) {
	for {
		clMsg, err := bufio.NewReader(*conn).ReadString('\t') // Reading from io buffer of open connection
		if err != nil {
			fmt.Println("Error while reading text:", err.Error())			
			break //breaking from infinite loop in case of error
		} else {
			res := strings.Split(clMsg, "\n")
			cl, _ := strconv.Atoi(strings.TrimSpace(res[1]))
			vectorClocks.Lock()
			if vectorClocks.clock[1] < cl {
				vectorClocks.clock[1] =  cl
			}
			vectorClocks.Unlock()
			fmt.Printf("Client>%s\tClock: %d\t%d\n", res[0], vectorClocks.clock[0], vectorClocks.clock[1]) // Printing client messages
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
			vectorClocks.Lock()
			vectorClocks.clock[0] += 1
			vectorClocks.Unlock()
			var sendVal = serMsg + strconv.Itoa(vectorClocks.clock[0]) + "\t"						
			fmt.Fprintf(*conn, sendVal)
		}
	}
	flag <- true
}