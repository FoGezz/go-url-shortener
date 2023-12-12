package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getURLHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		status   int
		location string
	}

	storage := storage.NewLinksMapping()
	storage.AddLink("https://ya.ru", "ya")
	cfg := &config.Config{}
	cfg.ResponseAddress = "http://localhost:8080"
	cfg.ServerAddress = ":8080"
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	app := config.NewApp(cfg, storage)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pos GET /ya",
			args: args{
				w:   httptest.NewRecorder(),
				req: AddChiURLParams(httptest.NewRequest(http.MethodGet, "/ya", nil), map[string]string{"id": "ya"}),
			},
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: "https://ya.ru",
			},
		},
		{
			name: "neg POST /ya WRONG METHOD",
			args: args{
				w:   httptest.NewRecorder(),
				req: AddChiURLParams(httptest.NewRequest(http.MethodPost, "/ya", nil), map[string]string{"id": "ya"}),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "neg GET / WITHOUT ID",
			args: args{
				w:   httptest.NewRecorder(),
				req: AddChiURLParams(httptest.NewRequest(http.MethodGet, "/", nil), map[string]string{"id": ""}),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "neg GET /qwerty NONEXISTENT ID",
			args: args{
				w:   httptest.NewRecorder(),
				req: AddChiURLParams(httptest.NewRequest(http.MethodGet, "/qwerty", nil), map[string]string{"id": "qwerty"}),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGetURLHandler(app)
			h.ServeHTTP(tt.args.w, tt.args.req)
			result := tt.args.w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want.status, result.StatusCode)
			if result.StatusCode != http.StatusBadRequest {
				require.Equal(t, tt.want.location, result.Header.Get("Location"))

				//GET with the same data and ensure that result is the same
				h.ServeHTTP(tt.args.w, tt.args.req)
				newResult := tt.args.w.Result()
				defer newResult.Body.Close()
				assert.Equal(t, result.StatusCode, newResult.StatusCode)
				assert.Equal(t, result.Header.Get("Location"), newResult.Header.Get("Location"))
			}
		})
	}
}
