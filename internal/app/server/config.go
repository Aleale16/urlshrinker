package server

import (
	"flag"
	"fmt"
	"os"
)

func InitFlags() {
	_, baseURLexists := os.LookupEnv("BASE_URL")
	_, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	_, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")

	if !srvAddressexists{
		SrvAddress := flag.String("a", "localhost:8080", "SERVER_ADDRESS")
		fmt.Println("Set from flag: SrvAddress:", *SrvAddress)
	}

	if !baseURLexists{
		BaseURL := flag.String("b", "http://localhost:8080", "BASE_URL")
		fmt.Println("Set from flag: BaseURL:", *BaseURL)
	}

	if !fileDBpathexists{
		fileDBpath := flag.String("f", "../../internal/app/storage/database.txt", "FILE_STORAGE_PATH")
		fmt.Println("Set from flag: fileDBpath:", *fileDBpath)
	}

}