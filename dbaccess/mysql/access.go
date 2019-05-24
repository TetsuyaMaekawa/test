package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// OpenMySQL MySQLとの接続を確立
func OpenMySQL() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:Ah4vn3tetsuya@/mydb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(3)
	return db, nil
}
