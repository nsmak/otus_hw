package rest

import (
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			start := time.Now()
			s.log.Info(
				"[REST]",
				s.log.String("addr", r.RemoteAddr),
				s.log.String("method", r.Method),
				s.log.String("path", r.URL.Path),
				s.log.String("proto", r.Proto),
				s.log.Duration("duration", time.Since(start)),
				s.log.String("user agent", r.UserAgent()),
			)
		}()
		next.ServeHTTP(w, r)
	})
}
