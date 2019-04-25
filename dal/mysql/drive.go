package mysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var DB *sql.DB
var NewPositionFMT *sql.Stmt

func InitMysql() {
	DB, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/ais")
	if err := DB.Ping(); err != nil {
		log.Println("connect database fail:", err)
		return
	}
	DB.SetConnMaxLifetime(time.Hour)
	DB.SetMaxIdleConns(100)
	DB.SetMaxOpenConns(100)
	log.Println("connect database success")

	err := errors.New("")
	NewPositionFMT, err = DB.Prepare("INSERT positiontest SET " +
		"MessageType=?," +
		"RepeatIndicator=?," +
		"MMSI=?," +
		"NavigationStatus=?," +
		"ROT=?," +
		"SOG=?," +
		"PositionAccuracy=?," +
		"Longitude=?," +
		"Latitude=?," +
		"COG=?," +
		"HDG=?," +
		"TimeStamp=?," +
		"ReservedForRegional=?," +
		"RAIMFlag=?," +
		"Year=?," +
		"Month=?," +
		"Day=?," +
		"Hour=?," +
		"Minute=?," +
		"Second=?")
	if err != nil {
		log.Println("NewPositionFMT error", err)
	}
}

func FreeConn() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("free database connect fail:", err)
		}
	}
}
