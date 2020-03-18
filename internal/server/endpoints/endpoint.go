package endpoints

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agubarev/tetest/internal/currency"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type contextKey int

type Handler func(e Endpoint, w http.ResponseWriter, r *http.Request) (result interface{}, code int, err error)

// sample context keys
const (
	keyUserID contextKey = iota
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Error      string      `json:"error,omitempty"`
	ExecTime   float64     `json:"exec_time"`
	Payload    interface{} `json:"payload,omitempty"`
}

// NOTE: usually I like to use that approach over canonical
// middleware function nesting, but either way is fine
type Endpoint struct {
	manager *currency.Manager
	handler Handler
}

func NewEndpoint(m *currency.Manager, h Handler) Endpoint {
	if m == nil {
		panic(currency.ErrNilManager)
	}

	return Endpoint{
		manager: m,
		handler: h,
	}
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: here you could handle authentication and possibly authorization
	/*
		u, err := e.core.UserManager().UserByID(
			r.Context(),
			r.Context().Value(keyUserID).(int64),
		)
	*/

	// ...

	l := e.manager.Logger().With(
		zap.String("uri", r.RequestURI),
	)

	// time mark just before the execution
	start := time.Now()

	// calling its respective handler and
	result, code, err := e.handler(e, w, r)

	if err != nil && code != http.StatusNotFound {
		l.Warn("endpoint handler returned with an error", zap.Error(errors.Cause(err)))
	}

	// handler must always return correct status code
	if code == 0 {
		l.Warn("endpoint handler returned with zero code; setting to 500")
		code = http.StatusInternalServerError
	}

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	// ... handle error or pass it by right into a response
	response, err := json.Marshal(Response{
		StatusCode: code,
		Error:      errMsg,
		ExecTime:   time.Since(start).Seconds(),
		Payload:    result,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
