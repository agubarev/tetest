package endpoints

import (
	"net/http"
)

func CurrencyGetLatest(e Endpoint, w http.ResponseWriter, r *http.Request) (result interface{}, code int, err error) {
	// obtaining latest currency values
	result, err = e.manager.GetLatest(r.Context())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}
