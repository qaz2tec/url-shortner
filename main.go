package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Url struct {
	ID          string    `json:"id"`
	OriginalUrl string    `json:"originalUrl"`
	ShortUrl    string    `json:"shortUrl"`
	Created_at  time.Time `json:"created_at"`
}

var Urldb = make(map[string]Url)

func generateShortUrl(OriginalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl))
	fmt.Println("Hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("Hasher data: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Encoded String is: ", hash)
	fmt.Println("Encoded String first 8 characters: ", hash[:8])
	return "Hey"
}

func createUrl(originalUrl string) string {
	shortUrl := generateShortUrl(originalUrl)
	id := shortUrl
	Urldb[id] = Url{
		ID:          id,
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		Created_at:  time.Now(),
	}
	return shortUrl
}

func getUrl(id string) (Url, error) {
	url, ok := Urldb[id]
	if !ok {
		return Url{}, errors.New("URL not found")
	}
	return url, nil
}

func handle(w http.ResponseWriter, r *http.Request) { //Sprintf
	fmt.Fprintf(w, "Hello World")
}

func main() {
	fmt.Println("Url-shortner....")
	OriginalUrl := "https://github.com/qaz2tec"
	generateShortUrl(OriginalUrl)

	http.HandleFunc("/", handle)

	fmt.Println("Starting the server at port 8080.....")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error in Starting server: ", err)
	}
}
