package mid

import (
	"log"
	"net/http"
	"time"
)

// customResponseWriter это http.ResponseWriter, который запоминает статус ответа
type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func newCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{
		ResponseWriter: w,
		status:         200,
	}
}

// WriteHeader реализует http.ResponseWriter и сохраняет статус
func (c *customResponseWriter) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

// Logger это middleware для логирования информации о запросе
// формат: (200) GET /foo -> IP ADDR (задержка)
func Logger(log *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := newCustomResponseWriter(w)
			h.ServeHTTP(ww, r)
			defer func() {
				var params string
				if len(r.URL.Query()) > 0 {
					params = "?" + r.URL.Query().Encode()
				}
				log.Printf("(%d) : %s %s%s -> %s (%s)",
					ww.status,
					r.Method, r.URL.Path, params,
					r.RemoteAddr, time.Since(t1),
				)
			}()
		}
		return http.HandlerFunc(fn)
	}
}
