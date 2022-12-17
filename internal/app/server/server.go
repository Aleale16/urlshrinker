package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Aleale16/urlshrinker/internal/app/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start(){
	//storage.Initdb() //Убрали управление инициализацией хранилища отсюда в storage

	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.Get("/", handler.GetHandler)
	r.Post("/", handler.PostHandler)
	r.Post("/api/shorten", handler.PostJSONHandler)
	//r.Get("/health-check", handler.StatusOKHandler)
	

	fmt.Println("Starting server...")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
	
	/*http.HandleFunc("/health-check", handler.StatusOKHandler)

	http.HandleFunc("/", handler.ReqHandler) //Мне так не нравится, хочется тип запроса обработать уже здесь.
											//Для этого есть методы в роутере chi

	fmt.Println("Starting server...")
	//запуск сервера с адресом localhost, порт 8080
		server := &http.Server{
			Addr: "localhost:8080",
			//Handler: handler1,
		}
		server.ListenAndServe()*/
		
}