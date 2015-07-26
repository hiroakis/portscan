package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPortLowerLimit := flag.Int("lower", 1, "scan range(lower limit)")
	optPortUpperLimit := flag.Int("upper", 63556, "scan range(upper limit)")
	optProtocol := flag.String("protocol", "tcp", "tcp or udp")
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
			_, err := net.Dial(*optProtocol, target)

			if err != nil && strings.Index(err.Error(), "connection refused") != -1 {
				closedPort <- port
			} else if err != nil {
				errorMessage <- err.Error()
			} else {
				openedPort <- port
			}
		}(i)
	}

	for i := *optPortLowerLimit; i <= *optPortUpperLimit; i++ {
		select {
		case o := <-openedPort:
			if *optPrint == "open" || *optPrint == "both" {
				msg := fmt.Sprintf("%s://%s:%d open", *optProtocol, *optHost, o)
				fmt.Println(msg)
			}
		case c := <-closedPort:
			if *optPrint == "close" || *optPrint == "both" {
				msg := fmt.Sprintf("%s://%s:%d closed", *optProtocol, *optHost, c)
				fmt.Println(msg)
			}
		case e := <-errorMessage:
			fmt.Println(e)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
