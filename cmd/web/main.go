package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TmzFranck/books-api-golang/internal/config"
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/jobs"
	"github.com/TmzFranck/books-api-golang/internal/utils"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator()
	app := config.NewRouter(log)
	redisClient := config.NewRedisClient(viperConfig, viperConfig.GetString("redis.address"))
	workerPool := jobs.NewWorkerPool(5, 10, log)

	workerPool.RegisterHandler("SendMail", utils.Send)
	workerPool.Start()

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Book{},
		&entity.Tag{},
		&entity.Review{},
	); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		WorkerPool:  workerPool,
		RedisClient: redisClient,
	})

	port := viperConfig.GetInt("server.port")
	addr := ":" + strconv.Itoa(port)

	server := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("HTTP server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Infof("received signal: %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("graceful shutdown failed: %v", err)
		if closeErr := server.Close(); closeErr != nil {
			log.Errorf("forced close failed: %v", closeErr)
		}
	}

	if err := workerPool.Shutdown(10 * time.Second); err != nil {
		log.Errorf("worker pool stop failed: %v", err)
	}

	log.Info("server stopped cleanly")
}
