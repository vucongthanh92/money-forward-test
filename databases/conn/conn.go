package conn

import (
	mysql "money-forward-test/databases/drivers"

	"github.com/jinzhu/gorm"
)

// Connect func
func Connect() *gorm.DB {
	return mysql.Connect()
}
