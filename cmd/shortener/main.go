//go:build !darwin

package main

import (
	"flag"
	"fmt"
	"log"

	"urlshrinker/internal/app/initconfig"
	"urlshrinker/internal/app/server"
	//"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	//"github.com/Aleale16/urlshrinker/internal/app/server"
)

// Init global initbuildinfo strings
// Output:
// Build version: <buildVersion> (or "N/A" if no value)
// Build date: <buildDate> (or "N/A" if no value)
// Build commit: <buildCommit> (or "N/A" if no value)
var (
	// buildVersion - global buildVersion value.
	buildVersion  string
	// buildDate - global buildDate value.
	buildDate  string
	// buildCommit - buildCommit value.
	buildCommit  string
)

func main() {


	// fmt.Printf("Build version: %s", initconfig.BuildVersion)
	// fmt.Printf("Build date: %s", initconfig.BuildDate)
	// fmt.Printf("Build commit: %s", initconfig.BuildCommit)
	if buildVersion == "" {buildVersion = "N/A"}
	if buildDate == "" {buildDate = "N/A"}
	if buildCommit == "" {buildCommit = "N/A"}
	_, err := fmt.Printf("Build version: %s", buildVersion)
	if err!=nil {
		log.Print(err)
	}
	_, err = fmt.Printf("Build date: %s", buildDate)
	if err!=nil {
		log.Print(err)
	}
	_, err = fmt.Printf("Build commit: %s", buildCommit)
	if err!=nil {
		log.Print(err)
	}

	//staticlint.Runinc18Checks()

	initconfig.InitFlags()

	flag.Parse()

	initconfig.SetinitVars()

	server.Start()

	//os.Exit(10)

}
