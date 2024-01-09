package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
)

type getUserURLsHandler struct {
	ShortenerHandler
}

type postGetUserURLsResponseUnit struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewGetUserURLsHandler(app *config.App) *getUserURLsHandler {
	return &getUserURLsHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *getUserURLsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	userUUIDAny := req.Context().Value(middleware.USER_ID_KEY)
	if userUUIDAny == nil || userUUIDAny == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userUUID := userUUIDAny.(string)

	links, err := h.app.Storage.GetByUserUUID(req.Context(), userUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(*links) == 0 {
		log.Println(userUUIDAny)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := make([]postGetUserURLsResponseUnit, 0, len(*links))
	for short, long := range *links {
		result = append(result, postGetUserURLsResponseUnit{
			ShortURL:    h.app.Cfg.ResponseAddress + "/" + string(short),
			OriginalURL: string(long),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
