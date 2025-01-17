package handler

import (
	"database/sql"
	"errors"
	"log"
	"online-shop-fastcampus/model"

	"github.com/gin-gonic/gin"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := model.SelectProduct(db)
		if err != nil {
			log.Printf("Something went wrong %v\n", err)
			c.JSON(500, gin.H{"error": err})
			return
		}

		c.JSON(200, products)
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		product, err := model.SelectProductByID(db, id)
		log.Printf("value test %v, test err %v\n", product, err)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(404, gin.H{"error": "Product not found"})
				return
			}
			log.Printf("Something went wrong %v\n", err)
			c.JSON(500, gin.H{"error": err})
			return
		}

		c.JSON(200, product)
	}
}
