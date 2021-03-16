package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pechorka/kv/internal/web"

	"github.com/pechorka/kv/internal/store"
)

type KVGroup struct {
	Store store.Store
}

func (kv KVGroup) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	val, ok, err := kv.Store.Get(key)
	if err != nil {
		web.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			web.ErrInternal,
			"ошибка при получении данных из хранилища",
		)
		return
	}

	if !ok {
		web.RespondError(
			w,
			http.StatusNotFound,
			errors.New("значения для ключа"+key+" не найдено"),
			web.ErrNotFound,
		)
		return
	}

	resp := map[string]interface{}{"value": val}

	web.Respond(w, resp, http.StatusOK)
}

func (kv KVGroup) Set(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}

	if err := web.DecodeJSON(r, &req); err != nil {
		web.RespondError(
			w,
			http.StatusBadRequest,
			err,
			web.ErrDecode,
		)
		return
	}

	if err := kv.Store.Set(req.Key, req.Value); err != nil {
		web.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			web.ErrInternal,
			"ошибка при записи в хранилище",
		)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}

func (kv KVGroup) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	if err := kv.Store.Delete(key); err != nil {
		web.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			web.ErrInternal,
			"ошибка при удалении данных из хранилища",
		)
		return
	}

	web.Respond(w, nil, http.StatusNoContent)
}

func (kv KVGroup) List(w http.ResponseWriter, r *http.Request) {
	storeContent, err := kv.Store.List()

	if err != nil {
		web.RespondError(
			w,
			http.StatusInternalServerError,
			err,
			web.ErrInternal,
			"ошибка при получении всех данных из хранилища",
		)
		return
	}

	web.Respond(w, storeContent, http.StatusOK)
}
