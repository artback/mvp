package main

import (
	"database/sql"
	"github.com/artback/mvp/pkg/config"
	"github.com/artback/mvp/pkg/handler"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
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

	router := handler.HttpRouter(db, *coins)

	s := http.Server{
		Addr:         *host,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
