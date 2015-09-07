//Transaction Processing System - Basic Implementation with example
//Taking lessons from the Producer Consumer Problem
/*The initialization of allAccnts(all account numbers with balances), src(source account no.) and transAccAmt(Dest account no and amount)
is just for explanatory purpose. For application, src and transAccAmt can be read from a form entry in HTML*/
package main

import (
	"fmt"
	"sync"
)

//All accounts, mapping Account Numbers with Account balances
//Maps are not safe in Go for concurrent write and read. Using Mutex as advised in Go docs.
//Refer to Concurrency in link https://blog.golang.org/go-maps-in-action
var allAccnts = struct{
	sync.RWMutex
	m map[int]int
}{m: make(map[int]int)}

var flag = make(chan bool)
var srcAccnt = make(chan int)
var destAccnts = make(chan map[int]int)

func main() {
	//Initial assignment
	allAccnts.Lock() //Write lock to initialize allAccounts
	allAccnts.m = map[int]int{
	//AccntNo. : Balance	
		1111 : 90000,
		2222 : 700,
		1212 : 80,
		3333 : 9,
		1313 : 100,
		2323 : 890,
	}
	allAccnts.Unlock() //Release write lock
	allAccnts.RLock() 
	for accnt, bal := range allAccnts.m {
		fmt.Println("AccountNo:",accnt, " Balance:", bal)
	}
	allAccnts.RUnlock()
	fmt.Println("--------------------------------------------------------------------------------------")
	go callerUpd()
	go transferAgent()
	<-flag
	allAccnts.RLock()
	for accnt, bal := range allAccnts.m {
		fmt.Println("AccountNo:",accnt, " Balance:", bal)
	}
	allAccnts.RUnlock()
}

func callerUpd() {
	src := 1111 //Source Account Number for amount transfer
	sumChk := 0 //Variable to check if total amount to be transfered exceeds account bal of source
	fmt.Println("Source Account No.:", src)	
	srcAccnt <-src //Passing src over the channel
	
	transAccAmt := map[int]int{
//Destination accounts : Amount to be transferred	
		2222 : 10,
		3333 : 1000,
		1313 : 90,
		//5555 : 10,//This entry is for testing validation of accounts. Uncomment and run to get error.
	}

	fmt.Println("Destination accounts and amounts to be transfered:")
	
	for d, t := range transAccAmt { //Looping through all destination accounts
		if _, ok := allAccnts.m[d]; ok { // Checking if dest account exists in allAccnts or not.
			fmt.Println("Account No:",d," Amount:",t)
			sumChk := sumChk + t	
			if (allAccnts.m[src] - sumChk) < 0 {
				fmt.Println("Insufficient funds.")
				flag <- true
			}	
		} else { //Printing error and returning.
			fmt.Printf("Destination account %d does not exist. Abort transaction.\n", d)
			flag <- true
		}
			
	}	
	
	fmt.Println("--------------------------------------------------------------------------------------")	
	
	destAccnts <- transAccAmt
	flag <- true
}

func transferAgent() {
	//Always running to check for any open channels
	for {
		src1 := <- srcAccnt 
		trans := <- destAccnts 
		
		for acc, amt := range trans { //Looping for all accounts and amounts mentioned
			allAccnts.RLock()
			if (allAccnts.m[src1] - amt) < 0 { //Checking if source has enough funds
				allAccnts.RUnlock()
				fmt.Println("Insufficient funds")				
			} else {
				allAccnts.RUnlock()
				allAccnts.Lock() //Write lock on allAccnts
				allAccnts.m[src1] = allAccnts.m[src1] - amt //Debit from Source account
				allAccnts.m[acc] = allAccnts.m[acc] + amt   //Credit to destination account		
				allAccnts.Unlock() //Write lock released
			}
		}		
		fmt.Println("Transfer Successful!")		
	}
}