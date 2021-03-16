package web

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Respond преобразует data в JSON и отправляет его клиенту
func Respond(w http.ResponseWriter, data interface{}, httpStatusCode int) {
	if httpStatusCode == http.StatusNoContent {
		w.WriteHeader(httpStatusCode)
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(res)
}

// RespondError отправляет ошибку клиенту
func RespondError(w http.ResponseWriter, httpStatusCode int, err error, errCode int, details ...string) {
	er := map[string]string{
		"errCode": strconv.Itoa(errCode),
		"error":   err.Error(),
	}
	if len(details) > 0 {
		er["details"] = details[0]
	}
	Respond(w, er, httpStatusCode)
}
