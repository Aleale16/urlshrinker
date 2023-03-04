package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aleale16/urlshrinker/internal/app/storage"
)

func BenchmarkServerstart(b *testing.B) {
	b.Run("variant01", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			checkSign("982c62b2730d5ffe217d69e8e82fd1e0aa4e0154a80323ea20678189b14b1e1d75736572333333")
		}
	})
	b.Run("variant02", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			checkSignOptimized("982c62b2730d5ffe217d69e8e82fd1e0aa4e0154a80323ea20678189b14b1e1d75736572333333")
		}
	})
} 

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StatusOKHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	fmt.Println(req.URL)
}

func TestReqHandlerPost(t *testing.T) {
	storage.Initdb()
	var body = []byte("https://ya.ru")	
	//Описание тела запроса в JSON
	//var body = []byte(`{"message": "mail.ru"}`)

	reqpost, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	//Передача JSON в запрос
	//reqpost, err := http.NewRequest("POST", "/", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}	
	//reqpost.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostHandler)
	handler.ServeHTTP(rr, reqpost)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	//expected := "http://localhost:8080/?id=..."
	expected := "http://localhost:8080/..."
	if rr.Body.String() == "" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	fmt.Println(reqpost.URL)
	fmt.Println(rr.Body.String())
	
}

func TestReqHandlerGet1(t *testing.T) {
//Запрос существующего id
	//reqget, err := http.NewRequest("GET", "/?id=7943", nil)
	reqget, err := http.NewRequest("GET", "/111", nil)
	if err != nil {
		t.Fatal(err)
	}	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHandler)
	handler.ServeHTTP(rr, reqget)

	// Check the status code is what we expect.
	status := rr.Code
	if  status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	// Check the response Header Location is what we expect.
	expected := "http://URL..."
	header := rr.Header()
	if header.Get("Location") == "" {
		t.Errorf("handler returned unexpected header: got %v want %v",
		header.Get("Location"), expected)
	}

	fmt.Println(reqget.URL)
	fmt.Println(status)
	fmt.Println(header.Get("Location"))
	
}

func TestReqHandlerGet2(t *testing.T) {
//запрос несуществующего id
	//reqget, err := http.NewRequest("GET", "/?id=xxxx", nil)
	reqget, err := http.NewRequest("GET", "/xxxx", nil)
	if err != nil {
		t.Fatal(err)
	}	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHandler)
	handler.ServeHTTP(rr, reqget)
	status := rr.Code
	if  status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected := "http://google.com/404"
	header := rr.Header()
	if header.Get("Location") != expected {
		t.Errorf("handler returned unexpected header: got %v want %v",
		header.Get("Location"), expected)
	}

	fmt.Println(reqget.URL)
	fmt.Println(status)
	fmt.Println(header.Get("Location"))
	
}

func TestReqHandlerGet3(t *testing.T) {
//id отсутствует
	reqget, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHandler)
	handler.ServeHTTP(rr, reqget)
	status := rr.Code
	if  status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	fmt.Println(reqget.URL)
	fmt.Println(status)
	
}