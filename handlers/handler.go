package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"lessons/calc"
	"lessons/models"
	"lessons/repo"
	"net/http"
)

func Transaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var transaction models.Transaction
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := repo.Create(transaction, db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			Commission(transaction, db)
			json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: transaction, Ok: true})

		case "GET":
			id := r.URL.Query().Get("id")
			if id != "" {
				transaction, err := repo.Read(id, db)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if transaction != nil {
					json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: *transaction, Ok: true})
				} else {
					http.NotFound(w, r)
				}
			} else {
				transactionList, err := repo.GetTransactionList(db)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(models.ListResponse{Transactions: transactionList, Ok: true})
			}

		case "PUT":
			id := r.URL.Query().Get("id")
			var alterTransaction models.AlterTransaction
			if err := json.NewDecoder(r.Body).Decode(&alterTransaction); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			transaction := models.Transaction{ID: id, Sum: alterTransaction.Sum,
				Currency: alterTransaction.Currency, Type: alterTransaction.Type, Category: alterTransaction.Category,
				Date: alterTransaction.Date, Description: alterTransaction.Description}

			if id != "" {
				repo.Alter(id, alterTransaction, db)
				json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: transaction, Ok: true})
				Commission(transaction, db)
			} else {
				json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: models.Transaction{}, Ok: false})
			}

		case "DELETE":
			id := r.URL.Query().Get("id")
			if id != "" {
				{
					repo.Drop(id, db)
					json.NewEncoder(w).Encode(models.DeleteResponse{ID: id, Ok: true})
				}
			} else {
				json.NewEncoder(w).Encode(models.TransactionResponse{Transaction: models.Transaction{}, Ok: false})
			}
		}
	}
}

func Commission(transaction models.Transaction, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		fmt.Print("Commission operation started")
		jsonRequest, err_ := json.Marshal(transaction)
		if err_ != nil {
			fmt.Print("Err with Marshalling")
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:8081/commissions/calculate", bytes.NewBuffer(jsonRequest))
		if err != nil {
			fmt.Print("Err with Request")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			json.NewEncoder(w).Encode(models.CommissionResponse{Commission: models.Commission{}, Ok: false})
		} else {
			commission := models.Commission{Transaction_id: transaction.ID, Commission: fmt.Sprintf("%f", calc.CalculateCommission(transaction)),
				Currency: transaction.Currency, Date: transaction.Date, Description: transaction.Description}
			if err := repo.CreateCommission(commission, db); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(models.CommissionResponse{Commission: commission, Ok: true})
		}
		defer resp.Body.Close()
	}
}
