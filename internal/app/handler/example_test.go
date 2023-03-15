package handler_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Example() {
	// Как CheckSign разбирает токен авторизации msg
	msg := "15b94b695561803cbf3bd2ef218518b3fce9661d0eba8ddf23fcd6deb556d0a939393939"
	var validSign bool
	var val string
	var key = []byte("secret key")
	var (
		data []byte // декодированное сообщение с подписью
		id   string // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)
	validSign = false
	data, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	//fmt.Println("data=" + string(data))
	id = string(data[sha256.Size:])
	val = id
	//id = binary.BigEndian.Uint32(data[:4])
	//id = binary.BigEndian.Uint32(data[sha256.Size:])
	h := hmac.New(sha256.New, key)
	h.Write(data[sha256.Size:])
	sign = h.Sum(nil)
	if hmac.Equal(sign, data[:sha256.Size]) {
		//fmt.Println("Подпись подлинная. ID:", id)
		validSign = true
	} /*else {
		fmt.Println("Подпись неверна. Где-то ошибка! ID:", id)
	}	*/
	fmt.Printf("Подпись верна?: %v. ID пользователя: %v", validSign, val)

}
