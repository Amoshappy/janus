package oauth

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// RevokeMiddleware prevents requests to an API from exceeding a specified rate limit.
type RevokeMiddleware struct {
	oauthServer *Spec
}

// NewRevokeMiddleware creates an instance of RevokeMiddleware
func NewRevokeMiddleware(oauthServer *Spec) *RevokeMiddleware {
	return &RevokeMiddleware{oauthServer}
}

// Handler is the middleware method.
func (m *RevokeMiddleware) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Starting RevokeMiddleware middleware")
		handler.ServeHTTP(w, r)

		if "" != r.Header.Get("Authorization") {
			log.Debug("Authorization header is empty")
			return
		}

		accessToken := r.Form.Get("access_token")
		if "" == accessToken {
			log.Warn("Token is empty or not set")
			return
		}

		log.Debug("Trying to remove the token")
		err := m.oauthServer.Manager.Remove(accessToken)
		if nil != err {
			log.WithError(err).Error("Not able to remove the token")
		}
	})
}
