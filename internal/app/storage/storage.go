package storage

import (
	"fmt"
	"math/rand"
	"strconv"
)

type URLrecord map[string]string

var URL URLrecord

func Initdb() {
	URL = make(URLrecord)
	fmt.Println("Storage ready!")
}

func Storerecord(fullURL string) string{
	//Once.Do(Initdb())
	id := strconv.Itoa(rand.Intn(9999))
	
	for (!isnewID(id)){
		id = strconv.Itoa(rand.Intn(9999))
	}
		URL[id] = fullURL
		return id
}

func Getrecord(id string) string {
	result := URL[id]

	if (result != ""){
		return result
	} else {
		return "http://google.com/404"
	}
}

func isnewID(id string) bool{
	result := URL[id]
	if (result == ""){
		return true
	} else {return false}
}