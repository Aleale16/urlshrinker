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

	"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	"github.com/Aleale16/urlshrinker/internal/app/storage"
)

func StatusOKHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)

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

	cookie := &http.Cookie{
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
	w.Header().Set("Authorization", cookie.Value)
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
		fmt.Println("Checking useridcookie:")
		useridcookie, err:= r.Cookie("userid")
		if err != nil{	
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		} else {	
			validSign, id := checkSign(useridcookie.Value)
			fmt.Println(id)
			fmt.Println(validSign)
			//if validSign {
				userURLS, noURLs := storage.GetuserURLS(id)
			if noURLs{
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(userURLS))
			}
			//}
		}
	} else {
		fmt.Println("Checking authorizationHeader:")
		validSign, id := checkSign(authorizationHeader)
		fmt.Println(id)
		fmt.Println(validSign)
		//if validSign {
			userURLS, noURLs := storage.GetuserURLS(id)
		if noURLs{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(userURLS))
		}
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")
	//q := r.URL.String()
    if q == "" {
        http.Error(w, "The query parameter is missing", http.StatusBadRequest)
        return
    }
	record := storage.Getrecord(q)	
	//if record != "http://google.com/404" {
		// устанавливаем заголовок Location	
		w.Header().Set("Location", record)
		// устанавливаем статус-код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
	//} else {
	//	http.Error(w, "Short URL with id=" + q + " not set", http.StatusBadRequest)
	//}

	fmt.Println("GET: " + q + " Redirect to " + record)
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

func PostHandler(w http.ResponseWriter, r *http.Request) /*(shortURL string)*/{
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
	fmt.Println(r.Cookie("userid"))
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
	
	//w.Write([]byte(useridcookie.Value))
	
	shortURLid := storage.Storerecord(string(body))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	//shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid	
	//shortURLpath := BaseURL + "/?id="+ shortURLid Как сюда передать переменную из server.go?	
	//вот так из пакета initconfig:
	shortURLpath :=initconfig.BaseURL + "/?id="+ shortURLid
	storage.AssignShortURLtouser(uid, shortURLid)
	
	//w.Header().Set("Content-Encoding", "gzip, deflate, br")
	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)
	//отладка что было в POST запросе
	//w.Write([]byte(b))

//типа return:
	w.Write([]byte(shortURLpath))

	fmt.Println("POST: " + string(body)+ " return id= "+ shortURLid)		

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

	shortURLid := storage.Storerecord(string(postJSON.URL))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	//shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid
	shortURLpath := initconfig.BaseURL + "/?id="+ shortURLid
	
	var shortURLpathJSON resultData
	shortURLpathJSON.ShortURL = shortURLpath


	w.Header().Set("Content-Type", "application/json")
	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	

	shortURLpathJSONBz, err := json.MarshalIndent(shortURLpathJSON, "", "  ")
	if err != nil {
        panic(err)
    }
//типа return:
	w.Write(shortURLpathJSONBz)
	//или?
	//w.Write([]byte(shortURLpathJSON))
	//w.Write([]byte(shortURLpath))

	fmt.Println("POST: " + string(b)+ " return id= "+ shortURLid + " return JSON= "+ string(shortURLpathJSONBz))	

	//return shortURLpath
}

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
				shortURLid := storage.Storerecord(string(v.URL))
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
