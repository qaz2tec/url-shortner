package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
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
	// fmt.Println("Hasher: ", hasher)
	data := hasher.Sum(nil)
	// fmt.Println("Hasher data: ", data)
	hash := hex.EncodeToString(data)
	// fmt.Println("Encoded String is: ", hash)
	fmt.Println("Encoded String first 8 characters: ", hash[:8])
	return hash[:8]
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

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shorturl := createUrl(data.URL)
	// fmt.Fprintf(w, shorturl)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shorturl}

	fmt.Println(Urldb)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var (
    redirectCount int
    mu            sync.Mutex
)

func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {	
	fmt.Println(redirectCount)
	mu.Lock()
    if redirectCount >= 1 {
        mu.Unlock()
        http.NotFound(w, r) // return 404 after first call
        return
    }
    redirectCount++
    mu.Unlock()

	id := r.URL.Path[len("/redirect/"):]
	// fmt.Println("ID is: ", id)
	url, err := getUrl(id)
	// fmt.Println(url, url.OriginalUrl)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func main() {
	// fmt.Println("Url-shortner....")
	// OriginalUrl := "https://github.com/qaz2tec"
	// generateShortUrl(OriginalUrl)

	http.HandleFunc("/", handle)
	http.HandleFunc("/shorten", ShortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlHandler)

	fmt.Println("Starting the server at port 8080.....")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error in Starting server: ", err)
	}
}
