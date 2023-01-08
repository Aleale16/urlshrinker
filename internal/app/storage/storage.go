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
					CONSTRAINT urls_pkey PRIMARY KEY (id)
				)
				ALTER TABLE urls ADD UNIQUE (fullurl)`)
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
			log.Println(initconfig.NextID)
			if err != nil {
				log.Println(err)
			} 

			err = PGdb.QueryRow(context.Background(), `select users.uid from users order by users.id desc limit 1`).Scan(&DBLastUID)
			log.Println("DBLastUID =" + DBLastUID)
			LastUID, _ := strconv.Atoi(DBLastUID)			
			initconfig.NextUID = LastUID + initconfig.Step
			log.Print(initconfig.NextUID)
			if err != nil {
				log.Println(err)
			} 
			
			PGdbOpened = true
			fmt.Println("PGdbOpened = TRUE") 
		}
	} else {
		fmt.Println("PGdbOpened = FALSE")
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
			log.Println(err)
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

func Storerecord(fullURL string) (ShortURLID, Status string){
	onlyOnce.Do(Initdb)
	//id := strconv.Itoa(rand.Intn(9999))
	id := strconv.Itoa(initconfig.NextID)
	
	/*for (!isnewID(id)){
		id = strconv.Itoa(rand.Intn(9999))
	}*/

	if RAMonly {
		URL[id] = fullURL
		initconfig.NextID = initconfig.NextID + initconfig.Step
	} else {
		URLJSONline := URLJSONrecord{
			ID:			id,
        	FullURL:	fullURL,
		}
		JSONdata, err := json.Marshal(&URLJSONline)
		if err != nil {
			return err.Error(), ""
		}
		JSONdata = append(JSONdata, '\n')
		//URL[id] = string(JSONdata)
		URL[id] = fullURL

		DBfile, _ := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_APPEND , 0777)
		_, err = DBfile.Write(JSONdata)	
		if err != nil {	
			return err.Error(), ""
		}
		DBfile.Close()
		initconfig.NextID = initconfig.NextID + initconfig.Step
	}

	if PGdbOpened {
		result, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl) values ($1, $2) on conflict (fullurl) DO NOTHING`, id, fullURL)
		if err == nil {
			if result.RowsAffected() == 0{
				var ShortID string
				err := PGdb.QueryRow(context.Background(), "SELECT urls.shortid FROM urls where fullurl=$1", fullURL).Scan(&ShortID)
				if err != nil {
					log.Println(err)
				}
				log.Printf("Value %q, already exist in DB, rows affected =%v, ShortURL id = %q", fullURL, result.RowsAffected(), ShortID)	
				id = ShortID	
				Status = "StatusConflict"	
			} else {
				log.Printf("Values %q, %q inserted successfully, rows affected =%v", id, fullURL, result.RowsAffected())
				initconfig.NextID = initconfig.NextID + initconfig.Step
			}
		} else {
			log.Println(err)
		}
	}	
	return id, Status
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
		err := PGdb.QueryRow(context.Background(), "SELECT urls.fullurl FROM urls where shortid=$1", id).Scan(&FullURL)
		if err != nil {
			return err.Error()
		}
		result = FullURL
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