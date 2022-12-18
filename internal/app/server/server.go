package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Aleale16/urlshrinker/internal/app/handler"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type ServerConfig struct {
	SrvAddress string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	fileDBpath    string `env:"FILE_STORAGE_PATH"`
	User       string `env:"USERNAME"`
}
func Start(){

	var SrvConfig ServerConfig
	var baseURL, UserName, fileDBpath string
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
	
    err := env.Parse(&SrvConfig)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(SrvConfig)

	_, baseURLexists := os.LookupEnv("BASE_URL")
	_, srvAddressexists := os.LookupEnv("SERVER_ADDRESS")
	_, fileDBpathexists := os.LookupEnv("FILE_STORAGE_PATH")

	if (SrvConfig.User=="") {
		UserName = "Noname"
		} else {
			UserName = SrvConfig.User
		}
	log.Println("USERNAME: " + UserName)

	if baseURLexists {
		baseURL = SrvConfig.BaseURL
		//baseURL = os.Getenv("BASE_URL")
		} else {
			//запишем в переменную значение по умолчанию
			//os.Setenv("BASE_URL", "http://localhost:8080")
			//baseURL = os.Getenv("BASE_URL")
			//Возьмем значение флага
			if BaseURL != ""{
				baseURL = BaseURL
			} else {//нет ни переменной окружения ни флага
				baseURL = "http://localhost:8080"
			}
    	}
	log.Println("BASE_URL: " + baseURL)

	if fileDBpathexists {
		fileDBpath = SrvConfig.fileDBpath
		//fileDBpath = os.Getenv("FILE_STORAGE_PATH")
		log.Print("FILE_STORAGE_PATH: " + "Loaded env: " + fileDBpath)
		} else {
			//Возьмем значение флага
			if FileDBpath != ""{
				fileDBpath = FileDBpath
				log.Print("fileDBpath Loaded from flag: " + fileDBpath)
			} else {//нет ни переменной окружения ни флага
				log.Print("FILE_STORAGE_PATH: not set")
			}
			
    	}

	if srvAddressexists {
		log.Print("SERVER_ADDRESS: " + "Loaded env: " + SrvConfig.SrvAddress)

		log.Fatal(http.ListenAndServe(SrvConfig.SrvAddress, r))
		} else {
			var srvaddr string
			//Возьмем значение флага
			if SrvAddress != ""{
				srvaddr = SrvAddress
				log.Print("SERVER_ADDRESS: "+"Loaded from flag: " + srvaddr)
			} else { //нет ни переменной окружения ни флага
				srvaddr ="localhost:8080" 
				log.Print("SERVER_ADDRESS: "+"Loaded default: " + srvaddr)
			}
			log.Fatal(http.ListenAndServe(srvaddr, r))
			//os.Setenv("SERVER_ADDRESS", "localhost:8080")
			//log.Print("SERVER_ADDRESS: "+"Loaded default: " + os.Getenv("SERVER_ADDRESS"))
			
			//log.Print("SERVER_ADDRESS: " + "Loaded env: " + os.Getenv("SERVER_ADDRESS"))
			//log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r))
    }





	
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