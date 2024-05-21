package cache

import (
	"FirstTry/internal/cache"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Key int64 `json:"user_id"`
}

type Response struct {
	Status int      `json:"status"`
	Error  error    `json:"error"`
	Key    int64    `json:"key"`
	Value  []string `json:"value"`
}

func ShowCache(logger *slog.Logger, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var key Request
		err := render.DecodeJSON(r.Body, &key)
		if err != nil {
			logger.Error("failed to decode JSON")
			response(w, r, http.StatusBadRequest, err, 0, nil)
		}
		value, ok := cache.ShowCache(key.Key)
		if !ok {
			logger.Info("key is not found")
			response(w, r, http.StatusOK, nil, key.Key, value)
		}

		logger.Info("key found")
		response(w, r, http.StatusOK, nil, key.Key, value)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error, key int64, value []string) {
	render.JSON(w, r, Response{
		Status: status,
		Error:  err,
		Key:    key,
		Value:  value,
	})
}
