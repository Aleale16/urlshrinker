package main

import (
	"flag"

	"urlshrinker/internal/app/initconfig"
	"urlshrinker/internal/app/server"
	//"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	//"github.com/Aleale16/urlshrinker/internal/app/server"
)

func main() {

	//staticlint.Runinc18Checks()

	initconfig.InitFlags()

	flag.Parse()

	initconfig.SetinitVars()

	server.Start()

}
