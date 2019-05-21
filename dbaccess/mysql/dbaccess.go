package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// GormConnect ...
func GormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"
	PASS := "pass"
	PROTOCOL := "tcp(127.0.0.1:3306)"
	DBNAME := "mydb"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		log.Print(err)
	}

	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(3)
	db.LogMode(true)
	return db
}
