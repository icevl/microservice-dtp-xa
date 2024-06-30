package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func DBConnect() {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"tcp",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	Database, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully")

	// defer Database.Close()
}

func StartTransaction(uuid string) error {
	_, err := Database.Exec(fmt.Sprintf("XA START '%s'", uuid))
	if err != nil {
		return err
	}

	return nil
}

func EndTransaction(uuid string) error {
	_, err := Database.Exec(fmt.Sprintf("XA END '%s'", uuid))
	if err != nil {
		return err
	}

	return nil
}

func PrepareTransaction(uuid string) error {
	_, err := Database.Exec(fmt.Sprintf("XA PREPARE '%s'", uuid))
	if err != nil {
		return err
	}

	return nil
}

func RollbackTransaction(uuid string) error {
	_, err := Database.Exec(fmt.Sprintf("XA ROLLBACK '%s'", uuid))
	if err != nil {
		return err
	}

	return nil
}

func CommitTransaction(uuid string) error {
	_, err := Database.Exec(fmt.Sprintf("XA COMMIT '%s'", uuid))
	if err != nil {
		return err
	}

	return nil
}
