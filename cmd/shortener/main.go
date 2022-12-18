package main

import (
	"github.com/Aleale16/urlshrinker/internal/app/server"
)

func main() {
	server.InitFlags()

	

	server.Start()
	
}
