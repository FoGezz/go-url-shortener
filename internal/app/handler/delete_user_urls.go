package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
)

type deleteUserURLsHandler struct {
	ShortenerHandler
}

func NewDeleteUserURLsHandler(app *config.App) *deleteUserURLsHandler {
	return &deleteUserURLsHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *deleteUserURLsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	userUUIDAny := req.Context().Value(middleware.UserIDKey)
	if userUUIDAny == nil || userUUIDAny == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID := userUUIDAny.(string)

	var shortURLs []string
	err := json.NewDecoder(req.Body).Decode(&shortURLs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.app.Storage.DeleteAsync(context.TODO(), shortURLs, userUUID)

	w.WriteHeader(http.StatusAccepted)
}
