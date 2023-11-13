package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
)

type postShortenHandler struct {
	ShortenerHandler
}

func NewPostShortenHandler(storage storage.ShortenerStorage, cfg *config.Config) *postShortenHandler {
	return &postShortenHandler{
		ShortenerHandler: ShortenerHandler{
			storage: storage,
			cfg:     cfg,
		},
	}
}

func (h *postShortenHandler) randShortUnique(n int) string {
	alphabet := h.cfg.Alphabet
	for {
		r := make([]rune, 0, n)
		for i := 0; i < n; i++ {
			randomSym := alphabet[rand.Intn(len(alphabet))]
			r = append(r, randomSym)
		}
		if _, exists := h.storage.GetByShort(string(r)); !exists {
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

	if short, exists := h.storage.GetByFull(string(full)); exists {
		fmt.Fprint(w, h.cfg.ResponseAddress+"/"+string(short))
	} else {
		short := h.randShortUnique(6)
		h.storage.AddLink(string(full), short)
		fmt.Fprint(w, h.cfg.ResponseAddress+"/"+string(short))
	}
}
