package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type DBConn struct {
	*gorm.DB
}

func Connection(host, user, password, dbName string) (*DBConn, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbName)
	conn, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	conn.SingularTable(true)
	return &DBConn{conn}, nil
}
