package payment

import (
	"bwa-startup/user"
	"errors"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct{}

type Service interface {
	GetPaymentResponse(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentResponse(transaction Transaction, user user.User) (string, error) {
	// 1. Initiate Snap client
	var sn = snap.Client{}
	sn.New("", midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).
	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		// CreditCard: &snap.CreditCardDetails{
		// 	Secure: true,
		// },
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			LName: "",
			Email: user.Email,
			Phone: "",
		},
	}
	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, err := sn.CreateTransaction(req)
	if err != nil {
		return "", errors.New(err.Message)
	}
	return snapResp.RedirectURL, nil
}
