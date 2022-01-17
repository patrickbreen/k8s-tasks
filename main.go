package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
)

func main() {
	log.Println("App is started")
	fmt.Println("fmt App is started")
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	log.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
