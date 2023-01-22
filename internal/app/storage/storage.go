package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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

//var onlyOnce sync.Once

var PGdb *pgxpool.Pool

//Будем фиксировать тот тип базы данных, который удалось подключить, с которым будем работать. Для каждого типа имплементируем методы интерфейса Storager в файле ifstorager.go
type connectRAM struct {}
type connectFileDB struct {}
type connectPGDB struct {}
// Если вот так не объявить переменные, то не проходит статический тест, подчеркивает желтым, говорит, что такие типы не используются
var DataBaseconnectRAM connectRAM
var DataBaseconnectFileDB connectFileDB
var DataBaseconnectPGDB connectPGDB
//Переменная типа базы данных, которую будем использовать:
//var DataBase struct{}
var S storager

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
			dbPathexists = false
			log.Println("DB file path is not set in env vars! Loading RAM storage...")
			URL = make(URLrecord)
	}	
	log.Println("Storage ready!")
	Usr = make(Userrecord)

	InitPGdb()
//Фиксируем тот тип БД, который получилось включить
	SetdbType()
}

func InitPGdb() {
	
//----------------------------//
//Подключаемся к СУБД postgres
//----------------------------//
	//urlExample := "postgres://postgres:1@localhost:5432/gotoschool"
    //os.Setenv("DATABASE_DSN", urlExample)
	//initconfig.PostgresDBURL = urlExample
	var DBLastURLID, DBLastUID string
	PGdbOpened = false
	if initconfig.PostgresDBURL != "" {
		poolConfig, err := pgxpool.ParseConfig(initconfig.PostgresDBURL)
		if err != nil {
			log.Fatalln("Unable to parse DATABASE_DSN:", err)
		}
		//fmt.Println(poolConfig)

		PGdb, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			PGdbOpened = false
			fmt.Println("ERROR! PGdbOpened = false")
			panic(err)
		} else {
			_, err := PGdb.Exec(context.Background(), `CREATE SEQUENCE IF NOT EXISTS public.urls_id_seq
				INCREMENT 1
				START 1
				MINVALUE 1
				MAXVALUE 2147483647
				CACHE 1;
		
				ALTER SEQUENCE urls_id_seq
					OWNER TO postgres;
		
				CREATE TABLE IF NOT EXISTS public.urls
				(
					shortid character varying(10) ,
					fullurl character varying(1000) ,
					id integer NOT NULL DEFAULT nextval('urls_id_seq'::regclass),
					uid character varying(10) DEFAULT 0,
					active boolean,
					uidint integer DEFAULT 0,
					CONSTRAINT urls_pkey PRIMARY KEY (id)
				);
				ALTER TABLE public.urls ADD UNIQUE (fullurl)`)
			if err != nil {
				log.Println(err)
			} 
			
			_, err = PGdb.Exec(context.Background(), `CREATE SEQUENCE IF NOT EXISTS public.users_id_seq
				INCREMENT 1
				START 1
				MINVALUE 1
				MAXVALUE 2147483647
				CACHE 1;
		
				ALTER SEQUENCE users_id_seq
					OWNER TO postgres;

				CREATE TABLE IF NOT EXISTS users
				(
					uid character varying(10),
					shortid character varying(10),
					id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
					
					CONSTRAINT users_pkey PRIMARY KEY (id)
				)`)
			if err != nil {
				log.Println(err)
			} 
			log.Println("Calculating NextID & NextUID:")
			err = PGdb.QueryRow(context.Background(), `select urls.shortid from urls order by urls.id desc limit 1`).Scan(&DBLastURLID)
			log.Println("DBLastURLID =" + DBLastURLID)
			LastID, _ := strconv.Atoi(DBLastURLID)			
			initconfig.NextID =LastID + initconfig.Step
			log.Printf("NextID= %v", initconfig.NextID)
			if err != nil {
				log.Println(err)
			} 

			//err = PGdb.QueryRow(context.Background(), `select users.uid from users order by users.id desc limit 1`).Scan(&DBLastUID)
			err = PGdb.QueryRow(context.Background(), `select urls.uid from urls order by urls.uidint desc limit 1`).Scan(&DBLastUID)
			log.Println("DBLastUID =" + DBLastUID)
			LastUID, _ := strconv.Atoi(DBLastUID)			
			initconfig.NextUID = LastUID + initconfig.Step
			log.Printf("NextUID= %v", initconfig.NextUID)
			if err != nil {
				log.Println(err)
			} 
			
			PGdbOpened = true
			log.Println("PGdbOpened = TRUE") 	
		}
	} else {
		log.Println("PGdbOpened = FALSE")
	}
	
}
func DelURLIDs(ch chan string){
	log.Println("Starting async 'delete URLIDs for user' routine")
  
		//time.Sleep(time.Second * 1)
        fmt.Printf("Length of channel Input is %v and capacity of channel c is %v\n", len(ch), cap(ch))
		if len(ch)>0{
			for shortURLID := range ch {
				S.deleteShortURLfromuser(shortURLID)
			}
		}
	
}

func SetdbType(){
	log.Println("TRY SetdbType")
	log.Printf("PGdbOpened= %v", PGdbOpened)
	log.Printf("dbPathexists= %v", dbPathexists)
	log.Printf("RAMonly= %v", RAMonly)
	switch true {
		case PGdbOpened:
			log.Println("case PGdbOpened")
			DataBase :=  DataBaseconnectPGDB		
			//DataBase =  connectFileDB{}  так говорит, что объявленные выше типы не используются
			S = DataBase
		case dbPathexists:
			log.Println("case dbPathexists")
			DataBase :=  DataBaseconnectFileDB
			S = DataBase
		case RAMonly:
			log.Println("case RAMonly")
			DataBase :=  DataBaseconnectRAM	
			S = DataBase		
	
		default:
			log.Println("case default (RAMonly)")
			DataBase :=  DataBaseconnectRAM	
			S = DataBase	
	}
	
}

type storager interface{
    storeURL(fullURL string) (ShortURLID, status string)
	storeShortURLtouser(userid, shortURLid string)
	deleteShortURLfromuser(shortURLid string)
    retrieveURL(id string) (FullURL string, Status string)
	retrieveUserURLS(userid string) (output string, noURLs bool, UsrShortURLsonly []string)
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

func CheckPGdbConn() (connected bool){
	//onlyOnce.Do(Initdb)
	//defer PGdb.Close()
    err := PGdb.Ping(context.Background())
    if err != nil {
        log.Println(err)
		return false
    } else {
        log.Println("Ping db is ok")
		return true
    }	
}

func AssignShortURLtouser(userid, shortURLid string){
	//onlyOnce.Do(Initdb)
	S.storeShortURLtouser(userid, shortURLid)
}

func DeleteShortURLfromuser(ch chan string){
	//onlyOnce.Do(Initdb)

	go DelURLIDs(ch)

}

func GetuserURLS(userid string) (output string, noURLs bool, arrayUserURLs []string){
	return S.retrieveUserURLS(userid)
}

func Storerecord(fullURL string) (ShortURLID, Status string){
	//onlyOnce.Do(Initdb)
	return S.storeURL(fullURL)
}



func Getrecord(id string) (FullURL, Status string) {
	//onlyOnce.Do(Initdb)
	return S.retrieveURL(id)
}