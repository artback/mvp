package main

import (
	"database/sql"
	"errors"
	"github.com/artback/mvp/internal/config"
	"github.com/artback/mvp/pkg/api/graceful"
	"github.com/artback/mvp/pkg/api/handler"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"time"
)

const timeout = 5 * time.Second

func main() {
	host := flag.String("http-host", ":7070", "http host")
	coins := flag.IntSlice("coins", []int{5, 10, 20, 50, 100}, "coins")
	flag.Parse()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", c.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	router, err := handler.HttpRouter(db, *coins)
	if err != nil {
		log.Fatal(err)
	}

	server := graceful.Server{
		Server: &http.Server{
			Addr:         *host,
			Handler:      router,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}
	server.RegisterOnShutdown(func() {
		err := db.Close()
		if err != nil {
			log.Printf("Database closed with: %v", err)
		}
	})

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server closed with: %v", err)
			os.Exit(1)
		}

		log.Print("HTTP server shut down")
	}
}
