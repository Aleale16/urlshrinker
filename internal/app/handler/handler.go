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

func ReqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
        postHandler(w, r)
    }
	if r.Method == http.MethodGet {
        getHandler(w, r)
    }
	fmt.Println(r.Method)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
        http.Error(w, "Only Get requests are allowed!", http.StatusBadRequest)
        return
    }
	
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
	//w.WriteHeader(http.StatusOK)

	fmt.Println("GET: " + q)
}

func postHandler(w http.ResponseWriter, r *http.Request) (shortURL string){
	// этот обработчик принимает только запросы, отправленные методом POST и GET
	if r.Method != http.MethodPost {
        http.Error(w, "Only Post requests are allowed!", http.StatusBadRequest)
        return
    }

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

	w.Write([]byte(shortURLpath))

	fmt.Println("POST: " + string(b)+ " return id= "+ shortURLid)	

	return shortURLpath
}