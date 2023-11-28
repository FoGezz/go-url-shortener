package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func Test_postShortenHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		status int
	}

	storage := storage.NewLinksMapping()
	cfg := &config.Config{}
	cfg.ResponseAddress = "http://localhost:8080"
	cfg.ServerAddress = "localhost:8080"
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pos POST / {https://ya.ru}",
			args: args{
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru")),
			},
			want: want{
				status: http.StatusCreated,
			},
		},
		{
			name: "neg GET / {https://ya.ru}",
			args: args{
				w:   httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/", strings.NewReader("https://ya.ru")),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPostShortenHandler(storage, cfg)
			h.ServeHTTP(tt.args.w, tt.args.req)
			result := tt.args.w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want.status, result.StatusCode)
			tt.args.w.Flush()
			if result.StatusCode != http.StatusBadRequest {
				bodyStr, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				require.NotEmpty(t, bodyStr)
			}
		})
	}
}

func Test_postApiShortenHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		status int
	}

	storage := storage.NewLinksMapping()
	cfg := &config.Config{}
	cfg.ResponseAddress = "http://localhost:8080"
	cfg.ServerAddress = "localhost:8080"
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pos POST /api/shorten {https://ya.ru}",
			args: args{
				w:   httptest.NewRecorder(),
				req: WithAddedHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("{\"url\":\"https://ya.ru\"}")), "Content-Type", "application/json"),
			},
			want: want{
				status: http.StatusCreated,
			},
		},
		{
			name: "neg POST /api/shorten {https://ya.ru}",
			args: args{
				w:   httptest.NewRecorder(),
				req: WithAddedHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("https://ya.ru")), "Content-Type", "application/json"),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "BAD JSON POST /api/shorten {https://ya.ru}",
			args: args{
				w:   httptest.NewRecorder(),
				req: WithAddedHeader(httptest.NewRequest(http.MethodPost, "/api/shorten/", strings.NewReader("hello: 'ping'}")), "Content-Type", "application/json"),
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPostShortenHandler(storage, cfg)
			h.ServeHTTP(tt.args.w, tt.args.req)
			result := tt.args.w.Result()
			defer result.Body.Close()
			require.Equal(t, tt.want.status, result.StatusCode)
			tt.args.w.Flush()
			if result.StatusCode != http.StatusBadRequest {
				resp := new(postShortenResponse)
				err := json.NewDecoder(result.Body).Decode(resp)
				require.NoError(t, err)
				require.NotEmpty(t, resp.ShortURL)
			}
		})
	}
}

func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	req := r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	return req
}

func WithAddedHeader(r *http.Request, key string, value string) *http.Request {
	r.Header.Add(key, value)
	return r
}
