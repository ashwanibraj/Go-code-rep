//Check out server file for comments on code. 
//This file is similar to the server file, as far as the syntax goes.
//Make sure you change the IP address in this file to the ip address of server.
//After starting the server, start the client. Command is "go run Chat_client.go"

package main
//Importing a bunch of libraries
import (
	"fmt"
	"bufio"
	"net"
	"os"
)

var flag = make(chan bool) //initializing a channel variable of boolean type

func main() {
	conn, err := net.Dial("tcp", "10.136.2.42:8080") //Initiating a new connection with the server at given ip and port using TCP
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
		serMsg, err := bufio.NewReader(*conn).ReadString('\n') //Reading from open connection
		if err != nil {
			fmt.Println("Error while reading text:", err.Error())			
			break
		} else {
			fmt.Println("Server>", string(serMsg))
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
			fmt.Fprintf(*conn, clMsg)
		}
	}
	flag <- true
}