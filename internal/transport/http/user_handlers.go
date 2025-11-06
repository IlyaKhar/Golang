package transporthttp

import (
    "encoding/json"
    "context"
    "net/http"
    "strconv"
    "strings"

    "day2/auth"
    "day2/internal/structure/domain"
)

type UserService interface {
    CreateUser(rctx context.Context, u *domain.User) error
    GetUser(rctx context.Context, id int) (domain.User, error)
    ListUsers(rctx context.Context, limit, offset int) ([]domain.User, error)
    UpdateUser(rctx context.Context, u *domain.User) error
    DeleteUser(rctx context.Context, id int) error
}

type UserHandlers struct { svc UserService }

func NewUserHandlers(s UserService) *UserHandlers { return &UserHandlers{svc: s} }

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
type LoginResponse struct { AccessToken string `json:"access_token"` }

func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request, jwtMgr *auth.Manager) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var in LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    if in.Email == "" || in.Password == "" {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    // роль по умолчанию user (для админ‑ручек можешь временно поменять на "admin")
    token, err := jwtMgr.GenerateAccessToken(1, in.Email, "user")
    if err != nil { http.Error(w, "cannot issue token", http.StatusInternalServerError); return }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(LoginResponse{AccessToken: token})
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
    var u domain.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    if err := h.svc.CreateUser(r.Context(), &u); err != nil {
        switch err {
        case domain.ErrInvalidInput:
            http.Error(w, err.Error(), http.StatusBadRequest)
        case domain.ErrEmailExists:
            http.Error(w, err.Error(), http.StatusConflict)
        default:
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(u)
}

func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil { http.Error(w, "bad id", http.StatusBadRequest); return }
    u, err := h.svc.GetUser(r.Context(), id)
    if err != nil {
        if err == domain.ErrUserNotFound { http.Error(w, err.Error(), http.StatusNotFound); return }
        http.Error(w, err.Error(), http.StatusInternalServerError); return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(u)
}

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil { http.Error(w, "bad id", http.StatusBadRequest); return }
    var u domain.User
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    u.ID = id
    if err := h.svc.UpdateUser(r.Context(), &u); err != nil {
        switch err {
        case domain.ErrInvalidInput:
            http.Error(w, err.Error(), http.StatusBadRequest)
        case domain.ErrUserNotFound:
            http.Error(w, err.Error(), http.StatusNotFound)
        case domain.ErrEmailExists:
            http.Error(w, err.Error(), http.StatusConflict)
        default:
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(u)
}

func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil { http.Error(w, "bad id", http.StatusBadRequest); return }
    if err := h.svc.DeleteUser(r.Context(), id); err != nil {
        if err == domain.ErrUserNotFound { http.Error(w, err.Error(), http.StatusNotFound); return }
        if err == domain.ErrInvalidInput { http.Error(w, err.Error(), http.StatusBadRequest); return }
        http.Error(w, err.Error(), http.StatusInternalServerError); return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandlers) ListUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.svc.ListUsers(r.Context(), 50, 0)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(users)
}


