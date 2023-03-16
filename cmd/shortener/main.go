package main

import (
	"flag"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"github.com/Aleale16/urlshrinker/internal/app/server"
)

func main() {

	initconfig.InitFlags()

	flag.Parse()

	initconfig.SetinitVars()

	server.Start()

}
