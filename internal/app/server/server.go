package server

import (
	"fmt"
	"net/http"

	"github.com/Aleale16/practicumDev/internal/app/handler"
	"github.com/Aleale16/practicumDev/internal/app/storage"
)

func Start(){
	storage.Initdb()
	
	http.HandleFunc("/health-check", handler.StatusOKHandler)

	http.HandleFunc("/", handler.ReqHandler) //Мне так не нравится, хочется тип запроса обработать уже здесь

	fmt.Println("Starting server...")
	//запуск сервера с адресом localhost, порт 8080
		server := &http.Server{
			Addr: "localhost:8080",
			//Handler: handler1,
		}
		server.ListenAndServe()
		
}