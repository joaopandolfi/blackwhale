package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaopandolfi/blackwhale/models/auth"
	"github.com/joaopandolfi/blackwhale/utils"
)

const (
	HEADER_USERID      = "_xid"
	HEADER_INSTITUTION = "_xinstitution"
	HEADER_PERMISSION  = "_xpermission"

	invalidPermissionMessage = "Not authorized"
)

// TokenHandler -
// @handler
// Intercept all transactions and check if is authenticated by token
func TokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		token := GetHeader(r, "token")
		t, err := utils.CheckJwtToken(token)

		if !t.Authorized || err != nil {
			utils.Debug("[TokenHandler]", "Auth Error", url)
			Response(w, invalidPermissionMessage, http.StatusForbidden)
			return
		}

		InjectHeader(r, HEADER_PERMISSION, t.Permission)
		InjectHeader(r, HEADER_INSTITUTION, t.Institution)
		InjectHeader(r, HEADER_USERID, t.ID)

		utils.Debug("[TokenHandler]", "Authenticated", url)
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

// HandlerTokenPermissions -
// check if the request contain the permissions
func HandleTokenPermissions(r *mux.Router, path string, f http.HandlerFunc, permissions []string, methods ...string) {
	r.HandleFunc(path, Chain(f, PermissionMiddleware(permissions), TokenHandler)).Methods(methods...)
}
