package middleware

import (
	"net/http"
	"songsapi/logger"
	"time"
)

type responseWrapper struct {
	http.ResponseWriter
	status      int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
  

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWrapper{ ResponseWriter: w, status: 200 }
		next.ServeHTTP(wrapper, r)
		logger.Info.Printf("[%s] %s - %d, %s %s\n",
			r.Method, r.URL.Path, wrapper.status, r.RemoteAddr, time.Since(start))
	})
}


func CORSMiddware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
