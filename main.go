package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPortLowerLimit := flag.Int("lower", 1, "scan range(lower limit)")
	optPortUpperLimit := flag.Int("upper", 1024, "scan range(upper limit)")
	optPrint := flag.String("print", "open", "open, close or both")
	flag.Parse()

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	if *optPrint != "open" && *optPrint != "close" && *optPrint != "both" {
		fmt.Println("-port must be open, close or both")
		os.Exit(1)
	}

	openedPort := make(chan int)
	closedPort := make(chan int)
	errorMessage := make(chan string)

	for i := *optPortLowerLimit; i <= *optPortUpperLimit; i++ {
		go func(port int) {
			target := fmt.Sprintf("%s:%d", *optHost, port)
			con, err := net.DialTimeout("tcp", target, 3*time.Second)

			if err != nil && strings.Index(err.Error(), "connection refused") != -1 {
				closedPort <- port
			} else if err != nil && strings.Index(err.Error(), "i/o timeout") != -1 {
				errorMessage <- err.Error()
			} else if err != nil {
				errorMessage <- err.Error()
			} else {
				defer con.Close()
				openedPort <- port
			}
		}(i)
	}

	for i := *optPortLowerLimit; i <= *optPortUpperLimit; i++ {
		select {
		case o := <-openedPort:
			if *optPrint == "open" || *optPrint == "both" {
				msg := fmt.Sprintf("tcp://%s:%d open", *optHost, o)
				fmt.Println(msg)
			}
		case c := <-closedPort:
			if *optPrint == "close" || *optPrint == "both" {
				msg := fmt.Sprintf("tcp://%s:%d closed", *optHost, c)
				fmt.Println(msg)
			}
		case e := <-errorMessage:
			fmt.Println(e)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
