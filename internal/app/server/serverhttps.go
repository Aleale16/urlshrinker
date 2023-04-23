package server

import (
	"log"
	"net/http"
	"os"
	"urlshrinker/internal/app/initconfig"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

func startHTTPS(r *chi.Mux) {
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
}
