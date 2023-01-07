package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLrecord map[string]string
type Userrecord map[string][]string
type URLJSONrecord struct {
    ID string `json:"id"`
    FullURL string `json:"fullurl"`    
}

type UsrURLJSONrecord struct {
    ShortURL string `json:"short_url"`
    FullURL string `json:"original_url"`    
}

var URL URLrecord
var Usr Userrecord
var dbPath string
var RAMonly, PGdbOpened, dbPathexists bool
var onlyOnce sync.Once

var PGdb *pgxpool.Pool

func InitPGdb() {
	
//----------------------------//
//Подключаемся к СУБД postgres
//----------------------------//
	//urlExample := "postgres://postgres:1@localhost:5432/gotoschool"
    //os.Setenv("DATABASE_DSN", urlExample)

    poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_DSN:", err)
	}
    fmt.Println(poolConfig)

    PGdb, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
    if err != nil {
		PGdbOpened = false
        fmt.Println("ERROR! PGdbOpened = false")
        panic(err)
    } else {
		PGdbOpened = true
	}
    
}

func Initdb() {
	dbPath, dbPathexists = os.LookupEnv("FILE_STORAGE_PATH")
	if (dbPathexists && dbPath != "") {
		RAMonly = false
		log.Println("Loading DB file...")
		log.Println("dbPath: " + dbPath)
		log.Println("Copying DB file to RAM storage...")
		URL = make(URLrecord)
		URL = copyFiletoRAM(dbPath, URL)

		} else {
			RAMonly = true
			fmt.Println("DB file path is not set in env vars! Loading RAM storage...")
			URL = make(URLrecord)
	}	
	fmt.Println("Storage ready!")
	Usr = make(Userrecord)

//???Возможно в отдельный пакет инициализацию Postgres надо вынести?
	InitPGdb()
}

func copyFiletoRAM(dbPath string, URLs URLrecord) URLrecord{
	DBfile, err := os.OpenFile(dbPath, os.O_RDONLY, 0777)
	if err != nil {
		log.Println("File does NOT EXIST")
		//result = ""
		log.Println(err)
		//idIsnew = false
		//panic(err)
	} else {
		scanner := bufio.NewScanner(DBfile)
		var postJSON URLJSONrecord
		var lastID int
		line := 0
		id := initconfig.NextID
		for scanner.Scan(){
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
				URLs[postJSON.ID] = postJSON.FullURL
				log.Println("Line " + strconv.Itoa(line) + "is loaded to RAM: " + scanner.Text())				
			}
			line++
		}
		if postJSON.ID != ""{
			lastID, _ = strconv.Atoi(postJSON.ID)
			initconfig.NextID = lastID + initconfig.Step
		}
		id = id + initconfig.Step
		initconfig.NextID = id
	}	
	DBfile.Close()	
	return URLs
}

func AssignShortURLtouser(userid string, shortURLid string){
	onlyOnce.Do(Initdb)
	uid := userid
	Usr[uid] = append(Usr[uid], shortURLid)
	fmt.Println("AssignShortURLtouser: " + string(uid)+ " shortURLid= " )	
	fmt.Println(Usr[uid])

	if PGdbOpened {
		_, err := PGdb.Exec(context.Background(), `insert into users(uid, shortid) values ($1, $2)`, uid, shortURLid)
		if err == nil {
			log.Println("w.WriteHeader(http.StatusOK)")
		} else {
			log.Println("http.Error(w, "+"Internal server error"+", http.StatusInternalServerError)")
		}
	}
}

func CheckPGdbConn() (connected bool){
	onlyOnce.Do(Initdb)
	//defer PGdb.Close()
    err := PGdb.Ping(context.Background())
    if err != nil {
        fmt.Println(err)
		return false
    } else {
        fmt.Println("Ping db is ok")
		return true
    }	
}

func GetuserURLS(userid string) (output string, noURLs bool){
	var UsrURLJSON []UsrURLJSONrecord
	var JSONresult []byte
	noURLs = true 
	UsrShortURLs := Usr[userid]
	if len(UsrShortURLs)>0{
		for _, v := range UsrShortURLs {	
				log.Println(v)
				//Так нормально заполнять JSON перед маршаллингом?
				UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
					ShortURL:	initconfig.BaseURL + "/?id=" + v,
					FullURL:	URL[v],
				})				
		}
		JSONdata, err := json.Marshal(&UsrURLJSON)
		if err != nil {
			return err.Error(), noURLs
		}
		//JSONdata = append(JSONdata, '\n')
		//URL[id] = string(JSONdata)
		JSONresult = JSONdata
		log.Println("JSONresult= ")		
		log.Println(JSONresult)		
		log.Println(string(JSONresult))		
		noURLs = false
	}
	

	if PGdbOpened {
		var (
			UID string
			shortID string
			FullURL string
		)
		rows, err := PGdb.Query(context.Background(), "SELECT usr.uid, usr.shortid, urls.fullurl FROM users as usr LEFT JOIN urls ON urls.shortid = usr.shortid where uid=$1", userid)
		if err != nil {
			return err.Error(), noURLs
		}
		// обязательно закрываем перед возвратом функции
		defer rows.Close()

		// пробегаем по всем записям
		for rows.Next() {
			err := rows.Scan(&UID, &shortID, &FullURL)
			if err != nil {
				log.Fatal(err)
			}
			UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
				ShortURL:	initconfig.BaseURL + "/?id=" + shortID,
				FullURL:	FullURL,
			})	
		}
		JSONdata, err := json.Marshal(&UsrURLJSON)
		if err != nil {
			return err.Error(), noURLs
		}
		//JSONdata = append(JSONdata, '\n')
		//URL[id] = string(JSONdata)
		JSONresult = JSONdata
		log.Println("JSONresult= ")		
		log.Println(JSONresult)		
		log.Println(string(JSONresult))		
		noURLs = false
	}

	return string(JSONresult), noURLs
}

func Storerecord(fullURL string) string{
	onlyOnce.Do(Initdb)
	//id := strconv.Itoa(rand.Intn(9999))
	id := strconv.Itoa(initconfig.NextID)
	
	/*for (!isnewID(id)){
		id = strconv.Itoa(rand.Intn(9999))
	}*/



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
		//URL[id] = string(JSONdata)
		URL[id] = fullURL

		DBfile, _ := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_APPEND , 0777)
		_, err = DBfile.Write(JSONdata)	
		if err != nil {	
			return err.Error()
		}
		DBfile.Close()

	}
	
	if PGdbOpened {
		_, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl) values ($1, $2)`, id, fullURL)
		if err == nil {
			log.Println("w.WriteHeader(http.StatusOK)")
		} else {
			log.Println("http.Error(w, "+"Internal server error"+", http.StatusInternalServerError)")
		}
	}

	initconfig.NextID = initconfig.NextID + initconfig.Step
	return id
}

func Getrecord(id string) string {
	var result string
	onlyOnce.Do(Initdb)

	if RAMonly {
		result = URL[id]
	} else {
		result = URL[id]
		/*
		DBfile, err := os.OpenFile(dbPath, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("File does NOT EXIST")
			result =""
			log.Println(err)
			//idIsnew = false
			//panic(err)
		} else {
			scanner := bufio.NewScanner(DBfile)
			line :=0
			var postJSON URLJSONrecord
			idIsfound := false
			for scanner.Scan() && !idIsfound{
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
						idIsfound = true
						log.Println("ID exists: " + postJSON.ID + "; FullURL: " + postJSON.FullURL)
					}
					line++
				}
			}
			result = postJSON.FullURL
			DBfile.Close()
		}*/
	}
	if PGdbOpened {
		var (
			FullURL string
		)
		rows, err := PGdb.Query(context.Background(), "SELECT urls.fullurl FROM urls where shortid=$1", id)
		if err != nil {
			return err.Error()
		}
		// обязательно закрываем перед возвратом функции
		defer rows.Close()

		// пробегаем по всем записям
		for rows.Next() {
			err := rows.Scan(&FullURL)
			if err != nil {
				log.Fatal(err)
			}
		}
		result = id
	}
	

	if (result != ""){
		return result
	} else {
		return "http://google.com/404"
		
	}
}
/*
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
			log.Println("File does NOT EXIST")
			log.Println(err)
			//idIsnew = false
			//panic(err)
		} else {
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
		}
		return idIsnew
	}
}*/