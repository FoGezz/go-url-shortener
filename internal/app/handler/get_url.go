package handler

import (
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/go-chi/chi/v5"
)

type getURLHandler struct {
	ShortenerHandler
}

func NewGetURLHandler(app *config.App) *getURLHandler {
	return &getURLHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *getURLHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		log.Println("getURL: method not GET but ", req.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	if id == "" {
		log.Println("getURL: no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if url, ok := h.app.Storage.GetByShort(req.Context(), id); !ok {
		log.Println("getURL: not found by ", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Location", string(url))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
