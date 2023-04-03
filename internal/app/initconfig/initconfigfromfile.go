package initconfig

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// JSON input record.
type ConfigJSONrecord struct {
	ServerAddress    string `json:"server_address"`
	BaseURL          string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN      string `json:"database_dsn"`
	EnableHTTPS      string `json:"enable_https"`
}

// config JSON input.
var configJSON ConfigJSONrecord

// addInitVarsFromConfigFile - adds Init Vars From Config File if they haven't been defined yet.
func addInitVarsFromConfigFile() {
	srvConfigFileENV, srvConfigFileexists := os.LookupEnv("CONFIG")
	srvConfigFile := ""
	if !srvConfigFileexists {
		if isFlagPassed("c") {
			srvConfigFile = *SrvConfigFileflag
			fmt.Println("Set from flag: CONFIG:", srvConfigFile)
		} else {
			fmt.Print("CONFIG file: not set, no flag, no ENV")
		}
	} else {
		srvConfigFile = srvConfigFileENV
		fmt.Println("Set from ENV: CONFIG file:", srvConfigFile)
	}

	if srvConfigFile != "" {
		configFile, err := os.OpenFile(srvConfigFile, os.O_RDONLY, 0777)
		if err != nil {
			fmt.Println("File does NOT EXIST")
			fmt.Println(err)
		} else {
			scanner := bufio.NewScanner(configFile)
			JSONstring := ""
			for scanner.Scan() {
				JSONstring += scanner.Text()
			}
			err = json.Unmarshal([]byte(JSONstring), &configJSON)
			if err != nil {
				panic(err)
			}
			if SrvAddress == "" {
				SrvAddress = configJSON.ServerAddress
			}
			if BaseURL == "" {
				BaseURL = configJSON.BaseURL
			}
			if FileDBpath == "" {
				FileDBpath = configJSON.FileStoragePath
			}
			if PostgresDBURL == "" {
				PostgresDBURL = configJSON.DatabaseDSN
			}
			if SrvRunHTTPS == "" {
				SrvRunHTTPS = configJSON.EnableHTTPS
			}
		}
	}
}
