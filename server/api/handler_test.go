package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Run("Positive router CRUD", func(t *testing.T) {
		r := NewRouter()
		require.NotNil(t, r)
	})
}

func TestController_GetPair(t *testing.T) {
	t.Run("Positive get", func(t *testing.T) {
		r := NewRouter()
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/rates?pairs=BTC-USDT,ETH-USDT", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})
	t.Run("Negative with another method", func(t *testing.T) {
		r := NewRouter()
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPatch, ts.URL+"/api/v1/rates?pairs=BTC-USDT,ETH-USDT", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "method does not allowed", string(body))
	})
}

func TestController_PostPair(t *testing.T) {
	t.Run("Positive post", func(t *testing.T) {
		r := NewRouter()
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/rates", bytes.NewBuffer([]byte(`{"pairs":["BTC-USDT","ETH-USDT"]}`)))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("Negative post with another route", func(t *testing.T) {
		r := NewRouter()
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/test", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "route does not exist", string(body))
	})
}
