package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	//"github.com/Aleale16/urlshrinker/internal/app/handler"
	//"github.com/Aleale16/urlshrinker/internal/app/initconfig"
	//"github.com/Aleale16/urlshrinker/internal/app/storage"
	"urlshrinker/internal/app/handler"
	"urlshrinker/internal/app/initconfig"
	"urlshrinker/internal/app/storage"

	//	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
	type ServerConfig struct {
		SrvAddress string `env:"SERVER_ADDRESS"`
		BaseURL    string `env:"BASE_URL"`
		fileDBpath    string `env:"FILE_STORAGE_PATH"`
		User       string `env:"USERNAME"`
	}
*/

const shutdownTimeout = 5 * time.Second

func Start(ctx context.Context) error {
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
	r.Get("/api/internal/stats", handler.GetStatsHandler)

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
	onlyOnce.Do(storage.Initdb)

	if initconfig.SrvRunHTTPS == "HTTPS_mode_enabled" {
		startHTTPS(r)
	} else {
		log.Print("ENABLE_HTTPS: " + "Loaded default: NO HTTPS")
		//log.Fatal(http.ListenAndServe("localhost:8080", r))
		//log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r))
		var srv = &http.Server{
			Addr:    os.Getenv("SERVER_ADDRESS"),
			Handler: r,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Listen and serve: %v", err)
			}
		}()

		log.Printf("Listening on %s", srv.Addr)
		<-ctx.Done()

		log.Println("Shutting down server gracefully")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}

		//Если какая-то из операций по очистке ресурсов повисла
		longShutdown := make(chan struct{}, 1)

		go func() {
			time.Sleep(3 * time.Second)
			longShutdown <- struct{}{}
		}()

		select {
		case <-shutdownCtx.Done():
			return fmt.Errorf("server shutdown: %w", ctx.Err())
		case <-longShutdown:
			log.Println("Finished")
		}
	}
	return nil
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
