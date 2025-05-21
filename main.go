package main

import (
	"fmt"
	"log"
	"net/http"

	routes "github.com/umerwaheed/backend_golang/routes"
)

func main() {
	fmt.Println("Server is getting started...")
	r := routes.Router()
	fmt.Println("Listening at port 5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}
