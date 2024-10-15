package main

import (
	api "go-fetch-backend/API"
	"net/http"
)

func main() {
	server := api.NewServer()
	http.ListenAndServe(":8000", server)
}
