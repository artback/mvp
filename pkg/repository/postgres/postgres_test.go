//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"database/sql"
	"github.com/artback/mvp/pkg/products"
	"github.com/artback/mvp/pkg/repository/postgres"
	"github.com/artback/mvp/pkg/users"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var (
	db             *sql.DB
	defaultBuyer   = users.User{Username: "defaultBuyer", Password: "pass", Role: users.Buyer}
	defaultSeller  = users.User{Username: "defaultSeller", Password: "pass", Role: users.Seller}
	defaultProduct = products.Product{Name: "claratin", SellerId: defaultSeller.Username, Price: 5, Amount: 100}
)

// TestMain run before and after all tests in the package
// The state in the db is a shared global state between the tests and this a big code smell ,
//which could make tests dependent on the order and each other.
//The alternative could be to spin up a new container for each test but this would quickly add up and take a large amount of resources and time
func TestMain(m *testing.M) {
	ctx := context.Background()
	teardown := startPG()
	defer func() {
		err := teardown()
		if err != nil {
			log.Fatal(err)
		}
	}()
	userRepo := postgres.UserRepository{DB: db}
	err := userRepo.Insert(ctx, defaultSeller)
	if err != nil {
		log.Fatal("TestMain: ", err)
	}
	err = userRepo.Insert(ctx, defaultBuyer)
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func startPG() func() error {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("myuser", "mypass"),
		Path:   "mydatabase",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	pw, _ := pgURL.User.Password()
	env := []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + pgURL.Path,
	}
	abs, err := filepath.Abs("../../../db")
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: "postgres", Tag: "14-alpine", Env: env}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.Mounts = []docker.HostMount{
			{
				Target: "/docker-entrypoint-initdb.d",
				Source: abs,
				Type:   "bind",
			},
		}
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start postgres container: %v", err)
	}

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() (err error) {
		db, err = sql.Open("postgres", pgURL.String())
		if err != nil {
			return err
		}
		return db.Ping()
	})
	return func() error {
		return pool.Purge(resource)
	}
}
