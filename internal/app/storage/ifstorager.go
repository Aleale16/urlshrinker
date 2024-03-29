package storage

import (
	"context"
	"encoding/json"

	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"

	//"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"urlshrinker/internal/app/initconfig"
)

var mu sync.Mutex

// storeURL - store URL if connectRAM.
func (conn connectRAM) storeURL(fullURL string) (ShortURLID, Status string) {
	//onlyOnce.Do(Initdb)
	log.Info().Msg("Running store connectRAM")
	id := strconv.Itoa(initconfig.NextID)
	URL[id] = fullURL
	initconfig.NextID = initconfig.NextID + initconfig.Step
	return id, ""
}

// storeURL - store URL if connectFileDB.
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

// storeURL - store URL if connectPGDB.
func (conn connectPGDB) storeURL(fullURL string) (ShortURLID, Status string) {
	//onlyOnce.Do(Initdb)
	id := strconv.Itoa(initconfig.NextID)
	result, err := PGdb.Exec(context.Background(), `insert into urls(shortid, fullurl, active) values ($1, $2, $3) on conflict (fullurl) DO NOTHING`, id, fullURL, true)
	if err == nil {
		if result.RowsAffected() == 0 {
			// ShortID load from request
			var ShortID string
			err := PGdb.QueryRow(context.Background(), "SELECT urls.shortid FROM urls where fullurl=$1", fullURL).Scan(&ShortID)
			if err != nil {
				log.Info().Msg(err.Error())
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
		log.Info().Msg(err.Error())
	}
	return id, Status
}

// storeShortURLtouser - store ShortURL to user if connectRAM.
func (conn connectRAM) storeShortURLtouser(userid, shortURLid string) {
	uid := userid
	Usr[uid] = append(Usr[uid], shortURLid)
	log.Info().Msgf("AssignShortURLtouser: "+string(uid)+" shortURLid= %v", Usr[uid])
}

// storeShortURLtouser - store ShortURL to user if connectPGDB.
func (conn connectPGDB) storeShortURLtouser(userid, shortURLid string) {
	uid := userid
	uidint, _ := strconv.Atoi(userid)
	//_, err := PGdb.Exec(context.Background(), `insert into users(uid, shortid, active) values ($1, $2, $3)`, uid, shortURLid, true)
	//_, err := PGdb.Exec(context.Background(), `insert into urls(uid, shortid, active) values ($1, $2, $3)`, uid, shortURLid, true)
	_, err := PGdb.Exec(context.Background(), `update urls set uid = $1, uidint = $2 where shortid=$3`, uid, uidint, shortURLid)
	if err == nil {
		log.Printf("User %v was created, URL %v assigned", uid, shortURLid)
	} else {
		log.Printf("User %v was not created (inserted), URL %v NOT assigned!", uid, shortURLid)
		log.Info().Msg(err.Error())
	}
}

// storeShortURLtouser - store ShortURL to user if connectFileDB.
// Deprecated: not required.
func (conn connectFileDB) storeShortURLtouser(userid, shortURLid string) {
}

// deleteShortURLfromuser - delete ShortURL to user if connectRAM.
func (conn connectRAM) deleteShortURLfromuser(shortURLid string) {
	URL[shortURLid] = "*" + URL[shortURLid]
	log.Printf("URL %v was disabled with *", shortURLid)
}

// deleteShortURLfromuser - delete ShortURL to user if connectPGDB.
func (conn connectPGDB) deleteShortURLfromuser(shortURLid string) {
	//uid := userid
	_, err := PGdb.Exec(context.Background(), `update urls set active = false where shortid=$1`, shortURLid)
	if err == nil {
		log.Printf("URL %v was disabled", shortURLid)
	} else {
		log.Info().Msg(err.Error())
	}
}

// deleteShortURLfromuser - delete ShortURL to user if connectFileDB.
func (conn connectFileDB) deleteShortURLfromuser(shortURLid string) {
}

// containsinStr - check if Substring exist in String.
func containsinStr(s string, e string) bool {
	for _, a := range s {
		if string(a) == e {
			return true
		}
	}
	return false
}

// retrieveURL - retrieve URL if connectRAM
func (conn connectRAM) retrieveURL(id string) (FullURL string, Status string) {
	log.Debug().Msgf("RAM retrieveURL ID=%v", id)
	FullURL = URL[id]
	withAsterisk := containsinStr(FullURL, "*")
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

// retrieveURL - retrieve URL if connectFileDB
func (conn connectFileDB) retrieveURL(id string) (FullURL string, Status string) {
	log.Debug().Msgf("FileDB retrieveURL ID=%v", id)
	FullURL = URL[id]
	return FullURL, "307"
}

// retrieveURL - retrieve URL if connectPGDB
func (conn connectPGDB) retrieveURL(id string) (FullURL string, Status string) {
	log.Debug().Msgf("PGDB retrieveURL ID=%v", id)
	// activelink load from request
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

// retrieveUserURLS - retrieve User URLs if connectRAM
func (conn connectRAM) retrieveUserURLS(userid string) (output string, noURLs bool, UsrShortURLsonly []string) {
	// JSON vars
	var (
		UsrURLJSON []UsrURLJSONrecord
		JSONresult []byte
	)
	noURLs = true
	UsrShortURLs := Usr[userid]
	if len(UsrShortURLs) > 0 {
		for _, v := range UsrShortURLs {
			log.Info().Msg(v)
			//Так нормально заполнять JSON перед маршаллингом?
			UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
				//ShortURL:	initconfig.BaseURL + "/?id=" + v,
				ShortURL: initconfig.BaseURL + "/" + v,
				FullURL:  URL[v],
			})
		}
		JSONdata, err := json.Marshal(&UsrURLJSON)
		if err != nil {
			return err.Error(), noURLs, UsrShortURLs
		}
		//JSONdata = append(JSONdata, '\n')
		//URL[id] = string(JSONdata)
		JSONresult = JSONdata
		//log.Info().Msg("JSONresult= ")
		//log.Info().Msg(JSONresult)
		shortURLpathJSONBz, err := json.MarshalIndent(&UsrURLJSON, "", "  ")
		if err != nil {
			panic(err.Error())
		}
		log.Printf("retrieveUserURLS JSONresult= %v", string(shortURLpathJSONBz))
		noURLs = false
	}
	return string(JSONresult), noURLs, UsrShortURLs
}

// retrieveUserURLS - retrieve User URLs if connectFileDB
func (conn connectFileDB) retrieveUserURLS(userid string) (output string, noURLs bool, UsrShortURLsonly []string) {
	return "", true, []string{}
}

// retrieveUserURLS - retrieve User URLs if connectPGDB
func (conn connectPGDB) retrieveUserURLS(userid string) (output string, noURLs bool, UsrShortURLsonly []string) {
	// JSON vars
	var (
		UsrURLJSON   []UsrURLJSONrecord
		JSONresult   []byte
		UID          string
		shortID      string
		FullURL      string
		UsrShortURLs []string
	)
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
			log.Fatal().Msg(err.Error())
		}
		UsrShortURLs = append(UsrShortURLs, shortID)
		UsrURLJSON = append(UsrURLJSON, UsrURLJSONrecord{
			//ShortURL:	initconfig.BaseURL + "/?id=" + shortID,
			ShortURL: initconfig.BaseURL + "/" + shortID,
			FullURL:  FullURL,
		})
	}
	JSONdata, err := json.Marshal(&UsrURLJSON)
	if err != nil {
		return err.Error(), noURLs, UsrShortURLs
	}
	//JSONdata = append(JSONdata, '\n')
	//URL[id] = string(JSONdata)
	JSONresult = JSONdata
	//log.Info().Msg("JSONresult= ")
	//log.Info().Msg(JSONresult)
	shortURLpathJSONBz, err := json.MarshalIndent(&UsrURLJSON, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("JSONresult= %v", string(shortURLpathJSONBz))
	noURLs = false
	return string(JSONresult), noURLs, UsrShortURLs
}
