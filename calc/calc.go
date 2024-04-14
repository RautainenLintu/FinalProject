package calc

import (
	"lessons/models"
	"strconv"
)

func CalculateCommission(transaction models.Transaction) float64 {
	// ● Если транзакция имеет type == перевод и валюта == USD, то комиссия составит 2% в валюте самой транзакции. (commission := req.Amount * 0.02)
	// ● Если транзакция имеет type == перевод и валюта == RUB, то комиссия составит 5% в валюте самой транзакции. (commission := req.Amount * 0.05)
	// ● Если транзакция имеет type == покупка или пополнение, то комиссию рассчитывать не надо.
	if transaction.Type != "transfer" {
		return 0.0
	}
	if amount, err := strconv.ParseFloat(transaction.Sum, 64); err == nil {
		switch transaction.Currency {
		case "USD":
			{
				return amount * 0.02
			}
		case "RUB":
			{
				return amount * 0.05
			}
		default:
			{
				return 0.0
			}
		}
	} else {
		return 0.0
	}
}
