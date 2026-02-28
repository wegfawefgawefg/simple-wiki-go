package main

import (
	"log"
	"net/http"
	"nofrills-wiki/internal/wiki"
)

func main() {
	addr := ":8080"
	log.Printf("wiki server running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, wiki.NewServer()))
}
