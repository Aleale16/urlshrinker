package storage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
)

func (conn connectRAM) storeURL(fullURL string) (ShortURLID, Status string) {
	onlyOnce.Do(Initdb)
	log.Println("Running store connectRAM")
	id := strconv.Itoa(initconfig.NextID)
	URL[id] = fullURL
	initconfig.NextID = initconfig.NextID + initconfig.Step
	return id, ""
}

func (conn connectFileDB) storeURL(fullURL string) (ShortURLID, Status string) {
	onlyOnce.Do(Initdb)
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
	onlyOnce.Do(Initdb)
	id := strconv.Itoa(initconfig.NextID)
	result, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl, active) values ($1, $2) on conflict (fullurl) DO NOTHING`, id, fullURL)
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
			initconfig.NextID = initconfig.NextID + initconfig.Step
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
	_, err := PGdb.Exec(context.Background(), `insert into users(uid, shortid, active) values ($1, $2, $3)`, uid, shortURLid, true)
	if err == nil {
		log.Println("User was created, URL assigned")
	} else {
		log.Println("User was not created (inserted), URL NOT assigned!")
		log.Println(err)
	}
} 
func (conn connectFileDB) storeShortURLtouser(userid, shortURLid string){
} 

func (conn connectRAM) deleteShortURLfromuser(shortURLid string){
} 
func (conn connectPGDB) deleteShortURLfromuser(shortURLid string){
	//uid := userid
	_, err := PGdb.Exec(context.Background(), `update users set active = false where shortid=$1`, shortURLid)
	if err == nil {
		log.Printf("URL %v was disabled", shortURLid)
	} else {
		log.Println(err)
	}
} 
func (conn connectFileDB) deleteShortURLfromuser(shortURLid string){
} 

func (conn connectRAM) retrieveURL(id string) (FullURL string) {
	FullURL = URL[id]
	if (FullURL != ""){
		return FullURL
	} else {
		return "http://google.com/404"		
	}
}
func (conn connectFileDB) retrieveURL(id string) (FullURL string) {
	FullURL = URL[id]
	return FullURL
}
func (conn connectPGDB) retrieveURL(id string) (FullURL string) {
	err := PGdb.QueryRow(context.Background(), "SELECT urls.fullurl FROM urls where shortid=$1", id).Scan(&FullURL)
	if err != nil {
		return err.Error()
	}
	return FullURL
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
					ShortURL:	initconfig.BaseURL + "/?id=" + v,
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
		log.Println("JSONresult= ")		
		log.Println(JSONresult)		
		log.Println(string(JSONresult))		
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
	rows, err := PGdb.Query(context.Background(), "SELECT usr.uid, usr.shortid, urls.fullurl FROM users as usr LEFT JOIN urls ON urls.shortid = usr.shortid where uid=$1", userid)
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
			ShortURL:	initconfig.BaseURL + "/?id=" + shortID,
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
	log.Println("JSONresult= ")		
	log.Println(JSONresult)		
	log.Println(string(JSONresult))		
	noURLs = false
	return string(JSONresult), noURLs, UsrShortURLs
}