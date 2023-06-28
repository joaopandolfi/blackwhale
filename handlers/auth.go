package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/models/auth"
	"github.com/joaopandolfi/blackwhale/remotes/jwt"
	"github.com/joaopandolfi/blackwhale/utils"
)

const (
	HEADER_USERID      = "_xid"
	HEADER_INSTITUTION = "_xinstitution"
	HEADER_PERMISSION  = "_xpermission"
	HEADER_BROKER      = "_xbroker"

	invalidPermissionMessage = "Not authorized"
	quietCtx                 = "quiet"
)

// TokenHandler -
// @handler
// Intercept all transactions and check if is authenticated by token
func TokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		token := GetHeader(r, "token")
		t, err := jwt.CheckJwtToken(token, configurations.Configuration.Security.JWTSecret)

		if !t.Authorized || err != nil {
			utils.Debug("[TokenHandler]", "Auth Error", url, err.Error())
			Response(w, invalidPermissionMessage, http.StatusForbidden)
			return
		}

		broker, err := json.Marshal(t.Broker)
		if err == nil {
			InjectHeader(r, HEADER_BROKER, string(broker))
		}

		InjectHeader(r, HEADER_PERMISSION, t.Permission)
		InjectHeader(r, HEADER_INSTITUTION, t.Institution)
		InjectHeader(r, HEADER_USERID, t.ID)

		ctx := r.Context()
		quiet := ctx.Value(quietCtx)
		if quiet == nil {
			utils.Debug("[TokenHandler]", "Authenticated", url)
		}

		next.ServeHTTP(w, r)
	})
}

// PermissionMiddleware -
// @middleware
// Intercept all transactions and check if user has perission
func PermissionMiddleware(expectedPermissions []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := GetHeader(r, HEADER_PERMISSION)

			for _, permission := range expectedPermissions {
				if auth.PermissionContain(permissions, permission) {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.CriticalError("[Permission][PermissionMiddleware]", "Permission Denied", r.URL.String(), expectedPermissions, permissions)
			Response(w, invalidPermissionMessage, http.StatusForbidden)
		})
	}
}

// QuietMiddleware -
// @middleware
// Intercept the request and add quiet flag
func QuietMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, quietCtx, true)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// HandlerTokenPermissions -
// check if the request contain the permissions
func HandleTokenPermissions(r *mux.Router, path string, f http.HandlerFunc, permissions []string, methods ...string) {
	r.HandleFunc(path, Chain(f, PermissionMiddleware(permissions), TokenHandler)).Methods(methods...)
}

// QuietHandlerTokenPermissions -
// Same as HandlerTokenPermissions but without log url
func QuietHandleTokenPermissions(r *mux.Router, path string, f http.HandlerFunc, permissions []string, methods ...string) {
	r.HandleFunc(path, Chain(f, PermissionMiddleware(permissions), TokenHandler, QuietMiddleware)).Methods(methods...)
}
