package handler

import (
	"database/sql"
	"log"
	"online-shop-fastcampus/model"

	"github.com/gin-gonic/gin"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: get order data from request
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			log.Printf("Something went wrong when read request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Product data is not valid"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)
		for _, o := range checkoutOrder.Products {
			ids = append(ids, o.ID)
			orderQty[o.ID] = o.Quantity
		}

		// TODO: get product data from db
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Something went wrong when get product: %v\n", err)
			c.JSON(500, gin.H{"error": "Something went wrong"})
			return
		}

		c.JSON(200, products)

		// TODO: Create a password

		// TODO: Hash a password

		// TODO: Create order & detail
	}
}
func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
