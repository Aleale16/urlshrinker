package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	//"github.com/Aleale16/urlshrinker/internal/app/handler"
	//"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	//"github.com/Aleale16/urlshrinker/internal/app/storage"
	"urlshrinker/internal/app/handler"
	"urlshrinker/internal/app/initconfig"
	"urlshrinker/internal/app/storage"

	//	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
)

/*
	type ServerConfig struct {
		SrvAddress string `env:"SERVER_ADDRESS"`
		BaseURL    string `env:"BASE_URL"`
		fileDBpath    string `env:"FILE_STORAGE_PATH"`
		User       string `env:"USERNAME"`
	}
*/
func Start() {
	var onlyOnce sync.Once
	//var SrvConfig ServerConfig
	//var UserName string
	//storage.Initdb() //Убрали управление инициализацией хранилища отсюда в storage

	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "gzip"))

	r.Get("/{id}", handler.GetHandler)
	r.Get("/api/user/urls", handler.GetUsrURLsHandler)
	r.Get("/ping", handler.GetPingHandler)

	r.Post("/", handler.PostHandler)
	r.Post("/api/shorten", handler.PostJSONHandler)
	r.Post("/api/shorten/batch", handler.PostJSONbatchHandler)

	r.Delete("/api/user/urls", handler.DeleteURLsHandler)

	r.Get("/health-check", handler.StatusOKHandler)

	fmt.Println()
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
	if initconfig.BaseURL == "" {
		//нет ни переменной окружения ни флага
		initconfig.BaseURL = "http://localhost:8080"
		log.Print("BASE_URL: " + "Loaded default: " + initconfig.BaseURL)
	}
	log.Println("BASE_URL: " + initconfig.BaseURL)

	if initconfig.FileDBpath == "" {
		//нет ни переменной окружения ни флага
		log.Print("FILE_STORAGE_PATH: not set")
	}

	if initconfig.SrvAddress == "" {
		//нет ни переменной окружения ни флага
		initconfig.SrvAddress = "localhost:8080"
		log.Print("SERVER_ADDRESS: " + "Loaded default: " + initconfig.SrvAddress)
	}

	os.Setenv("SERVER_ADDRESS", initconfig.SrvAddress)
	os.Setenv("BASE_URL", initconfig.BaseURL)
	os.Setenv("FILE_STORAGE_PATH", initconfig.FileDBpath)
	os.Setenv("DATABASE_DSN", initconfig.PostgresDBURL)
	onlyOnce.Do(storage.Initdb)

	if initconfig.SrvRunHTTPS == "HTTPS_mode_enabled" {
		log.Print("ENABLE_HTTPS: " + "HTTPS_mode_enabled")
		os.Setenv("ENABLE_HTTPS", initconfig.SrvRunHTTPS)
		// конструируем менеджер TLS-сертификатов
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist("localhost", "127.0.0.1"),
		}
		// конструируем сервер с поддержкой TLS
		server := &http.Server{
			Addr:    ":443",
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
		log.Print("ENABLE_HTTPS: " + "HTTPS_mode_enabled")
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Print("ENABLE_HTTPS: " + "Loaded default: NO HTTPS")
		//log.Fatal(http.ListenAndServe("localhost:8080", r))
		log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r))
	}

	//log.Fatal(http.ListenAndServe("localhost:8080", r))

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
