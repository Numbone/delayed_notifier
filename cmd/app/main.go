package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/delayed_notifier/internal/config"
	"github.com/delayed_notifier/internal/db"
	"github.com/delayed_notifier/internal/handler"
	"github.com/delayed_notifier/internal/repository"
	"github.com/delayed_notifier/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/delayed_notifier/internal/queue"
)

func main() {
	ctx ,cancel:= context.WithCancel(context.Background())
	defer cancel()


	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	conn, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе: %v", err)
	}
	defer conn.Close()

	mq,err:= queue.NewRabbitMQ(cfg.RabbitURL)
	if err != nil {
		log.Fatalf("can not connect rabbitmq")
	}
	defer mq.Close()

	queries := db.New(conn)

	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	notificationRepo := repository.NewNotificationRepository(queries)
	notificationSvc := service.NewNotificationService(notificationRepo)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)


	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)
	notificationHandler.RegisterRoutes(mux)


	srv := &http.Server{
		Addr:    ":8773",
		Handler: mux,
	}


	go func() {
		log.Println("Сервер запущен на :8773")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped gracefully")
}
