package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShorturl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) //It convert the orignalURL stirng to a byte slice
	fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data :", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Encode to String ", hash)
	fmt.Println("final string :", hash[:8])
	return hash[:8]

}
func createURL(originalURL string) string {
	shortURL := generateShorturl(originalURL)
	id := shortURL // Use the short URl as the ID for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL // return the newly created short URL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World !... ")
}
func shortURLhandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	//fmt.Fprintf(w, shortURL);
	response := struct {
		ShortURL string `json:"shorturl"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "applicationjson")
	json.NewEncoder(w).Encode(response)

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	/* fmt.Println("Starting URL shortener...")
	orignalURL := "https://github.com/Sankalp-Space"
	generateShorturl(orignalURL) */

	//Register the handler functionto handle all request to the root URL ("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", shortURLhandler)
	http.HandleFunc("/redirect/", redirectHandler)
	//start the HTTP server on port 3000
	fmt.Println("Starting server on port 3000.....")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on starting server :..", err)

	}
}
