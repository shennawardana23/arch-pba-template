package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

func TestOpenConnection(t *testing.T) {
	// config.Init(state.App.Environment)
	// var dbConfig = config.Get().Database

	// mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	// 	dbConfig.Username,
	// 	dbConfig.Password,
	// 	dbConfig.Host,
	// 	dbConfig.Port,
	// 	dbConfig.DbName,
	// )

	mysqlInfo := "shenna:Aqilah@21@tcp(localhost:3306)/arch_db"
	fmt.Println(mysqlInfo)

	db, err := sql.Open("mysql", mysqlInfo)
	panicOnError(err)
	if err = db.Ping(); err != nil {
		logger.Fatal(context.TODO(), err)
		if err = db.Close(); err != nil {
			logger.Fatal(context.TODO(), err)
		}
	}

	defer db.Close()
}

func panicOnError(err error) {
	if err != nil {
		log.Printf("panic on config %v", err)
		panic(err)
	}
}
