package main

import (
	"fmt"
	"log"
	"net/http"

	user_routes "github.com/umerwaheed/backend_golang/routes"
)

func main() {
	fmt.Println("Server is getting started...")
	r := user_routes.Router()
	fmt.Println("Listening at port 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))
}
