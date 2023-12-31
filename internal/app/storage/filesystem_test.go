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

	//todo
	t.Skip("Do after 1-2-1")

	fullURL := "https://testurl555.xyz"
	cfg := config.Config{}
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cfg.FileStoragePath = os.TempDir() + "short-url-db-test-" + uuid.NewString() + ".json"
	st := storage.NewLinksMapping()
	app := config.NewApp(&cfg, st)

	postHandler := handler.NewPostShortenHandler(app)
	getHandler := handler.NewGetURLHandler(app)
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
	response, err := client.NewRequest().SetBody(fullURL).Post(postServer.URL)
	require.NoError(t, err)
	shortened := response.Body()
	require.NotEmpty(t, shortened, shortened)

	require.FileExists(t, cfg.FileStoragePath)

	st = storage.NewLinksMapping()
	err = st.LoadFromJSONFile(cfg.FileStoragePath)
	require.NoError(t, err)

	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	response, err = client.NewRequest().Get(string(shortened))
	require.ErrorIs(t, err, resty.ErrAutoRedirectDisabled)
	redirectURL := response.Header().Get("Location")
	require.Equal(t, fullURL, redirectURL)
}
