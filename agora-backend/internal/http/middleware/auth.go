package middleware

import (
	"net/http"
	"strings"

	"github.com/GabrielPacotte/Agora/internal/auth"
	httpcommon "github.com/GabrielPacotte/Agora/internal/http/common"
)

func RequireAuth(verifier auth.AccessTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" {
				httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "missing authorization header")
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authz, prefix) {
				httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "invalid authorization scheme")
				return
			}

			raw := strings.TrimSpace(strings.TrimPrefix(authz, prefix))
			if raw == "" {
				httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "missing token")
				return
			}

			userID, err := verifier.Verify(raw)
			if err != nil {
				httpcommon.WriteError(w, http.StatusUnauthorized, httpcommon.ErrCodeUnauthorized, "invalid token")
				return
			}

			ctx := httpcommon.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
