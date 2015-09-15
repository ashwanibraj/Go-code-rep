/*Multithreading example- Running unix command in threads to print result for dirname and corresponding file-count 
for the current directory and all subdirectories recursively. Four goroutines are launched.
Copy and paste this in a file, say "Concurrent_unix_cmd.go".
Run the code in terminal with "go run Concurrent_unix_cmd.go"
*/

package main

import (
    "os/exec"
    "fmt"
    "sync"
    "strings"
)

func main() {

    //Defining a type struct
    type vars struct{ x *exec.Cmd //For external command execution
                      y string
    }

    tasks := make(chan vars, 64)

    // spawn four worker goroutines
    var wg sync.WaitGroup
    for i := 0; i < 4; i++ {
        wg.Add(1)
        go func() { // Begin of goroutine definition
            for fin1 := range tasks {
                fin, err1 := fin1.x.Output() //Execution the unix command and taking the solution in fin1, error in err
                if err1 != nil {
                    fmt.Println(fin1)                    
                } else {
                    sol := string(fin) 
                    result := strings.Fields(sol)
                    for ind, val := range result {
                        fmt.Println(fin1.y, " ", val)
                        ind += 1
                    }
                }
            }
            wg.Done()
        }()
    }
    
    cm := exec.Command("bash", "-c", "find . -type d") //Constructing command to find out all sub directories from current dir
    out, err := cm.Output() //Executing the unix command
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        outf := string(out)
        
        res := strings.Fields(outf)
        for index, each := range res {            
            s1 := fmt.Sprint("ls -1 ",each,"/ | wc -l")
            index += 1 
            //Initializing a variable of type vars with external command and directory name
            t1 := vars{x:exec.Command("bash", "-c", s1), y:each} 
            tasks <- t1 //Passing struct over the channel for goroutines
        }
        
        close(tasks)

        wg.Wait()
    }
}