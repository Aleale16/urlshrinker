package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RAM storage type.
func TestStorerecord(t *testing.T) {
	URL = make(URLrecord)
	RAMonly = true
	SetdbType()
	type want struct {
		shortid    string
		statusCode string
	}

	tests := []struct {
		name    string
		FullURL string
		want    want
	}{
		{
			name:    "store new url",
			FullURL: "https://ya.ru",
			want: want{
				shortid:    "111",
				statusCode: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortid, statusCode := Storerecord(tt.FullURL)
			// для сравнения двух чисел подойдет функция Equal
			assert.Equal(t, tt.want.shortid, shortid)
			assert.Equal(t, tt.want.statusCode, statusCode)
		})
	}

}

func TestGetrecord(t *testing.T) {
	type want struct {
		FullURL    string
		statusCode string
	}
	tests := []struct {
		name    string
		shortid string
		want    want
	}{
		{
			name:    "get existing url",
			shortid: "111",
			want: want{
				FullURL:    "https://ya.ru",
				statusCode: "307",
			},
		},
		{
			name:    "get not existing url",
			shortid: "/xxxx",
			want: want{
				FullURL:    "http://google.com/404",
				statusCode: "400",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FullURL, statusCode := Getrecord(tt.shortid)
			// для сравнения двух чисел подойдет функция Equal
			assert.Equal(t, tt.want.FullURL, FullURL)
			assert.Equal(t, tt.want.statusCode, statusCode)
		})
	}
}

func TestStoreShortURLtouser(t *testing.T) {

	Usr = make(Userrecord)
	RAMonly = true
	SetdbType()
	type want struct {
		JSONresponse string
		noURLs       bool
	}

	tests := []struct {
		name    string
		userID  string
		shortID string
		want    want
	}{
		{
			name:    "store new url",
			userID:  "9999",
			shortID: "111",
			want: want{
				JSONresponse: "[{\"short_url\":\"/111\",\"original_url\":\"https://ya.ru\"}]",
				noURLs:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssignShortURLtouser(tt.userID, tt.shortID)
		})
		t.Run(tt.name, func(t *testing.T) {
			output, noURLs, outputURLs := GetuserURLS(tt.userID)
			fmt.Printf("%s %v %s", output, noURLs, outputURLs)
			assert.Equal(t, tt.want.JSONresponse, output)
			assert.Equal(t, tt.want.noURLs, noURLs)
		})
	}

}
