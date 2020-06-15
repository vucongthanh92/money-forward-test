package routers

import (
	"money-forward-test/databases/conn"
	"money-forward-test/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidateUserID func
func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			status int = 200
			user   models.User
			msg    string
		)
		if c.Param("user_id") == "" {
			status = 400
			msg = "No userID information"
		} else {
			userID, err := strconv.Atoi(c.Param("user_id"))
			if err != nil {
				status = 400
				msg = "User ID is incorrect"
			} else {
				db := conn.Connect()
				defer db.Close()
				resultRow := db.Where("user_id = ?", userID).First(&user)
				if resultRow.RowsAffected == 0 {
					status = 400
					msg = "No user found"
				}
			}
		}
		if status == 200 {
			c.Next()
		} else {
			responseData := gin.H{
				"status": status,
				"msg":    msg,
			}
			c.JSON(status, responseData)
			c.Abort()
		}
	}
}
