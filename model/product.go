package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Product struct {
	ID        string `json:"id" binding:"len=0"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsDeleted *bool  `json:"is_deleted,omitempty" `
}

var (
	ErrorDBNil = errors.New("connection is not available")
)

func SelectProduct(db *sql.DB) ([]Product, error) {
	if db == nil {
		return nil, ErrorDBNil
	}
	query := `SELECT id, name, price FROM products WHERE is_deleted = false;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func SelectProductByID(db *sql.DB, id string) (Product, error) {
	if db == nil {
		return Product{}, ErrorDBNil
	}

	query := `SELECT id, name, price FROM products WHERE is_deleted = false AND id = $1;`

	var product Product
	row := db.QueryRow(query, id)
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func SelectProductIn(db *sql.DB, ids []string) ([]Product, error) {
	if db == nil {
		return nil, ErrorDBNil
	}

	placeholders := []string{}
	args := []any{}
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`SELECT id, name, price FROM products WHERE is_deleted = false AND id IN (%s);`, strings.Join(placeholders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func InsertProduct(db *sql.DB, product Product) error {
	if db == nil {
		return ErrorDBNil
	}

	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, &product.ID, &product.Name, &product.Price)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(db *sql.DB, product Product) error {
	if db == nil {
		return ErrorDBNil
	}
	query := `UPDATE products SET name=$1, price=$2 WHERE id=$3`
	_, err := db.Exec(query, &product.Name, &product.Price, &product.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(db *sql.DB, id string) error {
	if db == nil {
		return ErrorDBNil
	}
	query := `UPDATE products set is_deleted=TRUE WHERE id=$1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
