//Simple implementation of Producer Consumer Problem
//My first Go program with concurrency implementation
package main

import ("fmt")

var flag = make(chan bool)
var msg = make(chan int)

func main() {
	go producer()
	go consumer()
	<-flag
}

func producer() {
	for i := 0; i < 10; i++ {
		msg <- i
	}
	flag <- true
}

func consumer() {
	for {
		msg1 := <- msg
		fmt.Println(msg1)
	}
}