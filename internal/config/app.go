package config

import (
	"github.com/TmzFranck/books-api-golang/internal/delivery/http"
	"github.com/TmzFranck/books-api-golang/internal/delivery/http/middleware"
	"github.com/TmzFranck/books-api-golang/internal/delivery/http/route"
	"github.com/TmzFranck/books-api-golang/internal/jobs"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/TmzFranck/books-api-golang/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *chi.Mux
	Log         *logrus.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	WorkerPool  *jobs.WokerPool
	RedisClient *redis.Client
}

func Bootstrap(config *BootstrapConfig) {

	userRepository := repository.NewUserRepository(config.Log)
	bookRepository := repository.NewBookRepository(config.Log)
	tagRepository := repository.NewTagRepository(config.Log)
	reviewRepository := repository.NewReviewRepository(config.Log)

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Config, config.WorkerPool, userRepository, config.RedisClient)
	bookUseCase := usecase.NewBookUseCase(config.DB, config.Log, bookRepository)
	tagUseCase := usecase.NewTagUseCase(config.DB, config.Log, tagRepository)
	reviewUseCase := usecase.NewReviewUseCase(config.DB, config.Log, reviewRepository)

	userController := http.NewUserController(config.Log, userUseCase, config.Validate)
	bookController := http.NewBookController(config.Log, bookUseCase, config.Validate)
	tagController := http.NewTagController(config.Log, tagUseCase, config.Validate)
	reviewController := http.NewReviewController(config.Log, reviewUseCase, config.Validate)

	authMiddleware := middleware.NewAuthMiddleware(config.RedisClient, config.Log)

	routeConfig := route.RouteConfig{
		App:              config.App,
		UserController:   userController,
		BookController:   bookController,
		TagController:    tagController,
		ReviewController: reviewController,
		RedisClient:      config.RedisClient,
		Logger:           config.Log,
		AuthMiddleware:   authMiddleware,
	}

	routeConfig.Setup()
}
