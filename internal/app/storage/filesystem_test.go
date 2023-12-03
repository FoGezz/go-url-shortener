package storage_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

/*
*
Big system test of saving file from http on working server
*/
func TestLinksMapping_JSONFile_SaveAndLoad(t *testing.T) {
	fullUrl := "https://testurl555.xyz"
	cfg := config.Config{}
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cfg.FileStoragePath = os.TempDir() + "short-url-db-test-" + uuid.NewString() + ".json"
	st := storage.NewLinksMapping()

	postHandler := handler.NewPostShortenHandler(st, &cfg)
	getHandler := handler.NewGetURLHandler(st, &cfg)
	router := chi.NewRouter()
	router.Handle("/{id}", getHandler)
	router.Handle("/", postHandler)

	postServer := httptest.NewServer(router)
	defer postServer.Close()

	getServer := httptest.NewServer(router)
	defer getServer.Close()

	cfg.ServerAddress = postServer.URL
	cfg.ResponseAddress = getServer.URL

	client := resty.New()
	response, err := client.NewRequest().SetBody(fullUrl).Post(postServer.URL)
	require.NoError(t, err)
	shortened := response.Body()
	require.NotEmpty(t, shortened, shortened)

	require.FileExists(t, cfg.FileStoragePath)

	st = storage.NewLinksMapping()
	st.LoadFromJSONFile(cfg.FileStoragePath)

	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	response, err = client.NewRequest().Get(string(shortened))
	require.ErrorIs(t, err, resty.ErrAutoRedirectDisabled)
	redirectUrl := response.Header().Get("Location")
	require.Equal(t, fullUrl, redirectUrl)
}
