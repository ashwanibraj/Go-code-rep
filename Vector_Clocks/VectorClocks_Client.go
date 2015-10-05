//Make sure you change the IP address in this file to the ip address of server.
//After starting the server, start the client. Command is "go run VectorClocks_Client.go"
//Local arrays are maintained with time vectors of client and server. 

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

var flag = make(chan bool) //initializing a channel variable of boolean type

func main() {
	conn, err := net.Dial("tcp", "192.168.0.20:8080") //Initiating a new connection with the server at given ip and port using TCP
	if err != nil{
		conn.Close()
		fmt.Println("Error while connection:", err.Error())
	} else {
		fmt.Println("Connected to server.")
		go send2Server(&conn)
		go receiveFromServer(&conn)
		<- flag
	}		
}

func receiveFromServer(conn *net.Conn) {
	for {		
		serMsg, err := bufio.NewReader(*conn).ReadString('\t') //Reading from open connection
		if err != nil {
			fmt.Println("Error while reading text:", err.Error())			
			break
		} else {
			res := strings.Split(serMsg, "\n")
			cl, _ := strconv.Atoi(strings.TrimSpace(res[1]))			
			vectorClocks.Lock()
			if vectorClocks.clock[0] < cl {
				vectorClocks.clock[0] =  cl
			}
			vectorClocks.Unlock()
			fmt.Printf("Server>%s\tClock: %d\t%d\n", res[0], vectorClocks.clock[0], vectorClocks.clock[1]) // Printing client messages
		}
	}
	flag <- true
}

func send2Server(conn *net.Conn) {
	for {
		clMsg, err := bufio.NewReader(os.Stdin).ReadString('\n') //Reading from Stdin
		if err != nil {
			fmt.Println("Error while reading text from client:", err.Error())
		} else {
			vectorClocks.Lock()
			vectorClocks.clock[1] += 1
			vectorClocks.Unlock()
			var sendVal = clMsg + strconv.Itoa(vectorClocks.clock[1]) + "\t"			
			fmt.Fprintf(*conn, sendVal)
		}
	}
	flag <- true
}