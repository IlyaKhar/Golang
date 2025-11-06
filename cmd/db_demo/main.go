package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"day2/auth"
	"day2/config"
	"day2/db"
	"day2/internal/structure/domain"
	transporthttp "day2/internal/transport/http"
	"day2/logging"
	middleware "day2/middleware"
	"day2/monitoring"

	_ "modernc.org/sqlite"
)

// простая in-memory реализация UserService для запуска демо
// удовлетворяет интерфейсу, который ожидают хендлеры
type memUserService struct {
	users map[int]domain.User
}

func newMemUserService() *memUserService { return &memUserService{users: make(map[int]domain.User)} }

func (s *memUserService) CreateUser(rctx context.Context, u *domain.User) error {
	if u == nil || u.Name == "" || u.Email == "" || u.Age < 0 || u.Age > 150 {
		return domain.ErrInvalidInput
	}
	if _, ok := s.users[u.ID]; ok {
		return domain.ErrEmailExists
	}
	s.users[u.ID] = *u
	return nil
}
func (s *memUserService) GetUser(rctx context.Context, id int) (domain.User, error) {
	if v, ok := s.users[id]; ok {
		return v, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}
func (s *memUserService) ListUsers(rctx context.Context, limit, offset int) ([]domain.User, error) {
	out := make([]domain.User, 0, len(s.users))
	for _, v := range s.users {
		out = append(out, v)
	}
	return out, nil
}
func (s *memUserService) UpdateUser(rctx context.Context, u *domain.User) error {
	if u == nil || u.ID <= 0 {
		return domain.ErrInvalidInput
	}
	if _, ok := s.users[u.ID]; !ok {
		return domain.ErrUserNotFound
	}
	s.users[u.ID] = *u
	return nil
}
func (s *memUserService) DeleteUser(rctx context.Context, id int) error {
	if _, ok := s.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config:", err)
	}

	// 1) SQLite + миграции (SQL-файлы в коде)
	sqlDB, err := sql.Open("sqlite", cfg.SQLiteDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := db.ApplyMigrations(sqlDB); err != nil {
		log.Fatal("migrations:", err)
	}

	// 2) Service + Handlers (для простоты — in-memory; БД уже инициализирована для будущих шагов)
	userSvc := newMemUserService()
	userHandlers := transporthttp.NewUserHandlers(userSvc)

	// 3) JWT manager
	mgr := auth.NewManager(auth.Config{
		Secret:    []byte(cfg.JWTSecret),
		AccessTTL: 15 * time.Minute,
		Issuer:    cfg.JWTIssuer,
		Audience:  cfg.JWTAudience,
	})

	// 4) Router
	mux := http.NewServeMux()
	transporthttp.RegisterRoutes(mux, userHandlers, mgr)

	// 5) Logging + Recover middleware (zap)
	logg, err := logging.NewZap(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	handler := middleware.Recover(middleware.LoggingZap(logg)(mux))
	handler = monitoring.WithPrometheus(handler)

	mux.Handle("/metrics", monitoring.MetricsHandler())

	// pprof на cfg.PprofPort (в отдельной горутине)
	go func() {
		m := http.NewServeMux()
		m.HandleFunc("/debug/pprof/", pprof.Index)
		m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("/debug/pprof/profile", pprof.Profile)
		m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("/debug/pprof/trace", pprof.Trace)

		log.Printf("pprof on :%d\n", cfg.PprofPort)
		if err := http.ListenAndServe(
			":"+func(p int) string { return fmt.Sprintf("%d", p) }(cfg.PprofPort), m,
		); err != nil {
			log.Println("pprof server:", err)
		}
	}()

	// 6) HTTP server
	srv := &http.Server{
		Addr:    ":" + func(p int) string { return fmt.Sprintf("%d", p) }(cfg.Port),
		Handler: handler,
	}
	log.Fatal(srv.ListenAndServe())
}
