package encore

import (
	"net/http"

	"encore.dev/rlog"
)

// Fallback is called for the debug enpoints.
//
//encore:api public raw path=/!fallback
func (s *Service) Fallback(w http.ResponseWriter, req *http.Request) {
	rlog.Info("FALLBACK", "url", req.URL.String())
	s.debug.ServeHTTP(w, req)
}
