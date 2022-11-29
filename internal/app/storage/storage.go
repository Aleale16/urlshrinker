package storage

import (
	"math/rand"
	"strconv"
)

type URLrecord map[string]string

var URL URLrecord

func Initdb() {
	URL = make(URLrecord)
}

func Storerecord(fullURL string) string{
	id := strconv.Itoa(rand.Intn(9999))
	URL[id] = fullURL
	return id
}

func Getrecord(id string) string {
	result := URL[id]
	
	if (result != ""){
		return result
	} else {
		return "error404"
	}
}