package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"github.com/Aleale16/urlshrinker/internal/app/storage"
)

func StatusOKHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
	n := 3
	wg.Add(n)
	go func() {
		for i := 0; i < n; i++{
			time.Sleep(time.Second * 2)
			log.Println("Server is still alive!")
			wg.Done()
		}
	}()	
	wg.Wait()
}
/* Наверное, больше не пригодится, если и дальше использовать Chi
func ReqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
        PostHandler(w, r)
    }
	if r.Method == http.MethodGet {
        GetHandler(w, r)
    }
	fmt.Println(r.Method)
}
*/

func defineCookie(w http.ResponseWriter, r *http.Request)(uid string){

	var key = []byte("secret key")
	//userid := []byte(strconv.Itoa(rand.Intn(9999)))
	//userid := []byte("8888")
	userid := []byte(strconv.Itoa(initconfig.NextUID))
	fmt.Println("New userid=" + strconv.Itoa(initconfig.NextUID))
	initconfig.NextUID = initconfig.NextUID + initconfig.Step
	
      // подписываем алгоритмом HMAC, используя SHA256
	  h := hmac.New(sha256.New, key)
	  h.Write(userid)
	  dst := h.Sum(nil)

	 //вот это вообще было не очевидно:! 
	  signedcookie := string(dst) + string(userid)
  
	  fmt.Printf("%x", dst)
	  fmt.Printf("%v\n", dst)

/*	cookie := &http.Cookie{
        Name:   "userid",
        Value:  hex.EncodeToString([]byte(signedcookie)),
        MaxAge: 300,
		Path:  "/",
		HttpOnly: true,
        Secure:   true,
    }
	http.SetCookie(w, cookie)

	fmt.Println("cookie was set: " + cookie.Name + "; value= " + cookie.Value)

	//if cookie.Value != ""{
		//checkSign(cookie.Value)
	//}
	fmt.Println(r.Cookie("userid"))
	//w.Header().Set("Authorization", cookie.Value)
	*/
	w.Header().Set("Authorization", hex.EncodeToString([]byte(signedcookie)))
	return string(userid)
}

func checkSign(msg string) (validSign bool, val string){
	var key = []byte("secret key")
	var (
		data []byte // декодированное сообщение с подписью
		id   string // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)
	validSign = false
	data, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("data=" + string(data))
	id = string(data[sha256.Size:])
	val = id
	//id = binary.BigEndian.Uint32(data[:4])
	//id = binary.BigEndian.Uint32(data[sha256.Size:])
	h := hmac.New(sha256.New, key)
	h.Write(data[sha256.Size:])
	sign = h.Sum(nil) 
	if hmac.Equal(sign, data[:sha256.Size]) {
		fmt.Println("Подпись подлинная. ID:", id)
		validSign = true
	} else {
		fmt.Println("Подпись неверна. Где-то ошибка! ID:", id)
	}	
	return validSign, val
}

func GetUsrURLsHandler(w http.ResponseWriter, r *http.Request) {
	//То, что автотест ожидает, а затем отправляет токен в поле заголовка Authorization можно было узнать только в результате просмотра текста автотеста!
	authorizationHeader := r.Header.Get("Authorization")
	fmt.Println("authorizationHeader=" + authorizationHeader)
	if authorizationHeader == ""{
		fmt.Println("Empty authorizationHeader:")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		/*fmt.Println("Checking useridcookie:")
		useridcookie, err:= r.Cookie("userid")
		if err != nil{	
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		} else {	
			//validSign, id := checkSign(useridcookie.Value)
			validSign, id := checkSign(authorizationHeader)
			fmt.Println(id)
			fmt.Println(validSign)
			//if validSign {
				userURLS, noURLs, _ := storage.GetuserURLS(id)
			if noURLs{
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(userURLS))
			}*/
			//}
		//}
	} else {
		fmt.Println("Checking authorizationHeader:")
		validSign, id := checkSign(authorizationHeader)
		fmt.Println(id)
		fmt.Println(validSign)
		//if validSign {
			userURLS, noURLs, _ := storage.GetuserURLS(id)
		if noURLs{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(userURLS))
		}
	}
	fmt.Println("GET: /api/user/urls ")
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")
	//q := r.URL.String()
    if q == "" {
        http.Error(w, "The query parameter is missing", http.StatusBadRequest)
        return
    }
	record, Status := storage.Getrecord(q)	
	//if record != "http://google.com/404" {
		// устанавливаем заголовок Location	
		w.Header().Set("Location", record)
		switch Status{
			case "307": // устанавливаем статус-код 307
				w.WriteHeader(http.StatusTemporaryRedirect)
			case "410": // устанавливаем статус-код 410
				w.WriteHeader(http.StatusGone)
		}
		/*// устанавливаем статус-код 307
		w.WriteHeader(http.StatusTemporaryRedirect)*/
	//} else {
	//	http.Error(w, "Short URL with id=" + q + " not set", http.StatusBadRequest)
	//}

	fmt.Println("GET: / " + q + " Redirect to " + record)
}

func GetPingHandler(w http.ResponseWriter, r *http.Request) {
// работаем с базой storage.PGdb
	if storage.CheckPGdbConn(){
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println("GetPingHandler: finished")
}

//! POST /
func PostHandler(w http.ResponseWriter, r *http.Request) /*(shortURL string)*/{
	authorizationHeader := r.Header.Get("Authorization")
	fmt.Println("authorizationHeader=" + authorizationHeader)
	
/*	// читаем Body (Тело POST запроса)
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
*/	
	// обработаем ситуацию, если на вход может прийти сжатое содержимое	
	// переменная reader будет равна r.Body или *gzip.Reader
	var reader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		w.Header().Set("Accept-Encoding", "gzip")
		w.Header().Set("Content-Encoding", "gzip, deflate, br")
    // создаём *gzip.Reader, который будет читать тело запроса
    // и распаковывать его
    gz, err := gzip.NewReader(r.Body)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
	reader = gz
    // потом закрыть *gzip.Reader
    defer gz.Close()
    } else {
        reader = r.Body
    }
    // при чтении вернётся распакованный слайс байт
    body, err := io.ReadAll(reader)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    log.Println(w, "body: %d", body)

	uid := ""
/*	fmt.Println(r.Cookie("userid"))
	useridcookie, err:= r.Cookie("userid")
	if err != nil{	
		fmt.Println(err)
		uid = defineCookie(w, r)
	} else {
		//if useridcookie.Value == "" {
			validSign, id := checkSign(useridcookie.Value)
			fmt.Println(id)
			if !validSign {	
				uid = defineCookie(w, r)
			} else {
				uid = id
			}
		//} else {
		//	defineCookie(w, r)
		//}
	}
*/
	if authorizationHeader == ""{
		uid = defineCookie(w, r)
		} else {
			validSign, id := checkSign(authorizationHeader)
			fmt.Println(id)
			if !validSign {	
				uid = defineCookie(w, r)
			} else {
				uid = id
			}
		}
	
	//w.Write([]byte(useridcookie.Value))
	
	shortURLid, Status := storage.Storerecord(string(body))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	//shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid	
	//shortURLpath := BaseURL + "/?id="+ shortURLid Как сюда передать переменную из server.go?	
	//вот так из пакета initconfig:
	shortURLpath :=initconfig.BaseURL + "/?id="+ shortURLid
	storage.AssignShortURLtouser(uid, shortURLid)
	
	//w.Header().Set("Content-Encoding", "gzip, deflate, br")
	if Status == "StatusConflict"{
		// устанавливаем статус-код 409
		w.WriteHeader(http.StatusConflict)
	} else {
		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
	}
	//отладка что было в POST запросе
	//w.Write([]byte(b))

//типа return:
	w.Write([]byte(shortURLpath))

	fmt.Println("POST: / " + string(body)+ " return id= "+ shortURLid)		

	//return shortURLpath
}

//структура вводимого JSON
type inputData struct {
    //ID int `json:"ID"`
    URL string `json:"url,omitempty"`   
}

//структура выводимого JSON	 
type resultData struct {
    //ID int `json:"ID"`
    ShortURL string `json:"result"`    
}

//! POST /api/shorten
func PostJSONHandler(w http.ResponseWriter, r *http.Request) /*(shortURL string)*/{
	// читаем Body (Тело POST запроса)
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//отладка всё что было в POST запросе
	log.Println("PostJSONHandler body: " + string(b))
	//Добавить?
	//type Example struct {
	//	URL   string `valid:"url"`
	//}
	log.Println("Content-Encoding from req: " + r.Header.Get("Content-Encoding"))

	var postJSON inputData
	err = json.Unmarshal(b, &postJSON)
	if err != nil {
		panic(err)
	}
	//отладка что было в поле url в POST запросе
	log.Println(postJSON.URL)

	shortURLid, Status := storage.Storerecord(string(postJSON.URL))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	//shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid
	shortURLpath := initconfig.BaseURL + "/?id="+ shortURLid
	
	var shortURLpathJSON resultData
	shortURLpathJSON.ShortURL = shortURLpath


	w.Header().Set("Content-Type", "application/json")
	if Status == "StatusConflict"{
		// устанавливаем статус-код 409
		w.WriteHeader(http.StatusConflict)
	} else {
		// устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
	}
	

	shortURLpathJSONBz, err := json.MarshalIndent(shortURLpathJSON, "", "  ")
	if err != nil {
        panic(err)
    }
//типа return:
	w.Write(shortURLpathJSONBz)
	//или?
	//w.Write([]byte(shortURLpathJSON))
	//w.Write([]byte(shortURLpath))

	fmt.Println("POST: /api/shorten " + string(b)+ " return id= "+ shortURLid + " return JSON= "+ string(shortURLpathJSONBz))	

	//return shortURLpath
}

//! DELETE /api/user/urls

func getInputChan(listURLids []string) (ch chan string) {
    // make return channel
    //input := make(chan string, 100)
    ch = make(chan string, 100)
	//var numbers []string

    // sample numbers
    //numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	/*for _, v := range listURLids {	
		log.Println(v)	
		numbers = append(numbers, string(v))			
	}*/

    // run goroutine
    go func() {
        for _, URLid := range listURLids {
            //initconfig.InputIDstoDel <- URLid
            ch <- URLid
        }
        // close channel once all numbers are sent to channel
       // close(input)
    }()

    //return initconfig.InputIDstoDel
    return ch
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	var listURLids []string
	var InvalidURLIDexists bool
	//var IDstoDel = make(chan string, 7)
	//storage.DeleteShortURLfromuser()
	
	// читаем Body (Тело POST запроса)
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println("DeleteURLsHandler body: " + string(b))
	log.Println("Content-Encoding from req: " + r.Header.Get("Content-Encoding"))
	err = json.Unmarshal(b, &listURLids)
	if err != nil {
		panic(err)
	}
	log.Println(listURLids)	

	if len(listURLids)>0{
		authorizationHeader := r.Header.Get("Authorization")
		log.Println("authorizationHeader=" + authorizationHeader)
		if authorizationHeader != ""{
		//if authorizationHeader == ""{
			log.Println("Checking authorizationHeader:")
			validSign, id := checkSign(authorizationHeader)
			log.Printf("User with %v", id)
			log.Printf("Authenticated???: %v", validSign)
	//		validSign = true
			if validSign{
				userURLS, noURLs, arrayUserURLs := storage.GetuserURLS(id)

				if !noURLs && len(userURLS)>= len(listURLids){
					InvalidURLIDexists = false
					for _, v := range(listURLids){
						if !InvalidURLIDexists{
							if !contains(arrayUserURLs, v){
								InvalidURLIDexists = true
							}
						}
					}
					if !InvalidURLIDexists {
						IDstoDel := getInputChan(listURLids)
						// устанавливаем статус-код 202
						w.WriteHeader(http.StatusAccepted)
						storage.DeleteShortURLfromuser(IDstoDel)
					}
				} else {
					InvalidURLIDexists = true
					log.Println("No (invalid) ShortURLs to delete for user")
				}

			}
		} else {
			log.Println("Empty authorizationHeader for user")
		} 
		
		
		//uid := "9999"
		

		
	} else {
		InvalidURLIDexists = true
		log.Println("No (invalid) ShortURLs to delete for user")
	}
	fmt.Println("DELETE: " + string(b))
}
/*
func bgfunc(chanInputIDs){
	for shortURLID := range chanInputIDs {
		storage.DeleteShortURLfromuser(shortURLID)
	}

}*/

//! POST /api/shorten/batch
//структура вводимого JSON
type inputbatchData struct {
    ID string `json:"correlation_id"`
    URL string `json:"original_url"`   
}

//структура выводимого JSON	 
type resultbatchData struct {
    ID string `json:"correlation_id"`
    ShortURL string `json:"short_url"`    
}

func PostJSONbatchHandler(w http.ResponseWriter, r *http.Request) /*(shortURL string)*/{
	// читаем Body (Тело POST запроса)
	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//отладка всё что было в POST запросе
	log.Println("PostJSONHandler body: " + string(b))
	log.Println("Content-Encoding from req: " + r.Header.Get("Content-Encoding"))

	var inputbatchJSON []inputbatchData
	var resultbatchJSON []resultbatchData
	var JSONresult []byte
	err = json.Unmarshal(b, &inputbatchJSON)
	if err != nil {
		panic(err)
	}
	//отладка что было в поле url в POST запросе
	log.Println(inputbatchJSON)
	// Обработка входного JSON и выдача результирующего
	if len(inputbatchJSON)>0{
		for _, v := range inputbatchJSON {	
				log.Println(v)
				shortURLid, _ := storage.Storerecord(string(v.URL))
				resultbatchJSON = append(resultbatchJSON, resultbatchData{
					ID:	v.ID,
					ShortURL:	initconfig.BaseURL + "/?id=" + shortURLid,
				})	
		}
	}
	JSONdata, err := json.Marshal(&resultbatchJSON)
	if err != nil {
		log.Println(err.Error())
	}

	JSONresult = JSONdata
	w.Header().Set("Content-Type", "application/json")
	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
//типа return:
	w.Write(JSONresult)
	fmt.Println("POST: " + string(b) + " return JSON= "+ string(JSONresult))	
}
