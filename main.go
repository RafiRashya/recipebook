package main

import (
	"go-rest-modul/routes"
	"log"
	"net/http"
)

func main() {
	log.Println("Memulai server")

	routes := routes.RegisterRoutes()

	log.Fatal(http.ListenAndServe(":8080", routes))
}
