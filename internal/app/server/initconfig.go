package server

import (
	"flag"
	"fmt"
	"os"
)
var FileDBpath, BaseURL, SrvAddress string

func InitFlags() {

	baseURLENV, baseURLexists := os.LookupEnv("BASE_URL")
	srvAddressENV, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	fileDBpathENV, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")

	if !srvAddressexists{
		SrvAddressflag := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS")
		SrvAddress = *SrvAddressflag
		fmt.Println("Set from flag: SrvAddress:", SrvAddress)
	} else {
		SrvAddress = srvAddressENV
		fmt.Println("Set from ENV: SrvAddress:", SrvAddress)
	}

	if !baseURLexists{
		BaseURLflag := flag.String("b", "http://127.0.0.1:8080", "BASE_URL")
		BaseURL = *BaseURLflag
		fmt.Println("Set from flag: BaseURL:", BaseURL)
	} else {
		BaseURL = baseURLENV
		fmt.Println("Set from ENV: BaseURL:", BaseURL)
	}

	if !fileDBpathexists{
		var flagFound bool
		//слайс аргументов
		parameters := os.Args[1:]
		//среди них ищем передан ли -f
		flagFound = false
		fmt.Printf("Parameters: %v\n", parameters)
		for i, n := range parameters {
			if n == "-f" {
				flagFound = true
			}
			i++
		}
		if flagFound {
			FileDBpathflag := flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH")
			FileDBpath = *FileDBpathflag
			fmt.Println("Set from flag: FileDBpath:", FileDBpath)
		} else {
			fmt.Print("FILE_STORAGE_PATH: not set") 
		}
	} else {
		FileDBpath = fileDBpathENV
		fmt.Println("Set from ENV: FileDBpath:", FileDBpath)
	}
	

}