package main

import (
	"fmt"
	"log"
	"net/http"
	"sixam/go-ts-react-todo/router"
)

func main() {
	router := router.Router()

	fmt.Println("Starting server on port 6000...")
	log.Fatal(http.ListenAndServe(":6000", router))
}
