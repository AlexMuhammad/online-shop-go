package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"online-shop-fastcampus/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

		// TODO: Create a password
		passcode := generatePasscode(5)

		// TODO: Hash a password
		hashCode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			log.Printf("Something went wrong when get product: %v\n", err)
			c.JSON(500, gin.H{"error": "Something went wrong"})
			return
		}

		hashCodeString := string(hashCode)

		// TODO: Create order & detail
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			Passcode:   &hashCodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}
		for _, p := range products {
			total := p.Price * int64(orderQty[p.ID])
			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     total,
			}
			details = append(details, detail)
			order.GrandTotal += total
		}
		if err := model.CreateOrder(db, order, details); err != nil {
			log.Printf("Failed to create order: %v\n", err)
			c.JSON(500, gin.H{"error": "Failed to create order"})
			return
		}
		orderWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}
		orderWithDetail.Order.Passcode = &passcode
		c.JSON(200, orderWithDetail)
	}
}

func generatePasscode(length int) string {
	charset := "1234567890ABCDEFGHIJKLMNOPRSTUVWXYZ"
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
