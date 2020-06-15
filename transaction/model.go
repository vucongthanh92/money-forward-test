package transaction

import (
	"money-forward-test/models"

	"github.com/jinzhu/gorm"
)

// TxnController interface
type TxnController interface {
	Deposit(db *gorm.DB) error
	Withdraw(db *gorm.DB) error
	FormatTxn(db *gorm.DB)
	SetTransaction(models.Transaction)
	GetTransaction() models.Transaction
	UpdateTxn(db *gorm.DB, amount float64, transactionType string) error
	DeleteTxn(db *gorm.DB) error
	ValidateRequest(db *gorm.DB, requiredField map[string]bool) []string
	ResetTransaction()
}

var (
	txnController  TxnController
	txnControllers []models.Transaction
)

func init() {
	txnController = &models.Transaction{}
}
