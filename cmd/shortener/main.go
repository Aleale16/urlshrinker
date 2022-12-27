package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"github.com/Aleale16/urlshrinker/internal/app/server"
)

func handler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
}

func main() {

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	exitchnl := make(chan int)

	go func() {
		for {
			s := <-sigchnl
			handler(s)
		}
	}()

	exitcode := <-exitchnl
	os.Exit(exitcode)
	
	initconfig.InitFlags()

	flag.Parse()	

	initconfig.SetinitVars()
	
	server.Start()
	
}
