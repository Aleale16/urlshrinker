package main

import (
	"flag"

	"github.com/Aleale16/urlshrinker/internal/app/server"
)

func main() {
	server.InitFlags()

	flag.Parse()	

	server.SetinitVars()
	
	server.Start()
	
}
