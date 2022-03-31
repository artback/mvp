package handler

import (
	"database/sql"
	"fmt"
	"github.com/artback/mvp/pkg/api/handler/producthandler"
	"github.com/artback/mvp/pkg/api/handler/userhandler"
	"github.com/artback/mvp/pkg/api/handler/vendinghandler"
	"github.com/artback/mvp/pkg/api/middleware/logging"
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/api/middleware/security/basic"
	"github.com/artback/mvp/pkg/coin"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/artback/mvp/pkg/usecase"
	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"path/filepath"
)

func printWalk(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
	log.Printf("%s %s\n", method, route) // Walk and print out all routes

	return nil
}

func HttpRouter(db *sql.DB, coins coin.Coins) (chi.Router, error) {
	configPath, err := filepath.Abs("./config")
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(configPath+"/rbac_model.conf", configPath+"/auth_policy.csv")
	if err != nil {
		return nil, err
	}

	userService := usecase.UserService{Repository: postgres.UserRepository{DB: db}, Coins: coins}
	basicAuth := basic.Basic{Service: userService}

	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		logging.RequestLoggerMiddleware,
		middleware.Recoverer,
		security.Authenticate(basicAuth),
		security.Authorize(e),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			service := userService
			handler := userhandler.RestHandler{Service: service}
			r.Post("/", handler.CreateUser)
			r.Get("/{username}", handler.GetUser)
			r.Put("/", handler.UpdateUser)
			r.Delete("/", handler.DeleteUser)
		})
		r.Route("/product", func(r chi.Router) {
			service := usecase.ProductService{Repository: postgres.ProductRepository{DB: db}}
			handler := producthandler.RestHandler{Service: service}
			r.Get("/{product_name}", handler.GetProduct)
			r.Post("/", handler.CreateProduct)
			r.Put("/{product_name}", handler.UpdateProduct)
			r.Delete("/{product_name}", handler.DeleteProduct)
		})
		r.Route("/", func(r chi.Router) {
			service := usecase.VendingService{Repository: postgres.VendingRepository{DB: db}, Coins: coins}
			handler := vendinghandler.RestHandler{Service: service}
			r.Get("/deposit", handler.GetAccount)
			r.Put("/deposit", handler.Deposit)
			r.Post("/buy/{product_name}", handler.BuyProduct)
			r.Delete("/reset", handler.ResetDeposit)

		})
	})

	if err := chi.Walk(router, printWalk); err != nil {
		return router, fmt.Errorf("logging err: %v", err)
	}

	return router, nil
}
