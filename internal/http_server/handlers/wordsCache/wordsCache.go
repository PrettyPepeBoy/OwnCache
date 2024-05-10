package wordsCache

import (
	"FirstTry/internal/cache"
	"encoding/json"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type letter struct {
	Ltr string `json:"letter"`
}

func AddInCacheLetter(logger *slog.Logger, cache *cache.CacheForWords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		var ltr letter
		err = json.Unmarshal(rawByte, &ltr)
		if err != nil {
			logger.Error("failed to unmarshall json request", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		cache.Put(ltr.Ltr[0])
		logger.Info("successfully added in cache")
		response(w, r, http.StatusOK, nil)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
