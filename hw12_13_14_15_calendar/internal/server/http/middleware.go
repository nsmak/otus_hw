package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			start := time.Now()
			resp := fmt.Sprintf("%s %s %s %s %d %v %s", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, 200, time.Since(start), r.UserAgent())
			s.log.Info(resp)
		}()
		next.ServeHTTP(w, r)
	})
}
