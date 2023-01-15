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
	result, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl) values ($1, $2) on conflict (fullurl) DO NOTHING`, id, fullURL)
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
	_, err := PGdb.Exec(context.Background(), `insert into users(uid, shortid) values ($1, $2)`, uid, shortURLid)
	if err == nil {
		log.Println("User was created, URL assigned")
	} else {
		log.Println("User was not created (inserted), URL NOT assigned!")
		log.Println(err)
	}
} 
func (conn connectFileDB) storeShortURLtouser(userid, shortURLid string){
} 

func (conn connectRAM) retrieveURL(id string) (FullURL string) {
	FullURL = URL[id]
	return FullURL
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