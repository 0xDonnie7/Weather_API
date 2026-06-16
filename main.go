package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type config struct {
	port      int
	env       string
	redisAddr string
	APIKey    string
}

type application struct {
	config config
	rdb    *redis.Client
	ctx    context.Context
}

func main() {
	var cfg config
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	flag.StringVar(&cfg.env, "env", "development", "development|staging|production")
	flag.IntVar(&cfg.port, "port", 8080, "port-number")
	flag.StringVar(&cfg.redisAddr, "redis-address", "localhost:6379", "Redis server address")
	flag.StringVar(&cfg.APIKey, "APIKey", os.Getenv("API_KEY"), "visualCrossing API key")

	flag.Parse()

	app := &application{
		config: cfg,
		rdb:    newRedisClient(cfg.redisAddr),
		ctx:    context.Background(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", app.fetchWeather)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on port %d (%s)", cfg.port, cfg.APIKey)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
