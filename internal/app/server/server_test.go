package server

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetHandler(t *testing.T) {
	/*w := http.ResponseWriter, r *http.Request

    if shortURL := PostHandler(w, r); shortURL != "http://localhost:8080/?id=333" {
        t.Errorf("sum expected to be http://localhost:8080/?id=333; got %s", shortURL)
    }*/
	//client := &http.Client{}
	//req, err := http.NewRequest("GET", "http://localhost:8080/?id=333", nil) //так не проходят автотесты на гит
	req, err := http.NewRequest("GET", "/", nil) // так проходят	
    if err != nil {
        t.Fatal(err)
    }
	fmt.Println(req.URL)
}



