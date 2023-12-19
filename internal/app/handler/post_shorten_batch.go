package handler

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
)

type postShortenBatchRequestUnit struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type postShortenBatchResponseUnit struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type postShortenBatchHandler struct {
	ShortenerHandler
}

func NewPostShortenBatchHandler(app *config.App) *postShortenBatchHandler {
	return &postShortenBatchHandler{
		ShortenerHandler: ShortenerHandler{
			app,
		},
	}
}

func (h *postShortenBatchHandler) randShortUnique(ctx context.Context, n int) string {
	alphabet := h.app.Cfg.Alphabet
	for {
		r := make([]rune, 0, n)
		for i := 0; i < n; i++ {
			randomSym := alphabet[rand.Intn(len(alphabet))]
			r = append(r, randomSym)
		}
		if _, exists := h.app.Storage.GetByShort(ctx, string(r)); !exists {
			return string(r)
		}
	}

}

func (h *postShortenBatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()

	r := json.NewDecoder(req.Body)
	var parsedReq []postShortenBatchRequestUnit
	var resp []postShortenBatchResponseUnit
	if err := r.Decode(&parsedReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range parsedReq {

		s, found := h.app.Storage.GetByFull(req.Context(), v.OriginalURL)
		if !found {
			s = h.randShortUnique(req.Context(), 6)
			h.app.Storage.AddLink(req.Context(), v.OriginalURL, s)
		}
		resp = append(resp, postShortenBatchResponseUnit{ShortURL: h.app.Cfg.ResponseAddress + "/" + s, CorrelationID: v.CorrelationID})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)

	h.app.Storage.SaveJSONToFile(h.app.Cfg.FileStoragePath)
}
