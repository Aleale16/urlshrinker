package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
func GetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("id")
	//q := r.URL.String()
    if q == "" {
        http.Error(w, "The query parameter is missing", http.StatusBadRequest)
        return
    }	
	w.Header().Set("Content-Encoding", "gzip, deflate, br")
	w.Header().Set("Accept-Encoding", "gzip, deflate, br")
	// устанавливаем заголовок Location	
	w.Header().Set("Location", storage.Getrecord(q))
	// устанавливаем статус-код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	

	fmt.Println("GET: " + q + " Redirect to " + storage.Getrecord(q))
}

func PostHandler(w http.ResponseWriter, r *http.Request) /*(shortURL string)*/{
	// читаем Body (Тело POST запроса)
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		h :=r.Header.Get("Accept-Encoding")
		fmt.Println("Accept-Encoding= " + h)
		
	shortURLid := storage.Storerecord(string(b))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid
	//shortURLpath := BaseURL + "/?id="+ shortURLid Как сюда передать переменную из server.go?
	
	
	w.Header().Set("Content-Encoding", "gzip, deflate, br")
	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)




	//отладка что было в POST запросе
	//w.Write([]byte(b))

//типа return:
	w.Write([]byte(shortURLpath))

	fmt.Println("POST: " + string(b)+ " return id= "+ shortURLid)	
		

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
		log.Println(string(b))
		//Добавить?
		//type Example struct {
		//	URL   string `valid:"url"`
		//}
		var postJSON inputData
		err = json.Unmarshal(b, &postJSON)
		if err != nil {
			panic(err)
		}
		//отладка что было в поле url в POST запросе
		log.Println(postJSON.URL)

	shortURLid := storage.Storerecord(string(postJSON.URL))
	//shortURLpath := "http://localhost:8080/?id="+ shortURLid
	shortURLpath := os.Getenv("BASE_URL") + "/?id="+ shortURLid
	
	var shortURLpathJSON resultData
	shortURLpathJSON.ShortURL = shortURLpath

	w.Header().Set("Accept-Encoding", "gzip")
	w.Header().Set("Content-Encoding", "gzip, deflate, br")
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