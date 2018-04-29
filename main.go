package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/erikdubbelboer/fasthttp"
	"github.com/go-redis/redis"
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"

	_ "zxq.co/ripple/misirlou-api/http/api"
)

func main() {
	// set up mysql db
	db, err := models.CreateDB(getenvdefault("MYSQL_DSN",
		"root@/misirlou?multiStatements=true&parseTime=true"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// set up redis
	rdb := redis.NewClient(&redis.Options{
		Network:  os.Getenv("REDIS_NETWORK"),
		Addr:     getenvdefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getenvint("REDIS_DB", 0),
	})

	port := getenvdefault("PORT", ":8511")

	fmt.Println("Listening on", port)

	handler := http.Handler(http.Options{
		DB:                 db,
		Redis:              rdb,
		OAuth2ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		OAuth2ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		BaseURL:            getenvdefault("BASE_URL", "http://localhost"+port),
		StoreTokensURL:     getenvdefault("STORE_TOKENS_URL", "http://localhost"),
	})
	err = fasthttp.ListenAndServe(port, handler)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getenvdefault(env, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	return v
}

func getenvint(env string, def int) int {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	i, _ := strconv.Atoi(v)
	return i
}
