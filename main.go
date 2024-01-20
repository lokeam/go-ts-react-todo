package main

import (
	"fmt"
	"go-ts-react-todo/router"
	"log"
	"net/http"
)

func main() {
	router := router.Router()

	fmt.Println("Starting server on port 6000...")
	log.Fatal(http.ListenAndServe(":6000", router))
}
