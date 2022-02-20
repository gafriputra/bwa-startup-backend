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

type UserTransactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	campaign := CampaignFormatter{
		Name:     transaction.Campaign.Name,
		ImageURL: "",
	}
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaign.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	return UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaign,
	}
}

func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	var transactionsFormatter []UserTransactionFormatter
	for _, transaction := range transactions {
		transactionsFormatter = append(transactionsFormatter, FormatUserTransaction(transaction))
	}
	return transactionsFormatter
}

type TransactionFormatter struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	UserID     int       `json:"user_id"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       string    `json:"code"`
	PaymentURL string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	return TransactionFormatter{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserID:     transaction.UserID,
		Status:     transaction.Status,
		PaymentURL: transaction.PaymentURL,
		Code:       transaction.Code,
		Amount:     transaction.Amount,
		CreatedAt:  transaction.CreatedAt,
	}
}
