package main

import (
	"log"
	"net/http"

	"github.com/chankei613/execution-ledger/internal/api"
	"github.com/chankei613/execution-ledger/internal/db"
)

func main() {
	conn, err := db.Init("execution-ledger.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Println("execution-ledger backend listening on :8421")
	if err := http.ListenAndServe(":8421", router); err != nil {
		log.Fatal(err)
	}
}
