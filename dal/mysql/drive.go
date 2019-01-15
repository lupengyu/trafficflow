package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitMysql() {
	DB, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/ais")
	if err := DB.Ping(); err != nil {
		log.Println("connect database fail:", err)
		return
	}
	log.Println("connect database success")
}

func FreeConn() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("free database connect fail:", err)
		}
	}
}
