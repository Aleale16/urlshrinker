package server

import (
	"flag"
	"fmt"
	"os"
	"strings"
)
var FileDBpath, BaseURL, SrvAddress string
var SrvAddressflag, BaseURLflag, FileDBpathflag *string
func InitFlags() {	
	SrvAddressflag = flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS flag")
	BaseURLflag = flag.String("b", "http://127.0.0.1:8080", "BASE_URL flag")
	FileDBpathflag = flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH flag")
}
	

func SetinitVars() {

	baseURLENV, baseURLexists := os.LookupEnv("BASE_URL")
	srvAddressENV, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	fileDBpathENV, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")

	if !srvAddressexists{		
		SrvAddress = *SrvAddressflag
		fmt.Println("Set from flag: SrvAddress:", SrvAddress)
	} else {
		SrvAddress = srvAddressENV
		fmt.Println("Set from ENV: SrvAddress:", SrvAddress)
	}

	if !baseURLexists{		
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
			//это как-то некрасиво:
			argSplitted := strings.Split(n, "=")
			fmt.Println(argSplitted[0] + " знак равно " + argSplitted[1])
			if argSplitted[0] == "-f" {
				flagFound = true
			}
			fmt.Println(flagFound)
			i++
		}
		if flagFound {
			//FileDBpathflag := flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH")
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