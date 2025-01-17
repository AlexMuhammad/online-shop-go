package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"online-shop-fastcampus/handler"
	"os"

	"github.com/gin-gonic/gin"
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

	if _, err = migrate(db); err != nil {
		fmt.Printf("Failed to migrate table %v\n", err)
	}

	r := gin.Default()

	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout")

	r.POST("/api/v1/orders/:id/confirm")
	r.GET("/api/v1/orders/:id")

	r.GET("/admin/products")
	r.POST("/admin/products/:id")
	r.DELETE("/admin/products/:id")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err = server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to run server %v\n", err)
	}
}

func migrate(db *sql.DB) (sql.Result, error) {
	if db == nil {
		return nil, errors.New("connection is not available")
	}
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price BIGINT NOT NULL,
			is_deleted BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			address VARCHAR NOT NULL,
			passcode VARCHAR,
			paid_at TIMESTAMP,
			paid_bank VARCHAR(255),
			paid_account VARCHAR(255),
			grand_total BIGINT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS order_details (
			id VARCHAR(36) PRIMARY KEY,
			order_id VARCHAR(36) NOT NULL,
			product_id VARCHAR(36) NOT NULL,
			quantity INT NOT NULL,
			price BIGINT NOT NULL,
			total BIGINT NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE RESTRICT,
			FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT
		);
	`)
}
