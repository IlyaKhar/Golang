package auth

import (
    "context"
    "net/http"
    "strings"
)

type ctxKey string
const userCtxKey ctxKey = "auth_user"

func WithAuth(m *Manager, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authz := r.Header.Get("Authorization")
        if !strings.HasPrefix(authz, "Bearer ") {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        raw := strings.TrimPrefix(authz, "Bearer ")
        claims, err := m.ParseAndValidate(raw)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), userCtxKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func FromContext(r *http.Request) (*Claims, bool) {
    v := r.Context().Value(userCtxKey)
    c, ok := v.(*Claims)
    return c, ok
}

func RequireRole(role string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims, ok := FromContext(r)
        if !ok || claims.Role != role {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}


