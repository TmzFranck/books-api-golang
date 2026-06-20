package route

import (
	"net/http"

	controller "github.com/TmzFranck/books-api-golang/internal/delivery/http"
	"github.com/TmzFranck/books-api-golang/internal/delivery/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	App              *chi.Mux
	UserController   *controller.UserController
	BookController   *controller.BookController
	ReviewController *controller.ReviewController
	TagController    *controller.TagController
	RedisClient      *redis.Client
	Logger           *logrus.Logger
	AuthMiddleware   middleware.HandlerFunc
}

func (c *RouteConfig) Setup() {
	c.RegisterGuestsRoute()
	c.SetupAuthRoute()
	c.Health()
}

func (c *RouteConfig) RegisterGuestsRoute() {
	c.App.Post("/api/v1/login", c.UserController.Login)
	c.App.Post("/api/v1/register", c.UserController.Register)
	c.App.Post("/api/v1/auth/verify/{token}", c.UserController.VerifyUserAccount)
	c.App.Post("/api/v1/auth/password-reset-confirm/{token}", c.UserController.ConfirmPasswordReset)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Group(func(r chi.Router) {
		r.Use(c.AuthMiddleware)

		// Authentication
		r.Post("/api/v1/auth/logout", c.UserController.Logout)
		r.Post("/api/v1/auth/refresh-token", c.UserController.RefreshToken)
		r.Get("/api/v1/auth/me", c.UserController.GetCurrentUser)
		r.Post("/api/v1/auth/password-reset-request", c.UserController.PasswordReset)
		r.Post("/api/v1/send-mail", c.UserController.SendMail)

		// Books
		r.Get("/api/v1/books", c.BookController.GetBooks)
		r.Get("/api/v1/books/{book_id}", c.BookController.GetBook)
		r.Post("/api/v1/books", c.BookController.CreateBook)
		r.Put("/api/v1/books/{book_id}", c.BookController.UpdateBook)
		r.Delete("/api/v1/books/{book_id}", c.BookController.DeleteBook)
		r.Get("/api/v1/books/user", c.BookController.GetUserBooks)

		// Reviews
		r.Get("/api/v1/reviews", c.ReviewController.GetReviews)
		r.Get("/api/v1/reviews/{review_id}", c.ReviewController.GetReview)
		r.Delete("/api/v1/reviews/{review_id}", c.ReviewController.DeleteReviewFromBook)
		r.Post("/api/v1/reviews/book/{book_id}", c.ReviewController.AddReviewToBook)

		// Tags
		r.Get("/api/v1/tags", c.TagController.GetTags)
		r.Post("/api/v1/tags", c.TagController.CreateTag)
		r.Put("/api/v1/tags/{tag_id}", c.TagController.UpdateTag)
		r.Delete("/api/v1/tags/{tag_id}", c.TagController.DeleteTag)
		r.Post("/api/v1/tags/books/{book_id}", c.TagController.AddTagToBook)
	})
}

func (c *RouteConfig) Health() {

	c.App.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
