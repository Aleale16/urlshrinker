package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type URLrecord map[string]string
type URLJSONrecord struct {
    ID string `json:"id"`
    FullURL string `json:"fullurl"`    
}

var URL URLrecord
var dbPath string
var RAMonly, dbPathexists bool
var onlyOnce sync.Once

func Initdb() {
	dbPath, dbPathexists = os.LookupEnv("FILE_STORAGE_PATH")
	if (dbPathexists && dbPath != "") {
		RAMonly = false
		log.Println("Loading DB file...")
		log.Println("dbPath: " + dbPath)
		} else {
			RAMonly = true
			fmt.Println("DB file path is not set in env vars! Loading RAM storage...")
			URL = make(URLrecord)
	}	
	fmt.Println("Storage ready!")
}

func Storerecord(fullURL string) string{
	onlyOnce.Do(Initdb)
	id := strconv.Itoa(rand.Intn(9999))
	
	for (!isnewID(id)){
		id = strconv.Itoa(rand.Intn(9999))
	}

	if RAMonly {
		URL[id] = fullURL
	} else {
		URLJSONline := URLJSONrecord{
			ID:			id,
        	FullURL:	fullURL,
		}
		JSONdata, err := json.Marshal(&URLJSONline)
		if err != nil {
			return err.Error()
		}
		JSONdata = append(JSONdata, '\n')
		

		DBfile, _ := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_APPEND , 0777)
		_, err = DBfile.Write(JSONdata)	
		if err != nil {	
			return err.Error()
		}
		DBfile.Close()
	}
	return id
}

func Getrecord(id string) string {
	onlyOnce.Do(Initdb)
	result := URL[id]

	if (result != ""){
		return result
	} else {
		return "http://google.com/404"
	}
}

func isnewID(id string) bool{
	if RAMonly {
	result := URL[id]
	if (result == ""){
		return true
	} else {return false}
	}else {
		var idIsnew bool
		idIsnew = true
		DBfile, err := os.OpenFile(dbPath, os.O_RDONLY, 0777)
		if err != nil {
			//log.Println(err)
			//idIsnew = false
			panic(err)
		}
		scanner := bufio.NewScanner(DBfile)
		line :=0
		var postJSON URLJSONrecord

		for scanner.Scan() && idIsnew{
			//log.Println(line)
			//log.Println("lineStr: " + scanner.Text())
			if scanner.Text() != "" {
				err = json.Unmarshal([]byte(scanner.Text()), &postJSON)
				if err != nil {
					panic(err)
				}
			//отладка что было в поле FullURL в строке файла
				log.Println(postJSON.ID)
				log.Println(postJSON.FullURL)
				if postJSON.ID == id {
					idIsnew = false
					log.Println("ID exists: " + postJSON.ID)
				}
				line++
			}
		}
		DBfile.Close()
		return idIsnew
	}
}