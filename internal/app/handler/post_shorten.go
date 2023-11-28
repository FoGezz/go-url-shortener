package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
)

type postShortenRequest struct {
	FullURL string `json:"url"`
}

type postShortenResponse struct {
	ShortURL string `json:"result"`
}

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
	defer func() {
		_ = req.Body.Close()
	}()
	if req.Method != http.MethodPost {
		log.Println("postShorten: method not POST but ", req.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	full, err := parseFullFromRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")

	if short, exists := h.storage.GetByFull(full); exists {
		err := printResponse(w, req, h.cfg.ResponseAddress+"/"+short)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		short := h.randShortUnique(6)
		h.storage.AddLink(full, short)
		err := printResponse(w, req, h.cfg.ResponseAddress+"/"+short)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	h.storage.SaveJsonToFile(h.cfg.FileStoragePath)
}

func parseFullFromRequest(req *http.Request) (string, error) {
	if req.Header.Get("Content-Type") == "application/json" {
		jsonReq := new(postShortenRequest)
		decodeEdd := json.NewDecoder(req.Body).Decode(jsonReq)
		if decodeEdd != nil {
			return "", decodeEdd
		}
		return jsonReq.FullURL, nil
	}
	full, err := io.ReadAll(req.Body)
	return string(full), err
}

func printResponse(w http.ResponseWriter, req *http.Request, shortAddress string) error {
	if req.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(postShortenResponse{ShortURL: shortAddress})
		return err
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortAddress)

	return nil
}
