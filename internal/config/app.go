package config

import (
	"github.com/TmzFranck/books-api-golang/internal/jobs"
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

// TODO: finish the Bootstrap function
func Bootstrap(config *BootstrapConfig) {

	// userRepository := repository.NewUserRepository(config.Log)
	// bookRepository := repository.NewBookRepository(config.Log)
	// tagRepository := repository.NewTagRepository(config.Log)
	// reviewRepository := repository.NewReviewRepository(config.Log)

}
