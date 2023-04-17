// Package initconfig declares all global variables and init functions for service.
package initconfig

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Variables using across the service.
var (
	// FileDBpath, BaseURL, SrvAddress - store database options. SrvRunHTTPS - store server HTTPS mode.
	FileDBpath, BaseURL, SrvAddress, SrvRunHTTPS, SrvConfigFile, TrustedSubnet string
	// SrvAddressflag, BaseURLflag, FileDBpathflag, PostgresDBURLflag, SrvRunHTTPSflag - store possible flags.
	SrvAddressflag, BaseURLflag, FileDBpathflag, PostgresDBURLflag, SrvConfigFileflag, TrustedSubnetflag *string
	SrvRunHTTPSflag                                                                                      *bool
)

// PostgresDBURL - init database URL string.
var PostgresDBURL string

// InputIDstoDel - init channel with ids to delete.
var InputIDstoDel = make(chan string, 7)

// WG - init waitgroup.
var WG sync.WaitGroup

// InitFlags - init flags.
func InitFlags() {

	SrvAddressflag = flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS flag")
	BaseURLflag = flag.String("b", "http://127.0.0.1:8080", "BASE_URL flag")
	PostgresDBURLflag = flag.String("d", "postgres://postgres:1@localhost:5432/gotoschool", "DATABASE_DSN flag")
	FileDBpathflag = flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH flag")
	SrvRunHTTPSflag = flag.Bool("s", false, "ENABLE_HTTPS flag")
	SrvConfigFileflag = flag.String("c", "", "CONFIG flag")
	TrustedSubnetflag = flag.String("t", "", "TRUSTED_SUBNET flag")
}

// SetinitVars - init global vars according to ENV vars and flags passed.
func SetinitVars() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	baseURLENV, baseURLexists := os.LookupEnv("BASE_URL")
	srvAddressENV, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	fileDBpathENV, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")
	postgresDBURLENV, postgresDBURLexists := os.LookupEnv("DATABASE_DSN")
	srvRunHTTPSENV, srvRunHTTPSexists := os.LookupEnv("ENABLE_HTTPS")
	trustedSubnetENV, trustedSubnetexists := os.LookupEnv("TRUSTED_SUBNET")

	if !srvAddressexists {
		SrvAddress = *SrvAddressflag
		fmt.Println("Set from flag: SrvAddress:", SrvAddress)
	} else {
		SrvAddress = srvAddressENV
		fmt.Println("Set from ENV: SrvAddress:", SrvAddress)
	}

	if !baseURLexists {
		BaseURL = *BaseURLflag
		fmt.Println("Set from flag: BaseURL:", BaseURL)
	} else {
		BaseURL = baseURLENV
		fmt.Println("Set from ENV: BaseURL:", BaseURL)
	}

	if !postgresDBURLexists {
		if isFlagPassed("d") {
			PostgresDBURL = *PostgresDBURLflag
			fmt.Println("Set from flag: PostgresDBURL:", PostgresDBURL)
		} else {
			fmt.Print("DATABASE_DSN: not set, no flag, no ENV")
		}

	} else {
		PostgresDBURL = postgresDBURLENV
		fmt.Println("Set from ENV: PostgresDBURL:", PostgresDBURL)
	}

	if !fileDBpathexists {

		/*var flagFound bool
		//слайс аргументов
		parameters := os.Args[1:]
		//среди них ищем передан ли -f
		flagFound = false
		fmt.Printf("Parameters: %v\n", parameters)
		for i, n := range parameters {
			//это как-то некрасиво:
			argSplitted := strings.Split(n, "=")
			fmt.Println(argSplitted[0] + " знак равно " + argSplitted[1])
			if argSplitted[0] == "-f" {
				flagFound = true
			}
			fmt.Println(flagFound)
			i++
		}*/

		//if flagFound{
		if isFlagPassed("f") {
			//FileDBpathflag := flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH")
			FileDBpath = *FileDBpathflag
			fmt.Println("Set from flag: FileDBpath:", FileDBpath)
		} else {
			fmt.Print("FILE_STORAGE_PATH: not set, no flag, no ENV")
		}
	} else {
		FileDBpath = fileDBpathENV
		fmt.Println("Set from ENV: FileDBpath:", FileDBpath)
	}
	if !srvRunHTTPSexists {
		if *SrvRunHTTPSflag {
			SrvRunHTTPS = "HTTPS_mode_enabled"
			fmt.Print("Set from flag: SrvRunHTTPS:", *SrvRunHTTPSflag)
		} else {
			fmt.Print("ENABLE_HTTPS: not set ")
		}
	} else {
		SrvRunHTTPS = srvRunHTTPSENV
		fmt.Println("Set from ENV: SrvRunHTTPS:", SrvRunHTTPS)
	}
	if !trustedSubnetexists {
		if isFlagPassed("t") {
			TrustedSubnet = *TrustedSubnetflag
			fmt.Print("Set from flag: TrustedSubnet:", *TrustedSubnetflag)
		} else {
			fmt.Print("TRUSTED_SUBNET: not set ")
		}
	} else {
		TrustedSubnet = trustedSubnetENV
		fmt.Println("Set from ENV: TrustedSubnet:", TrustedSubnet)
	}

	addInitVarsFromConfigFile()

}

// isFlagPassed - checks if flag is passed
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
