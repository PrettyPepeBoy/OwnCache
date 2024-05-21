package wordsCache

import (
	"FirstTry/internal/cache"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

type Word struct {
	ProductName string `json:"product_name"`
}

func AddInCacheWord(logger *slog.Logger, cache *cache.WordsCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		if err != nil {
			logger.Error("failed to unmarshall json request", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}

		cache.PutWord(rawByte)

		logger.Info("successfully added in cache")
		response(w, r, http.StatusOK, nil)
	}
}

func GetWord(logger *slog.Logger, cache *cache.WordsCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		var word Word
		err = json.Unmarshal(rawByte, &word)
		if err != nil {
			logger.Error("failed to unmarshall json request", err)
			response(w, r, http.StatusInternalServerError, err)
			return
		}
		slc := cache.GetWord([]byte(word.ProductName))
		logger.Info(fmt.Sprintf("succesfully get word %s from %s", slc, word.ProductName))
		response(w, r, http.StatusOK, nil)
	}
}

func response(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.JSON(w, r, Response{
		StatusResp: status,
		Err:        err,
	})
}
