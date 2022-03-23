package handler

import (
	"database/sql"
	"github.com/artback/mvp/pkg/authentication/basic"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/handler/producthandler"
	"github.com/artback/mvp/pkg/handler/userhandler"
	"github.com/artback/mvp/pkg/handler/vendinghandler"
	"github.com/artback/mvp/pkg/logging"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func printWalk(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
	log.Printf("%s %s\n", method, route) // Walk and print out all routes
	return nil
}

func HttpRouter(db *sql.DB, coins coin.Coins) chi.Router {
	productRepository := postgres.ProductRepository{DB: db}
	userRepository := postgres.UserRepository{DB: db, Coins: coins}
	vendingRepository := postgres.VendingRepository{DB: db, Coins: coins}
	auth := basic.Auth{Repository: userRepository}

	r := chi.NewRouter()
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		logging.RequestLoggerMiddleware,
		middleware.Recoverer,
	)
	r.Route("/v1", func(rc chi.Router) {
		r.Mount("/user", userhandler.Routes(auth, userRepository))
		rc.Mount("/product", producthandler.Routes(auth, productRepository))
		rc.Mount("/", vendinghandler.Routes(auth, vendingRepository))
	})

	if err := chi.Walk(r, printWalk); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	return r
}
