package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type fullUrl string
type shortUrl string

type shortToFullMap map[shortUrl]fullUrl
type fullToShortMap map[fullUrl]shortUrl

var byShortMap shortToFullMap
var byFullMap fullToShortMap
var alphabet []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	byShortMap = make(shortToFullMap, 0)
	byFullMap = make(fullToShortMap, 0)
	mux := mux.NewRouter()
	mux.HandleFunc("/", postShorten)
	mux.HandleFunc("/{id}", getUrl)
	http.ListenAndServe(":8080", mux)
}

func postShorten(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	full, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	if short, exists := byFullMap[fullUrl(full)]; exists {
		fmt.Fprint(w, req.Host+"/"+string(short))
	} else {
		short := randShortUnique(6)
		addToMaps(shortUrl(short), fullUrl(full))
		fmt.Fprint(w, req.Host+"/"+string(short))
	}
}

func getUrl(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url, ok := byShortMap[shortUrl(id)]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", string(url))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func randShortUnique(n int) shortUrl {

	for {
		r := make([]rune, 0, n)
		for i := 0; i < n; i++ {
			randomSym := alphabet[rand.Intn(len(alphabet))]
			r = append(r, randomSym)
		}
		if _, exists := byShortMap[shortUrl(r)]; !exists {
			return shortUrl(r)
		}
	}

}

func addToMaps(s shortUrl, f fullUrl) {
	byShortMap[s] = f
	byFullMap[f] = s
}
