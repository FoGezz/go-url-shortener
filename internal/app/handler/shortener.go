package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

var alphabet []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ShortenerHandler struct {
	http.Handler
	linksContainer *storage.LinksContainer
}

type postShortenHandler struct {
	ShortenerHandler
}

func NewPostShortenHandler(container *storage.LinksContainer) *postShortenHandler {
	s := &postShortenHandler{} //у меня не получилось здесь сразу передать linksContainer
	s.linksContainer = container

	return s
}

func (h *postShortenHandler) randShortUnique(n int) string {
	for {
		r := make([]rune, 0, n)
		for i := 0; i < n; i++ {
			randomSym := alphabet[rand.Intn(len(alphabet))]
			r = append(r, randomSym)
		}
		if _, exists := h.linksContainer.GetByShort(string(r)); !exists {
			return string(r)
		}
	}

}

func (h *postShortenHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Println("postShorten: method not POST but ", req.Method)
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

	if short, exists := h.linksContainer.GetByFull(string(full)); exists {
		fmt.Fprint(w, "http://"+req.Host+"/"+string(short))
	} else {
		short := h.randShortUnique(6)
		h.linksContainer.AddLink(string(full), short)
		fmt.Fprint(w, "http://"+req.Host+"/"+string(short))
	}
}

type getURLHandler struct {
	ShortenerHandler
}

func NewGetURLHandler(container *storage.LinksContainer) *getURLHandler {
	s := &getURLHandler{} //у меня не получилось здесь сразу передать linksContainer
	s.linksContainer = container

	return s
}

func (h *getURLHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("getURL: method not GET but ", req.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, exist := mux.Vars(req)["id"]
	if !exist {
		log.Println("getURL: no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url, ok := h.linksContainer.GetByShort(id); !ok {
		log.Println("getURL: not found by ", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", string(url))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
