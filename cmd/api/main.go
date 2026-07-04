package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"svitlonav/internal/config"
	"svitlonav/internal/handler"
	"svitlonav/internal/repository"
	"svitlonav/internal/service"
	"svitlonav/internal/valhalla"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := pgxpool.New(connectCtx, cfg.DatabaseURL)
	cancelConnect()
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	lampRepo := repository.NewLampRepository(pool)
	valhallaClient := valhalla.NewClient(cfg.ValhallaURL)

	routeService := service.NewRouteService(lampRepo, valhallaClient)
	routeHandler := handler.NewRouteHandler(routeService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/route", routeHandler.GetRoute)
	mux.HandleFunc("GET /api/health", routeHandler.HealthCheck)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		http.NotFound(w, r)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("svitlonav listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
