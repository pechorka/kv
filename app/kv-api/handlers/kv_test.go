package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/pechorka/kv/internal/store"
	"github.com/stretchr/testify/require"
)

func TestKV(t *testing.T) {
	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := KVTests{api: API(log, store.NewMemory())}

	t.Run("List", tests.List)
	t.Run("CRUD", tests.CRUD)
}

type KVTests struct {
	api http.Handler
}

func (kvt *KVTests) List(t *testing.T) {
	testData := map[string]interface{}{
		"int":   1,
		"float": 1.0,
		"bool":  true,
	}

	// заполняем хранилище данными
	{
		const setReqTemplate = `{"key":"%s", "value":%v}`
		for key, val := range testData {
			body := fmt.Sprintf(setReqTemplate, key, val)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/kv", strings.NewReader(body))
			resp := httptest.NewRecorder()
			kvt.api.ServeHTTP(resp, req)
			require.Equal(t, http.StatusNoContent, resp.Code)
		}

		// string отдельно добавляем, он не вписывается в шаблон
		testData["string"] = "string"
		body := fmt.Sprintf(`{"key":"%s", "value":"%s"}`, "string", "string")
		req := httptest.NewRequest(http.MethodPost, "/api/v1/kv", strings.NewReader(body))
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNoContent, resp.Code)
	}

	// проверяем что все правильно сохранилось
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kv", nil)
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		storeContent := make(map[string]interface{})
		err := json.NewDecoder(resp.Body).Decode(&storeContent)
		require.NoError(t, err)

		for key, val := range testData {
			// EqualValues, а не Equals, потому что go декодит int как float64
			require.EqualValues(t, val, storeContent[key])
		}
	}
}

func (kvt *KVTests) CRUD(t *testing.T) {
	const (
		testKey   = "testKey"
		testValue = "testValue"
	)

	// SET
	{
		body := fmt.Sprintf(`{"key":"%s", "value":"%s"}`, testKey, testValue)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/kv", strings.NewReader(body))
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNoContent, resp.Code)
	}

	// GET
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kv/"+testKey, nil)
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		respBody := make(map[string]interface{})
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		require.Equal(t, testValue, respBody["value"])
	}

	// UPDATE
	{
		const testNewValue = "some new value"
		body := fmt.Sprintf(`{"key":"%s", "value":"%s"}`, testKey, testNewValue)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/kv", strings.NewReader(body))
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNoContent, resp.Code)

		req = httptest.NewRequest(http.MethodGet, "/api/v1/kv/"+testKey, nil)
		resp = httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		respBody := make(map[string]interface{})
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		require.Equal(t, testNewValue, respBody["value"])
	}

	// DELETE
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/kv/"+testKey, nil)
		resp := httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNoContent, resp.Code)

		req = httptest.NewRequest(http.MethodGet, "/api/v1/kv/"+testKey, nil)
		resp = httptest.NewRecorder()
		kvt.api.ServeHTTP(resp, req)
		require.Equal(t, http.StatusNotFound, resp.Code)
	}
}
