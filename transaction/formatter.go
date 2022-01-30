package transaction

import "time"

type CampaignTransationFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransationFormatter {
	return CampaignTransationFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransationFormatter {
	if len(transactions) == 0 {
		return []CampaignTransationFormatter{}
	}

	var transactionsFormatter []CampaignTransationFormatter
	for _, transaction := range transactions {
		transactionsFormatter = append(transactionsFormatter, FormatCampaignTransaction(transaction))
	}
	return transactionsFormatter
}
