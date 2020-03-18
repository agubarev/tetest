package endpoints

import (
	"net/http"

	"github.com/agubarev/tetest/internal/currency"
	"github.com/go-chi/chi"
)

func CurrencyGetByID(e Endpoint, w http.ResponseWriter, r *http.Request) (result interface{}, code int, err error) {
	// obtaining currency history by ID
	switch result, err = e.manager.GetByID(r.Context(), chi.URLParam(r, "id")); err {
	case nil: // all good
		return result, http.StatusOK, nil
	case currency.ErrCurrencyNotFound: // handling 404
		return nil, http.StatusNotFound, err
	default: // regular error
		return nil, http.StatusInternalServerError, err
	}
}
