package handler

import (
	"database/sql"
	"online-shop-fastcampus/model"

	"github.com/gin-gonic/gin"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := model.SelectProduct(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		c.JSON(200, products)
	}
}
