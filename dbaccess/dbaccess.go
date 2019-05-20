package dbaccess

import (
	"log"

	"github.com/jinzhu/gorm"
)

// GormConnect ...
func GormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"
	PASS := "Ah4vn3tetsuya"
	PROTOCOL := "tcp(192.168.9.4:3306)"
	DBNAME := "mydb"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		log.Print(err)
	}
	return db
}
