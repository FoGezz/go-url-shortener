package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
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

func NewPostShortenHandler(app *config.App) *postShortenHandler {
	return &postShortenHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *postShortenHandler) randShortUnique(n int) string {
	alphabet := h.app.Cfg.Alphabet
	for {
		r := make([]rune, 0, n)
		for i := 0; i < n; i++ {
			randomSym := alphabet[rand.Intn(len(alphabet))]
			r = append(r, randomSym)
		}
		if _, exists := h.app.Storage.GetByShort(string(r)); !exists {
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

	if short, exists := h.app.Storage.GetByFull(full); exists {
		err := printResponse(w, req, h.app.Cfg.ResponseAddress+"/"+short, exists)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		short := h.randShortUnique(6)
		h.app.Storage.AddLink(full, short)
		err := printResponse(w, req, h.app.Cfg.ResponseAddress+"/"+short, exists)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	h.app.Storage.SaveJSONToFile(h.app.Cfg.FileStoragePath)
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

func printResponse(w http.ResponseWriter, req *http.Request, shortAddress string, conflicted bool) error {
	if req.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		if conflicted {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		err := json.NewEncoder(w).Encode(postShortenResponse{ShortURL: shortAddress})
		return err
	}
	if conflicted {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	fmt.Fprint(w, shortAddress)

	return nil
}
