package repo

import (
	"database/sql"
	"errors"
	"fmt"
	configs "lessons/config"
	"lessons/models"

	_ "github.com/lib/pq"
)

func InitDB(config *configs.Config) (*sql.DB, error) {
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName, config.Database.SSLMode)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println("Open err")
		return nil, err
	}

	createDB := `DROP TABLE IF EXISTS transactions;
	CREATE TABLE IF NOT EXISTS transactions (
		transaction_id VARCHAR(255),
		sum DECIMAL(10, 2),
		currency VARCHAR(255),
		type VARCHAR(255),
		category VARCHAR(255),
		date DATE,
		description VARCHAR(255)
	);
	DROP TABLE IF EXISTS commissions;
	CREATE TABLE IF NOT EXISTS commissions (
		transaction_id VARCHAR(255),
		commission DECIMAL(10, 2),
		currency VARCHAR(255),
		date DATE,
		description VARCHAR(255)
	);
	`

	_, err = db.Exec(createDB)
	if err != nil {
		fmt.Println("Exec err")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping err")
		return nil, err
	}

	return db, nil
}

func Create(transaction models.Transaction, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO transactions (transaction_id, sum, currency, type, category, date, description) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		transaction.ID, transaction.Sum, transaction.Currency, transaction.Type, transaction.Category, transaction.Date, transaction.Description)
	if err != nil {
		fmt.Println("Exec INSERT")
		return err
	}

	return nil
}

func Read(id string, db *sql.DB) (*models.Transaction, error) {
	// var result models.Transaction
	// db.QueryRow("SELECT transaction_id, sum, currency, type, category, date, description FROM transactions WHERE transaction_id = $1", id).Scan(&result)
	// return &result, nil
	rows, err := db.Query("SELECT * FROM transactions WHERE transaction_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	var transaction_id string
	var sum string
	var currency string
	var type_ string
	var category string
	var date string
	var description string
	err = rows.Scan(&transaction_id, &sum, &currency, &type_, &category, &date, &description)
	if err != nil {
		return nil, err
	}
	transaction := models.Transaction{ID: transaction_id, Sum: sum, Currency: currency, Type: type_, Category: category, Date: date, Description: description}
	if !rows.Next() {
		fmt.Printf("transaction_id: %s, sum: %s, currency: %s, type: %s, category: %s, date: %s, description: %s\n", transaction_id, sum, currency, type_, category, date, description)
		return &transaction, nil
	} else {
		err = errors.New("non-unique id error")
	}
	return nil, err
}

func Alter(id string, alterTransaction models.AlterTransaction, db *sql.DB) error {
	_, err := db.Exec("UPDATE transactions SET sum = $1, currency = $2, type = $3, category = $4, date = $5, description = $6 WHERE transaction_id = $7",
		alterTransaction.Sum, alterTransaction.Currency, alterTransaction.Type, alterTransaction.Category, alterTransaction.Date, alterTransaction.Description, id)
	if err != nil {
		fmt.Println("Exec UPDATE")
		return err
	}
	return nil
}

func Drop(id string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM transactions WHERE transaction_id = $1", id)
	if err != nil {
		fmt.Println("Exec DELETE")
		return err
	}
	return nil
}

func GetTransactionList(db *sql.DB) ([]models.Transaction, error) {
	rows, err := db.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactionList []models.Transaction
	for rows.Next() {
		var transaction_id string
		var sum string
		var currency string
		var type_ string
		var category string
		var date string
		var description string
		err := rows.Scan(&transaction_id, &sum, &currency, &type_, &category, &date, &description)
		if err != nil {
			return nil, err
		}
		transactionList = append(transactionList, models.Transaction{ID: transaction_id, Sum: sum, Currency: currency,
			Type: type_, Category: category, Date: date, Description: description})
		//fmt.Printf("transaction_id: %s, sum: %s, currency: %s, type: %s, category: %s, date: %s, description: %s\n", transaction_id, sum, currency, type_, category, date, description)
	}
	return transactionList, nil
}

func CreateCommission(commission models.Commission, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO commisssions (transaction_id, commmission, currency, date, description) VALUES ($1, $2, $3, $4, $5)",
		commission.Transaction_id, commission.Commission, commission.Currency, commission.Date, commission.Description)
	if err != nil {
		fmt.Println("Exec INSERT COMMISSION")
		return err
	}
	return nil
}
