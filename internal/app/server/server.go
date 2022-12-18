package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Aleale16/urlshrinker/internal/app/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*type ServerConfig struct {
	SrvAddress string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	fileDBpath    string `env:"FILE_STORAGE_PATH"`
	User       string `env:"USERNAME"`
}*/
func Start(){

	//var SrvConfig ServerConfig
	//var UserName string
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
/*	
    err := env.Parse(&SrvConfig)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(SrvConfig)

	if (SrvConfig.User=="") {
		UserName = "Noname"
		} else {
			UserName = SrvConfig.User
		}
	log.Println("USERNAME: " + UserName)
*/
	if BaseURL == ""{
		//нет ни переменной окружения ни флага
		BaseURL = "http://localhost:8080"
		log.Print("BASE_URL: " + "Loaded default: " + BaseURL)
	}
	log.Println("BASE_URL: " + BaseURL)

	if FileDBpath == ""{
		//нет ни переменной окружения ни флага
		log.Print("FILE_STORAGE_PATH: not set")
	}

	if SrvAddress == ""{
		//нет ни переменной окружения ни флага
		SrvAddress ="localhost:8080" 
		log.Print("SERVER_ADDRESS: " + "Loaded default: " + SrvAddress)
	}

	os.Setenv("SERVER_ADDRESS", SrvAddress)
	os.Setenv("BASE_URL", BaseURL)
	os.Setenv("FILE_STORAGE_PATH", FileDBpath)

	log.Fatal(http.ListenAndServe(SrvAddress, r))

			//os.Setenv("SERVER_ADDRESS", "localhost:8080")
			//log.Print("SERVER_ADDRESS: "+"Loaded default: " + os.Getenv("SERVER_ADDRESS"))
			
			//log.Print("SERVER_ADDRESS: " + "Loaded env: " + os.Getenv("SERVER_ADDRESS"))
			//log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r))
    





	
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