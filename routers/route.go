package routers

import (
	"money-forward-test/transaction"

	"github.com/gin-gonic/gin"
)

// Route api
func Route(r *gin.RouterGroup) {
	r.GET("/users/:user_id/transactions", ValidateUserID(), transaction.GetTransactions)
	r.POST("/users/:user_id/transactions", ValidateUserID(), transaction.CreateTransaction)
	r.PATCH("/users/:user_id/transactions/:transaction_id", ValidateUserID(), transaction.UpdateTransaction)
	r.DELETE("/users/:user_id/transactions/:transaction_id", ValidateUserID(), transaction.DeleteTransaction)
}
