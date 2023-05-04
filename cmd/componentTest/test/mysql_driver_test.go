package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func StartMySQLServer() {
	username := "root"
	password := "123456"
	dbname := "mall"
	dataSourceName := fmt.Sprintf("%s:%s@/%s", username, password, dbname)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

// StartMySQLServer deprecated, for the reason that package "github.com/go-sql-driver/mysql" is difficult to use,
// use "gorm.io/driver/mysql" instead, see gorm_test.go in detail
func TestMySQLDriver(t *testing.T) {
	go StartMySQLServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	exitCode := <-ch
	log.Printf("Exit Code:%v", exitCode)
}
