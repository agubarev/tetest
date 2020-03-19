package endpoints_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agubarev/tetest/internal/currency"
	"github.com/agubarev/tetest/internal/server/endpoints"
	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// NOTE: since this is a simple test project, I'm not using
// any abstraction for the feed source for this test, only a
// simple in-memory store to fetch test values from

func init() {
	os.Setenv("FEED_URL", "https://www.bank.lv/vk/ecb_rss.xml")
}

func TestEndpointGetLatest(t *testing.T) {
	a := assert.New(t)

	// initializing currency manager
	m, err := currency.NewManager(currency.NewMemoryStore(), os.Getenv("FEED_URL"))
	a.NoError(err)
	a.NotNil(m)

	// importing real values from a real external endpoint
	a.NoError(m.Import(context.Background()))

	req, err := http.NewRequest("GET", "/api/v1/currency", nil)
	a.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// calling endpoint
	endpoints.NewEndpoint(m, endpoints.CurrencyGetLatest).ServeHTTP(rr, req)

	// unmarshaling response payload
	resp := endpoints.Response{}
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &resp))
	a.Equal(http.StatusOK, resp.StatusCode)
	a.Empty(resp.Error)
	a.NotEmpty(resp.Payload)
}

func TestEndpointGetByID(t *testing.T) {
	a := assert.New(t)

	// initializing currency manager
	m, err := currency.NewManager(currency.NewMemoryStore(), os.Getenv("FEED_URL"))
	a.NoError(err)
	a.NotNil(m)

	// importing real values from a real external endpoint
	a.NoError(m.Import(context.Background()))

	req, err := http.NewRequest("GET", "/api/v1/currency/USD", nil)
	a.NoError(err)

	// hacky way to fully enable chi router with this test request
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "USD")

	// injecting manager request's context
	req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, chiCtx))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// calling endpoint
	endpoints.NewEndpoint(m, endpoints.CurrencyGetByID).ServeHTTP(rr, req)

	// unmarshaling response payload
	resp := endpoints.Response{}
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &resp))
	a.Equal(http.StatusOK, resp.StatusCode)
	a.Empty(resp.Error)
	a.NotEmpty(resp.Payload)
}
