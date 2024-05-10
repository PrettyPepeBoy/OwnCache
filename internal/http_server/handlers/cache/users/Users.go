package users

import (
	"FirstTry/internal/cache"
	"encoding/json"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type usersAndProductId struct {
	UserId    int64 `json:"user_id"`
	ProductId int64 `json:"product_id"`
}

type Response struct {
	StatusResp int   `json:"statusResp"`
	Err        error `json:"err,omitempty"`
}

func AddInCacheUsersAndProductId(logger *slog.Logger, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var usrProduct usersAndProductId
		rawByte, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			response(w, r, http.StatusBadRequest, err)
			return
		}
		logger.Info("raw request", string(rawByte))
		err = json.Unmarshal(rawByte, &usrProduct)
		if err != nil {
			logger.Error("failed to decode JSON")
			response(w, r, http.StatusBadRequest, err)
			return
		}
		cache.SetKey(usrProduct.UserId, usrProduct.ProductId)
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
