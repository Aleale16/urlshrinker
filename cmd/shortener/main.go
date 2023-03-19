//go:build !darwin

package main

import (
	"flag"
	"fmt"

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
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)

	//staticlint.Runinc18Checks()

	initconfig.InitFlags()

	flag.Parse()

	initconfig.SetinitVars()

	server.Start()

}
