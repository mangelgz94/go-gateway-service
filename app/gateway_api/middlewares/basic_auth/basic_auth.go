package basic_auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/mangelgz94/thinksurance-miguel-angel-gonzalez-morera/app/gateway_api/middlewares"
)

type BasicAuthMiddleware struct {
	config *BasicAuthConfig
}

func (m *BasicAuthMiddleware) HandleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isUserAuthorized bool
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(m.config.AuthUser))
			expectedPasswordHash := sha256.Sum256([]byte(m.config.AuthPassword))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				isUserAuthorized = true
			}
		}

		requestContext := context.WithValue(r.Context(), middlewares.IsUserAuthorized, isUserAuthorized)
		r = r.WithContext(requestContext)
		next.ServeHTTP(w, r)
	})
}

func New(config *BasicAuthConfig) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		config: config,
	}
}
