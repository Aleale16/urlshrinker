package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Aleale16/practicumDev/internal/app/storage"
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
	shortURLid := storage.Storerecord(string(b))
	shortURLpath := "http://localhost:8080/?id="+ shortURLid

	// устанавливаем статус-код 201
	w.WriteHeader(http.StatusCreated)

	//отладка что было в POST запросе
	//w.Write([]byte(b))

//типа return:
	w.Write([]byte(shortURLpath))

	fmt.Println("POST: " + string(b)+ " return id= "+ shortURLid)	
		

	//return shortURLpath
}