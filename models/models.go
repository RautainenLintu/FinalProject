package models

type Transaction struct {
	ID          string `json:"id"`
	Sum         string `json:"sum"`
	Currency    string `json:"currency"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type AlterTransaction struct {
	Sum         string `json:"sum"`
	Currency    string `json:"currency"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	Ok          bool        `json:"ok"`
}

type ListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Ok           bool          `json:"ok"`
}

type DeleteResponse struct {
	ID string `json:"id"`
	Ok bool   `json:"ok"`
}

type Commission struct {
	Transaction_id string `json:"transaction_id"`
	Commission     string `json:"commission"`
	Currency       string `json:"currency"`
	Date           string `json:"date"`
	Description    string `json:"description"`
}

type CommissionResponse struct {
	Commission Commission `json:"commission"`
	Ok         bool       `json:"ok"`
}
