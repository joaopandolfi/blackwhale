package handlers

import (
	"context"
	"net/http"
)

const (
	OperatorID   ContextInjection = ":ci:operatorID"
	OperatorRole ContextInjection = ":ci:operatorRole"
)

// InjectOperatorOnContext get user_id extracted from jwt and injected on headers
func InjectOperatorOnContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		injectedAuth := ctx.Value(InjectedAuth)
		if injectedAuth != nil {
			operatorID := r.Header.Get(HEADER_USERID)
			if operatorID != "" {
				ctx := r.Context()
				ctx = context.WithValue(ctx, OperatorID, operatorID)
				r = r.WithContext(ctx)
			}

			role := r.Header.Get(HEADER_PERMISSION)
			if role != "" {
				ctx := r.Context()
				ctx = context.WithValue(ctx, OperatorRole, role)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ExtractValueFromContext - grab injected data on context
// returns "" if does not exists
func ExtractValueFromContext(ctx context.Context, v ContextInjection) string {
	converted := ""
	raw := ctx.Value(v)
	if raw != nil {
		value, ok := raw.(string)
		if ok {
			converted = value
		}
	}
	return converted
}
