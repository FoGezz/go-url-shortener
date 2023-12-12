package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
)

type getPingHandler struct {
	ShortenerHandler
}

func NewGetPingHandler(app *config.App) *getPingHandler {
	return &getPingHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *getPingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.app.DBPool != nil {
		err := h.app.DBPool.Ping(context.Background())
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
}
