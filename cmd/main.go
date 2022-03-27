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
	"os/signal"
	"syscall"
	"time"
)

var (
	host  *string
	coins *[]int
)

func init() {
	host = flag.String("http-host", ":7070", "http host")
	coins = flag.IntSlice("coins", []int{5, 10, 20, 50, 100}, "coins")
}
func main() {
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

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	s := graceful.Server{
		Server: &http.Server{
			Addr:         *host,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.RegisterOnShutdown(func() {
		err := db.Close()
		if err != nil {
			log.Printf("Database closed with: %v", err)
		}
	})
	if err := s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server closed with: %v", err)
			os.Exit(1)
		}
		log.Printf("HTTP server shut down")
	}
}
