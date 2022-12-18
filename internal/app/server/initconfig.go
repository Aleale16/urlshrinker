package server

import (
	"flag"
	"fmt"
	"os"
)
var FileDBpath, BaseURL, SrvAddress string

func InitFlags() {

	_, baseURLexists := os.LookupEnv("BASE_URL")
	_, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	_, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")

	if !srvAddressexists{
		SrvAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS")
		fmt.Println("Set from flag: SrvAddress:", *SrvAddress)
	}

	if !baseURLexists{
		BaseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL")
		fmt.Println("Set from flag: BaseURL:", *BaseURL)
	}

	if !fileDBpathexists{
		FileDBpath := flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH")
		fmt.Println("Set from flag: FileDBpath:", *FileDBpath)
		os.Setenv("FILE_STORAGE_PATH", *FileDBpath)
	}

}