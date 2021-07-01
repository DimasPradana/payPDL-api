package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	"github.com/sirupsen/logrus"
)

func KonekMysql(host, username, pwd, port, dbname string) (*sql.DB, error) { // epajak
	connString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, pwd, host, port, dbname)
	db, err := sql.Open("mysql", connString)
	// TODO check if error
	if err = db.Ping(); err != nil {
		logrus.Fatalf("koneksi error bosku: %v", err)
		db.Close()
	}
	return db, nil
}

func KonekOracle(username, pwd, host, port, sn string) (*sql.DB, error) {
	connString := fmt.Sprintf("%v/%v@%v:%v/%v", username, pwd, host, port, sn)
	//connString := fmt.Sprintf("%v/%v@%v:%v/%v?charset=utf8mb4&parseTime=True&loc=Asia/Jakarta", username, pwd, host, port, sn)
	db, err := sql.Open("godror", connString)
	// db.SetMaxOpenConns(20) //  TODO snub on Wed 30 Jun 2021 11:17:55 : limit opened connections //
	// if err != nil {
	if err = db.Ping(); err != nil {
		logrus.Fatalf("koneksi error bosku: %v", err)
		// logrus.Infof("koneksi error bosku: %v", err)
		db.Close()
	}
	return db, nil
}

func KonekBphtb(host, username, pwd, port, dbname string) (*sql.DB, error) {
	connString := fmt.Sprintf("%v:%v@(%v:%v)/%v", username, pwd, host, port, dbname)
	db, err := sql.Open("mysql", connString)
	// if err != nil {
	if err = db.Ping(); err != nil {
		logrus.Fatalf("koneksi error bosku: %v", err)
		db.Close()
	}
	return db, nil
}
