package setup

import (
	"context"
	"net/http"
	"strings"

	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/token"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			response.SendWithError(w, http.StatusUnauthorized, *errors.AuthorizationError("Missing authorization header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		
		e := config.NewEnv()
		tCtx := token.New(e.JwtSecretKey)
		c, err := tCtx.Validate(tokenString)
		if err != nil {
			response.SendWithError(w, http.StatusUnauthorized, *errors.AuthorizationError(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), "claims", c.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}