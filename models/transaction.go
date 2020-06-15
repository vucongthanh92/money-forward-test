package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// User struct
type User struct {
	UserID   int    `gorm:"column:user_id;primary_key;AUTO_INCREMENT;not null" json:"user_id"`
	Username string `gorm:"column:username;type:nvarchar(100)" json:"username"`
	Password string `gorm:"column:password;type:nvarchar(100)" json:"password"`
}

// TableName func
func (User) TableName() string {
	return "users"
}

// Account struct
type Account struct {
	AccountID int     `gorm:"column:account_id;primary_key;AUTO_INCREMENT;not null" json:"account_id"`
	UserID    int     `gorm:"column:user_id;type:int(11)" json:"user_id"`
	Balance   float64 `gorm:"column:balance;type:decimal(18,5)" json:"balance"`
	Bank      string  `gorm:"column:bank;type:enum('vcb','acb','vib')" json:"bank"`
	Name      string  `gorm:"column:name;type:varchar(100)" json:"name"`
}

// TableName func
func (Account) TableName() string {
	return "accounts"
}

// Transaction struct
type Transaction struct {
	TransactionID   int     `gorm:"column:transaction_id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	AccountID       int     `gorm:"column:account_id;type:int(11)" json:"account_id"`
	Amount          float64 `gorm:"column:amount;type:decimal(18,5)" json:"amount"`
	Bank            string  `sql:"-" json:"bank"`
	TransactionType string  `gorm:"column:transaction_type;type:enum('withdraw','deposit')" json:"transaction_type"`
	CreatedAt       string  `gorm:"column:created_at;type:varchar(100)" json:"created_at"`
	Deleted         bool    `gorm:"column:deleted;type:varchar(100)" json:"-"`
	Updated         bool    `gorm:"column:updated;type:varchar(100)" json:"-"`
	ParentTxn       int     `gorm:"column:parent_txn;type:varchar(100)" json:"-"`
}

// TableName func
func (Transaction) TableName() string {
	return "transactions"
}

// Deposit method of transaction
func (txn *Transaction) Deposit(db *gorm.DB) error {
	var (
		account Account
		err     error
	)
	txn.getCreatedTime()
	err = db.Transaction(func(transDB *gorm.DB) error {
		var errTransDB error
		transDB.Exec(`SET AUTOCOMMIT = 0;`)
		errTransDB = transDB.Exec(`set global transaction isolation level repeatable read`).Error
		if errTransDB != nil {
			return errTransDB
		}
		errTransDB = transDB.Where("account_id = ?", txn.AccountID).First(&account).Error
		if errTransDB != nil {
			return errTransDB
		}
		account.Balance = account.Balance + txn.Amount
		errTransDB = transDB.Save(&account).Error
		if errTransDB != nil {
			return errTransDB
		}
		errTransDB = transDB.Create(&txn).Error
		if errTransDB != nil {
			return errTransDB
		}
		return nil
	})
	return err
}

// Withdraw method of transaction
func (txn *Transaction) Withdraw(db *gorm.DB) error {
	var (
		account Account
		err     error
	)
	txn.getCreatedTime()
	err = db.Transaction(func(transDB *gorm.DB) error {
		var errTransDB error
		transDB.Exec(`SET AUTOCOMMIT = 0;`)
		errTransDB = transDB.Exec(`set global transaction isolation level repeatable read`).Error
		if errTransDB != nil {
			return errTransDB
		}
		errTransDB = transDB.Where("account_id = ?", txn.AccountID).First(&account).Error
		if errTransDB != nil {
			return errTransDB
		}
		if account.Balance < txn.Amount {
			return errors.New("Balance is not enough")
		}
		account.Balance = account.Balance - txn.Amount
		errTransDB = transDB.Save(&account).Error
		if errTransDB != nil {
			return errTransDB
		}
		errTransDB = transDB.Create(&txn).Error
		if errTransDB != nil {
			return errTransDB
		}
		return nil
	})
	return err
}

// GetCreatedTime method
func (txn *Transaction) getCreatedTime() {
	txn.CreatedAt = fmt.Sprint(time.Now().Format("2006-02-01 15:04:00 -0700"))
}

func (txn *Transaction) getAmountTxn(db *gorm.DB) {
	var (
		subTxn []Transaction
	)
	db.Where("parent_txn = ?", txn.TransactionID).Find(&subTxn)
	if len(subTxn) > 0 {
		for _, item := range subTxn {
			if item.TransactionType == txn.TransactionType {
				txn.Amount += item.Amount
			} else {
				txn.Amount -= item.Amount
			}
		}
	}
}

// FormatTxn method transaction
func (txn *Transaction) FormatTxn(db *gorm.DB) {
	var (
		accout Account
	)
	db.Where("account_id = ?", txn.AccountID).First(&accout)
	db.First(&txn)
	txn.getAmountTxn(db)
	txn.Bank = accout.Bank
}

// UpdateTxn method of transaction
func (txn *Transaction) UpdateTxn(db *gorm.DB, amount float64, transactionType string) error {
	var (
		newTxn Transaction
		err    error
	)
	err = db.Where("deleted = 0").First(&txn).Error
	if err != nil {
		return err
	}
	newTxn.AccountID = txn.AccountID
	newTxn.Amount = amount
	newTxn.TransactionType = transactionType
	newTxn.ParentTxn = txn.TransactionID
	newTxn.getCreatedTime()
	txn.getAmountTxn(db)
	if txn.Amount < newTxn.Amount && txn.TransactionType != newTxn.TransactionType {
		return errors.New("Cannot update the amount greater than the transaction amount")
	}
	if newTxn.TransactionType == "deposit" {
		err = newTxn.Deposit(db)
	} else {
		err = newTxn.Withdraw(db)
	}
	if err == nil {
		db.Model(&txn).Update("updated", true)
	}
	return err
}

// DeleteTxn method transaction
func (txn *Transaction) DeleteTxn(db *gorm.DB) error {
	var (
		inverseTxn Transaction
		err        error
	)
	err = db.Where("parent_txn = 0 AND deleted = 0").First(&txn).Error
	if err != nil {
		return err
	}
	txn.getAmountTxn(db)
	inverseTxn.AccountID = txn.AccountID
	inverseTxn.Amount = txn.Amount
	inverseTxn.ParentTxn = txn.TransactionID
	inverseTxn.getCreatedTime()
	if txn.TransactionType == "deposit" {
		inverseTxn.TransactionType = "withdraw"
		err = inverseTxn.Withdraw(db)
	} else {
		inverseTxn.TransactionType = "deposit"
		err = inverseTxn.Deposit(db)
	}
	if err == nil {
		db.Model(&txn).Update("deleted", true)
	}
	return err
}

// ValidateRequest func
func (txn *Transaction) ValidateRequest(db *gorm.DB, requiredField map[string]bool) []string {
	var (
		err []string
	)
	if requiredField["transaction_id"] == true {
		resultRow := db.Where("transaction_id = ?", txn.TransactionID).First(&Transaction{})
		if resultRow.RowsAffected == 0 {
			err = append(err, "No transaction found")
		}
	}
	if requiredField["account_id"] == true {
		resultRow := db.Where("account_id = ?", txn.AccountID).First(&Account{})
		if resultRow.RowsAffected == 0 {
			err = append(err, "No account found")
		}
	}
	if requiredField["amount"] == true {
		if txn.Amount == 0 || txn.Amount > 9999999 {
			err = append(err, "Amount not allowed")
		}
	}
	if requiredField["transaction_type"] == true {
		if txn.TransactionType != "deposit" && txn.TransactionType != "withdraw" {
			err = append(err, "Transaction Type is incorrect")
		}
	}
	return err
}

// GetTransaction method
func (txn *Transaction) GetTransaction() Transaction {
	return *txn
}

// SetTransaction method
func (txn *Transaction) SetTransaction(paramTransaction Transaction) {
	txn.TransactionID = paramTransaction.TransactionID
	txn.AccountID = paramTransaction.AccountID
	txn.Amount = paramTransaction.Amount
	txn.Bank = paramTransaction.Bank
	txn.TransactionType = paramTransaction.TransactionType
	txn.CreatedAt = paramTransaction.CreatedAt
	txn.Deleted = paramTransaction.Deleted
	txn.Updated = paramTransaction.Updated
	txn.ParentTxn = paramTransaction.ParentTxn
}

// ResetTransaction method
func (txn *Transaction) ResetTransaction() {
	txn.TransactionID = 0
	txn.AccountID = 0
	txn.Amount = 0
	txn.Bank = ""
	txn.TransactionType = ""
	txn.CreatedAt = ""
	txn.Deleted = false
	txn.Updated = false
	txn.ParentTxn = 0
}
