package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshrinker/internal/app/storage"
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

func TestReqHandlerPostJSON(t *testing.T) {
	//storage.Initdb()
	//var body = []byte("https://ya.ru")
	//Описание тела запроса в JSON
	var body = []byte(`{"url": "yajson.ru"}`)

	reqpost, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	//Передача JSON в запрос
	//reqpost, err := http.NewRequest("POST", "/", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}
	//reqpost.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostJSONHandler)
	handler.ServeHTTP(rr, reqpost)

	if contenttype := rr.Header().Get("Content-Type"); contenttype != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			contenttype, "application/json")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	fmt.Println(reqpost.URL)
	fmt.Println(rr.Body.String())
	fmt.Println(rr.Header().Get("Content-Type"))

}

func TestReqHandlerPostbatchJSON(t *testing.T) {
	//storage.Initdb()
	//var body = []byte("https://ya.ru")
	//Описание тела запроса в JSON

	type shortenRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	requestData := []shortenRequest{
		{
			CorrelationID: "123456",
			OriginalURL:   "https://yabatch.ru",
		},
	}
	//var body = []byte(`{{"correlation_id":"123456","original_url":"https://yabatch.ru"},}`)
	body, _ := json.Marshal(requestData)
	reqpost, err := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
	//Передача JSON в запрос
	//reqpost, err := http.NewRequest("POST", "/", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}
	//reqpost.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostJSONbatchHandler)
	handler.ServeHTTP(rr, reqpost)

	if contenttype := rr.Header().Get("Content-Type"); contenttype != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			contenttype, "application/json")
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	fmt.Println(reqpost.URL)
	fmt.Println(rr.Body.String())
	fmt.Println(rr.Header().Get("Content-Type"))

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
	if status != http.StatusTemporaryRedirect {
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
	if status != http.StatusBadRequest {
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
	if status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	fmt.Println(reqget.URL)
	fmt.Println(status)

}

func TestReqHandlerGetURLforUser(t *testing.T) {
	//id отсутствует
	reqget, err := http.NewRequest("GET", "/api/user/urls", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	reqget.Header.Set("Authorization", "15b94b695561803cbf3bd2ef218518b3fce9661d0eba8ddf23fcd6deb556d0a939393939")
	handler := http.HandlerFunc(GetUsrURLsHandler)
	handler.ServeHTTP(rr, reqget)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if contenttype := rr.Header().Get("Content-Type"); contenttype != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			contenttype, "application/json")
	}
	fmt.Println(reqget.URL)
	fmt.Println(status)

}
/*
func TestGetPingHandler(t *testing.T) {
	// type args struct {
	// 	w http.ResponseWriter
	// 	r *http.Request
	// }
	type want struct {
        //contentType string
        statusCode  int
    }
	tests := []struct {
		name string
		//args args
        want    want
	}{
		// TODO: Add test cases.
		{
			name: "connection is OK",
			want: want{
                //contentType: "application/json",
                statusCode:  200,
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, "/", nil)

            // создаём новый Recorder
            w := httptest.NewRecorder()
            // определяем хендлер
            h := http.HandlerFunc(GetPingHandler)
            // запускаем сервер
            h.ServeHTTP(w, request)
            res := w.Result()

            // проверяем код ответа
            if res.StatusCode != tt.want.statusCode {
                t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
            }
			
		})
	}
}
*/
func TestReqHandlerDelURLforUser(t *testing.T) {

	var body = []byte(`["222"]`)
	reqget, err := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	reqget.Header.Set("Authorization", "15b94b695561803cbf3bd2ef218518b3fce9661d0eba8ddf23fcd6deb556d0a939393939")

	handler := http.HandlerFunc(DeleteURLsHandler)
	handler.ServeHTTP(rr, reqget)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	fmt.Println(reqget.URL)
	fmt.Println(status)

}

func TestSign(t *testing.T) {
	validSign, val := checkSign("15b94b695561803cbf3bd2ef218518b3fce9661d0eba8ddf23fcd6deb556d0a939393939")
	if validSign != true {
		t.Errorf("Wrong Sign state: got %v want %v",
			validSign, true)
	}
	if val != "9999" {
		t.Errorf("Wrong id retrieved: got %v want %v",
			val, "9999")
	}
	fmt.Println(validSign)
	fmt.Println(val)

}
func TestSignOptimized(t *testing.T) {	
	validSignopt, valopt := checkSignOptimized("15b94b695561803cbf3bd2ef218518b3fce9661d0eba8ddf23fcd6deb556d0a939393939")
	if validSignopt != true {
		t.Errorf("Wrong SignOptimized state: got %v want %v",
		validSignopt, true)
	}
	if valopt != "9999" {
		t.Errorf("Wrong id retrieved: got %v want %v",
		valopt, "9999")
	}
	fmt.Println(validSignopt)
	fmt.Println(valopt)
}
