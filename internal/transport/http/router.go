package transporthttp

import (
    "net/http"
    "day2/auth"
)

func RegisterRoutes(mux *http.ServeMux, h *UserHandlers, mgr *auth.Manager) {
    // /login (public)
    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        h.Login(w, r, mgr)
    })

    // защищённые маршруты
    mux.Handle("/users", auth.WithAuth(mgr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.ListUsers(w, r)
        case http.MethodPost:
            auth.RequireRole("admin", http.HandlerFunc(h.CreateUser)).ServeHTTP(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })))

    mux.Handle("/users/", auth.WithAuth(mgr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.GetUser(w, r)
        case http.MethodPut:
            auth.RequireRole("admin", http.HandlerFunc(h.UpdateUser)).ServeHTTP(w, r)
        case http.MethodDelete:
            auth.RequireRole("admin", http.HandlerFunc(h.DeleteUser)).ServeHTTP(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })))
}


