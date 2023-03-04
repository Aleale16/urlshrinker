package storage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
)

var mu sync.Mutex

func (conn connectRAM) storeURL(fullURL string) (ShortURLID, Status string) {
	//onlyOnce.Do(Initdb)
	log.Println("Running store connectRAM")
	id := strconv.Itoa(initconfig.NextID)
	URL[id] = fullURL
	initconfig.NextID = initconfig.NextID + initconfig.Step
	return id, ""
}

func (conn connectFileDB) storeURL(fullURL string) (ShortURLID, Status string) {
	//onlyOnce.Do(Initdb)
	id := strconv.Itoa(initconfig.NextID)
	URLJSONline := URLJSONrecord{
		ID:      id,
		FullURL: fullURL,
	}
	JSONdata, err := json.Marshal(&URLJSONline)
	if err != nil {
		return err.Error(), ""
	}
	JSONdata = append(JSONdata, '\n')
	//URL[id] = string(JSONdata)
	URL[id] = fullURL

	DBfile, _ := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	_, err = DBfile.Write(JSONdata)
	if err != nil {
		return err.Error(), ""
	}
	DBfile.Close()
	initconfig.NextID = initconfig.NextID + initconfig.Step
	return id, ""
}

func (conn connectPGDB) storeURL(fullURL string) (ShortURLID, Status string) {
	//onlyOnce.Do(Initdb)
	id := strconv.Itoa(initconfig.NextID)
	result, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl, active) values ($1, $2, $3) on conflict (fullurl) DO NOTHING`, id, fullURL, true)
	if err == nil {
		if result.RowsAffected() == 0 {
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
			mu.Lock()
			initconfig.NextID = initconfig.NextID + initconfig.Step
			mu.Unlock()
		}
	} else {
		log.Println(err)
	}
	return id, Status
}

func (conn connectRAM) storeShortURLtouser(userid, shortURLid string){
	uid := userid
	Usr[uid] = append(Usr[uid], shortURLid)
	log.Println("AssignShortURLtouser: " + string(uid)+ " shortURLid= " )	
	log.Println(Usr[uid])
} 
func (conn connectPGDB) storeShortURLtouser(userid, shortURLid string){
	uid := userid
	uidint, _ := strconv.Atoi(userid)
	//_, err := PGdb.Exec(context.Background(), `insert into users(uid, shortid, active) values ($1, $2, $3)`, uid, shortURLid, true)
	//_, err := PGdb.Exec(context.Background(), `insert into urls(uid, shortid, active) values ($1, $2, $3)`, uid, shortURLid, true)
	_, err := PGdb.Exec(context.Background(), `update urls set uid = $1, uidint = $2 where shortid=$3`, uid, uidint, shortURLid)
	if err == nil {
		log.Printf("User %v was created, URL %v assigned", uid, shortURLid)
	} else {
		log.Printf("User %v was not created (inserted), URL %v NOT assigned!", uid, shortURLid)
		log.Println(err)
	}
} 
func (conn connectFileDB) storeShortURLtouser(userid, shortURLid string){
} 

func (conn connectRAM) deleteShortURLfromuser(shortURLid string){
	URL[shortURLid] = "*" + URL[shortURLid]
	log.Printf("URL %v was disabled with *", shortURLid)
} 
func (conn connectPGDB) deleteShortURLfromuser(shortURLid string){
	//uid := userid
	_, err := PGdb.Exec(context.Background(), `update urls set active = false where shortid=$1`, shortURLid)
	if err == nil {
		log.Printf("URL %v was disabled", shortURLid)
	} else {
		log.Println(err)
	}
} 
func (conn connectFileDB) deleteShortURLfromuser(shortURLid string){
} 

func containsinStr(s string, e string) bool {
    for _, a := range s {
        if string(a) == e {
            return true
        }
    }
    return false
}

func (conn connectRAM) retrieveURL(id string) (FullURL string, Status string) {
	FullURL = URL[id]
	withAsterisk := containsinStr(FullURL , "*")
	log.Printf("withAsterisk = %v", withAsterisk)
	switch true {
		case (FullURL == ""):
			return "http://google.com/404", "400"
		case (withAsterisk):
			return "", "410"
		default:
			return FullURL, "307"
	}
	/*
	if (FullURL != ""){
		return FullURL, "307"
	} else {
		return "http://google.com/404", "400"		
		//return "", "400"		
	}
	*/
}
func (conn connectFileDB) retrieveURL(id string) (FullURL string, Status string) {
	FullURL = URL[id]
	return FullURL, "307"
}
func (conn connectPGDB) retrieveURL(id string) (FullURL string, Status string) {
	var activelink bool
	err := PGdb.QueryRow(context.Background(), "SELECT urls.fullurl, urls.active FROM urls where shortid=$1", id).Scan(&FullURL, &activelink)
	if err != nil {
		return err.Error(), ""
	}
	if activelink {
		return FullURL, "307"
		} else {
			return FullURL, "410"
		}
	
}

func (conn connectRAM) retrieveUserURLS (userid string) (output string, noURLs bool, UsrShortURLsonly []string){
	var UsrURLJSON []UsrURLJSONrecord
	var JSONresult []byte
	noURLs = true 
	UsrShortURLs := Usr[userid]
	if len(UsrShortURLs)>0{
		for _, v := range UsrShortURLs {	
				log.Println(v)
				//Так нормально заполнять JSON перед маршаллингом?
				UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
					//ShortURL:	initconfig.BaseURL + "/?id=" + v,
					ShortURL:	initconfig.BaseURL + "/" + v,
					FullURL:	URL[v],
				})				
		}
		JSONdata, err := json.Marshal(&UsrURLJSON)
		if err != nil {
			return err.Error(), noURLs, UsrShortURLs
		}
		//JSONdata = append(JSONdata, '\n')
		//URL[id] = string(JSONdata)
		JSONresult = JSONdata
		//log.Println("JSONresult= ")		
		//log.Println(JSONresult)
		shortURLpathJSONBz, err := json.MarshalIndent(&UsrURLJSON, "", "  ")
		if err != nil {
			panic(err)
		}		
		log.Printf("retrieveUserURLS JSONresult= %v", string(shortURLpathJSONBz))		
		noURLs = false
	}
	return string(JSONresult), noURLs, UsrShortURLs
}
func (conn connectFileDB) retrieveUserURLS (userid string) (output string, noURLs bool, UsrShortURLsonly []string){
	return "", true, []string{}
}
func (conn connectPGDB) retrieveUserURLS (userid string) (output string, noURLs bool, UsrShortURLsonly []string){
	var (
		UsrURLJSON []UsrURLJSONrecord
		JSONresult []byte
		UID string
		shortID string
		FullURL string
	)
	var UsrShortURLs []string
	noURLs = true
	//rows, err := PGdb.Query(context.Background(), "SELECT usr.uid, usr.shortid, urls.fullurl FROM users as usr LEFT JOIN urls ON urls.shortid = usr.shortid where uid=$1", userid)
	rows, err := PGdb.Query(context.Background(), "SELECT urls.uid, urls.shortid, urls.fullurl FROM urls where uid=$1", userid)
	if err != nil {
		return err.Error(), noURLs, []string{}
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		err := rows.Scan(&UID, &shortID, &FullURL)
		if err != nil {
			log.Fatal(err)
		}
		UsrShortURLs = append(UsrShortURLs, shortID)
		UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
			//ShortURL:	initconfig.BaseURL + "/?id=" + shortID,
			ShortURL:	initconfig.BaseURL + "/" + shortID,
			FullURL:	FullURL,
		})	
	}
	JSONdata, err := json.Marshal(&UsrURLJSON)
	if err != nil {
		return err.Error(), noURLs, UsrShortURLs
	}
	//JSONdata = append(JSONdata, '\n')
	//URL[id] = string(JSONdata)
	JSONresult = JSONdata
	//log.Println("JSONresult= ")		
	//log.Println(JSONresult)		
	shortURLpathJSONBz, err := json.MarshalIndent(&UsrURLJSON, "", "  ")
	if err != nil {
		panic(err)
	}		
	log.Printf("JSONresult= %v", string(shortURLpathJSONBz))		
	noURLs = false
	return string(JSONresult), noURLs, UsrShortURLs
}