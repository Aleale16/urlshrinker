package initconfig

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// JSON input record.
type ConfigJSONrecord struct {
	SERVER_ADDRESS    string `json:"server_address"`
	BASE_URL          string `json:"base_url"`
	FILE_STORAGE_PATH string `json:"file_storage_path"`
	DATABASE_DSN      string `json:"database_dsn"`
	ENABLE_HTTPS      string `json:"enable_https"`
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
				SrvAddress = configJSON.SERVER_ADDRESS
			}
			if BaseURL == "" {
				BaseURL = configJSON.BASE_URL
			}
			if PostgresDBURL == "" {
				PostgresDBURL = configJSON.DATABASE_DSN
			}
			if FileDBpath == "" {
				FileDBpath = configJSON.FILE_STORAGE_PATH
			}
			if SrvRunHTTPS == "" {
				SrvRunHTTPS = configJSON.ENABLE_HTTPS
			}
		}
	}
}
