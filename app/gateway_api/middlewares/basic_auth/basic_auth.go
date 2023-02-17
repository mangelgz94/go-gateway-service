package basic_auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
)

type BasicAuthMiddleware struct {
	config *BasicAuthConfig
}

func (m *BasicAuthMiddleware) HandleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(m.config.AuthUser))
			expectedPasswordHash := sha256.Sum256([]byte(m.config.AuthPassword))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		response, _ := json.Marshal(map[string]string{"error": "Unauthorized"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(response)

		return
	})
}

func New(config *BasicAuthConfig) *BasicAuthMiddleware {
	return &BasicAuthMiddleware{
		config: config,
	}
}
