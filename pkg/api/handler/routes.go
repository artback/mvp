package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/artback/mvp/pkg/api/handler/producthandler"
	"github.com/artback/mvp/pkg/api/handler/userhandler"
	"github.com/artback/mvp/pkg/api/handler/vendinghandler"
	"github.com/artback/mvp/pkg/api/middleware/authentication/basic"
	"github.com/artback/mvp/pkg/api/middleware/logging"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/artback/mvp/pkg/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func printWalk(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
	log.Printf("%s %s\n", method, route) // Walk and print out all routes

	return nil
}

func HttpRouter(db *sql.DB, coins coin.Coins) (chi.Router, error) {
	userService := service.UserService{Repository: postgres.UserRepository{DB: db}, Coins: coins}
	auth := basic.Auth{Service: userService}

	productService := service.ProductService{Repository: postgres.ProductRepository{DB: db}}
	vendingService := service.VendingService{Repository: postgres.VendingRepository{DB: db}, Coins: coins}

	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		logging.RequestLoggerMiddleware,
		middleware.Recoverer,
	)
	router.Route("/v1", func(rc chi.Router) {
		rc.Mount("/user", userhandler.Routes(auth, userService))
		rc.Mount("/product", producthandler.Routes(auth, productService))
		rc.Mount("/", vendinghandler.Routes(auth, vendingService))
	})

	if err := chi.Walk(router, printWalk); err != nil {
		return router, fmt.Errorf("logging err: %v", err)
	}

	return router, nil
}
