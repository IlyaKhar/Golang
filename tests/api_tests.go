package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"database/sql"

	"day2/auth"
	"day2/db"
	domain "day2/internal/structure/domain"
	transporthttp "day2/internal/transport/http"

	_ "modernc.org/sqlite"
)

// простая in-memory реализация сервиса для тестов (удовлетворяет интерфейсу из хендлеров)
type testMemService struct{ users map[int]domain.User }

func newTestMemService() *testMemService { return &testMemService{users: make(map[int]domain.User)} }

func (s *testMemService) CreateUser(ctx context.Context, u *domain.User) error {
	if u == nil {
		return domain.ErrInvalidInput
	}
	if _, ok := s.users[u.ID]; ok {
		return domain.ErrEmailExists
	}
	s.users[u.ID] = *u
	return nil
}

func (s *testMemService) GetUser(ctx context.Context, id int) (domain.User, error) {
	if v, ok := s.users[id]; ok {
		return v, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (s *testMemService) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	out := make([]domain.User, 0, len(s.users))
	for _, v := range s.users {
		out = append(out, v)
	}
	return out, nil
}

func (s *testMemService) UpdateUser(ctx context.Context, u *domain.User) error {
	if u == nil || u.ID <= 0 {
		return domain.ErrInvalidInput
	}
	if _, ok := s.users[u.ID]; !ok {
		return domain.ErrUserNotFound
	}
	s.users[u.ID] = *u
	return nil
}

func (s *testMemService) DeleteUser(ctx context.Context, id int) error {
	if _, ok := s.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	// изолированная БД в памяти
	sqlDB, err := sql.Open("sqlite", "file:test.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.ApplyMigrations(sqlDB); err != nil {
		t.Fatal(err)
	}

	// DI: упрощённый in-memory сервис для тестов
	h := transporthttp.NewUserHandlers(newTestMemService())

	// JWT
	mgr := auth.NewManager(auth.Config{
		Secret:    []byte("test-secret"),
		AccessTTL: time.Hour,
		Issuer:    "day2",
		Audience:  "day2-clients",
	})

	// router
	mux := http.NewServeMux()
	transporthttp.RegisterRoutes(mux, h, mgr)

	// без внешних мидлварей — напрямую mux
	return httptest.NewServer(mux)
}

func TestLogin_OK(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	reqBody := `{"email":"a@b.com","password":"123"}`
	res, err := http.Post(ts.URL+"/login", "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.AccessToken == "" {
		t.Fatalf("empty access token")
	}
}

func TestUsers_Unauthorized(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", res.StatusCode)
	}
}

func getToken(t *testing.T, ts *httptest.Server) string {
	t.Helper()
	// чтобы получить admin-токен — временно сделай роль в Login = role (или заведи отдельный login handler для админа)

	body := `{"email":"a@b.com","password":"123"}`
	res, err := http.Post(ts.URL+"/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("login status %d", res.StatusCode)
	}

	var out struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.NewDecoder(res.Body).Decode(&out)
	return out.AccessToken
}

func TestUsers_List_WithToken(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	token := getToken(t, ts)
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}
}

func TestUsers_Create_ForbiddenForUser(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	token := getToken(t, ts)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/users", strings.NewReader(`{"id":1,"name":"Alice","email":"alice@example.com","age":30,"is_active":true,"created_at":"2025-01-01T00:00:00Z"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status: got %d, want 403", res.StatusCode)
	}
}
