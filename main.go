package main

import (
	"fmt"
	"os"

	"github.com/erikdubbelboer/fasthttp"
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
)

func main() {
	db, err := models.CreateDB(os.Getenv("MYSQL_DSN"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8511"
	}

	fmt.Println("Listening on", port)

	handler := http.Handler(http.Options{
		DB: db,
	})
	err = fasthttp.ListenAndServe(port, handler)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
