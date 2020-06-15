package mysql

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Connect func
func Connect() *gorm.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	for i := 1; i <= 3; i++ {
		db, err := gorm.Open("mysql", connString)
		if err == nil {
			return db
		}
		if i == 3 {
			panic(err.Error())
		}
	}
	return nil
}
