package cmd

import (
	"database/sql"
	"lessons/handlers"
	"log"
	"net/http"
)

func Run(db *sql.DB) {
	http.HandleFunc("/transaction", handlers.Transaction(db))
	/// http.HandleFunc("/commissions/calculate", handlers.Commission(db))
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
