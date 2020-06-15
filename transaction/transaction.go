package transaction

import (
	"encoding/json"
	"io/ioutil"
	"money-forward-test/databases/conn"
	libs "money-forward-test/helpers"
	"strconv"

	"money-forward-test/models"

	"github.com/gin-gonic/gin"
)

// GetTransactions API
func GetTransactions(c *gin.Context) {
	var (
		status       = 200
		responseData = gin.H{}
	)
	db := conn.Connect()
	defer db.Close()
	userID := c.Param("user_id")
	accountID := c.Query("account_id")
	if accountID == "" {
		var (
			arrAccount   []models.Account
			arrAccountID []int
		)
		db.Where("user_id = ?", userID).Find(&arrAccount)
		for _, item := range arrAccount {
			arrAccountID = append(arrAccountID, item.AccountID)
		}
		db.Where("account_id IN (?) AND parent_txn = 0 AND deleted = 0", arrAccountID).Find(&txnControllers)
	} else {
		db.Where("account_id = ? AND parent_txn = 0", accountID).Find(&txnControllers)
	}
	for index := range txnControllers {
		txnControllers[index].FormatTxn(db)
	}
	responseData = gin.H{
		"status": status,
		"data":   txnControllers,
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// CreateTransaction func
func CreateTransaction(c *gin.Context) {
	var (
		status       = 200
		responseData = gin.H{}
		err          error
	)
	db := conn.Connect()
	defer db.Close()
	txnController.ResetTransaction()
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &txnController)
	arrError := txnController.ValidateRequest(db, map[string]bool{
		"account_id":       true,
		"amount":           true,
		"transaction_type": true,
	})
	if len(arrError) > 0 {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  arrError,
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}

	if txnController.GetTransaction().TransactionType == "deposit" {
		err = txnController.Deposit(db)
	} else if txnController.GetTransaction().TransactionType == "withdraw" {
		err = txnController.Withdraw(db)
	}

	if err != nil {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  err.Error(),
		}
	} else {
		status = 200
		txnController.FormatTxn(db)
		responseData = gin.H{
			"status": status,
			"data":   txnController.GetTransaction(),
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// UpdateTransaction func
func UpdateTransaction(c *gin.Context) {
	var (
		status       = 200
		responseData = gin.H{}
		err          error
		tempTxn      models.Transaction
	)
	// connect db and create avairiable db
	db := conn.Connect()
	defer db.Close()
	// reset value txnController
	txnController.ResetTransaction()
	transactionID := c.Param("transaction_id")
	body, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(string(body)), &tempTxn)
	tempTxn.TransactionID, err = strconv.Atoi(transactionID)
	if err != nil {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  err.Error(),
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}
	// set value txnController
	txnController.SetTransaction(tempTxn)
	arrError := txnController.ValidateRequest(db, map[string]bool{
		"transaction_id":   true,
		"amount":           true,
		"transaction_type": true,
	})
	if len(arrError) > 0 {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  arrError,
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}
	err = txnController.UpdateTxn(db, tempTxn.Amount, tempTxn.TransactionType)
	if err != nil {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  err.Error(),
		}
	} else {
		status = 200
		txnController.FormatTxn(db)
		responseData = gin.H{
			"status": status,
			"data":   txnController.GetTransaction(),
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}

// DeleteTransaction func
func DeleteTransaction(c *gin.Context) {
	var (
		status       = 200
		responseData = gin.H{}
		tempTxn      models.Transaction
		err          error
	)
	// connect to db and create avariabla db
	db := conn.Connect()
	defer db.Close()
	// reset value transaction current
	txnController.ResetTransaction()
	// get transaction ID and check value ID
	transactionID := c.Param("transaction_id")
	tempTxn.TransactionID, err = strconv.Atoi(transactionID)
	if err != nil {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  err.Error(),
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}
	// add value temTxn into txnCOntroller
	txnController.SetTransaction(tempTxn)
	arrError := txnController.ValidateRequest(db, map[string]bool{
		"transaction_id": true,
	})
	if len(arrError) > 0 {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  arrError,
		}
		libs.ResponseRESTAPI(responseData, c, status)
		return
	}
	err = txnController.DeleteTxn(db)
	if err != nil {
		status = 400
		responseData = gin.H{
			"status": 400,
			"error":  err.Error(),
		}
	} else {
		status = 200
		txnController.FormatTxn(db)
		responseData = gin.H{
			"status": status,
			"data":   txnController.GetTransaction(),
		}
	}
	libs.ResponseRESTAPI(responseData, c, status)
}
