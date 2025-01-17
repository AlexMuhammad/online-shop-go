package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/onlineshopfastcampus?sslmode=disable")
	if err != nil {
		fmt.Printf("Failed to connect db %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Failed to verify db connection %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Success to connect db \n")

}
