package mid

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recoverer это middleware для восстановления после паник
func Recoverer(l *log.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					l.Printf("request panic, %v", rvr)
					l.Println(string(debug.Stack()))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
