package handler

import (
	"database/sql"
	"errors"
	"log"
	"online-shop-fastcampus/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request body
		var product model.Product
		if err := c.Bind(&product); err != nil {
			log.Printf("Something went wrong %v\n", err)
			c.JSON(400, gin.H{"error": err})
			return
		}
		product.ID = uuid.New().String()
		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Something went wrong when read request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data is not valid"})
			return
		}

		c.JSON(201, product)
	}
}
func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var reqBody model.Product
		if err := c.Bind(&reqBody); err != nil {
			log.Printf("Something went wrong %v\n", err)
			c.JSON(400, gin.H{"error": err})
			return
		}
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			log.Printf("Something went wrong when get product %v\n", err)
			c.JSON(500, gin.H{"error": "Something went wrong with server"})
			return
		}

		if reqBody.Name != "" {
			product.Name = reqBody.Name
		}
		if reqBody.Price != 0 {
			product.Price = reqBody.Price
		}
		if err := model.UpdateProduct(db, product); err != nil {
			log.Printf("Something went wrong when read request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data is not valid"})
			return
		}
		c.JSON(200, product)
	}
}
func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
