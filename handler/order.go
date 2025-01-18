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
		// TODO: get id param
		id := c.Param("id")
		// TODO: read request body
		var requstBody model.Confirm
		if err := c.BindJSON(&requstBody); err != nil {
			log.Printf("Something went wrong %v\n", err)
			c.JSON(400, gin.H{"error": err})
			return
		}
		// TODO: get order data from db
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Something went wrong when get product %v\n", err)
			c.JSON(500, gin.H{"error": "Something went wrong with server"})
			return
		}
		if order.Passcode == nil {
			log.Println("Passcode is not valid")
			c.JSON(401, gin.H{"error": "not autorize to order"})
			return
		}

		// TODO: match passcode
		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(requstBody.Passcode)); err != nil {
			log.Println("Passcode is not match")
			c.JSON(401, gin.H{"error": "not autorize to order"})
			return
		}
		// TODO: make sure order which unpaid
		if order.PaidAt != nil {
			log.Println("Order had been paid")
			c.JSON(400, gin.H{"error": "Order had been paid"})
			return
		}
		// TODO: match amount orders
		if order.GrandTotal != requstBody.Amount {
			log.Println("Check your order amout")
			c.JSON(400, gin.H{"error": "Check your order amout"})
			return
		}
		// TODO: update order information
		current := time.Now()
		if err = model.UpdateOrderByID(db, id, requstBody, current); err != nil {
			log.Println("Something went wrong when update order data")
			c.JSON(500, gin.H{"error": "Something went wrong with server"})
			return
		}
		order.PaidAccountNumber = &requstBody.AccountNumber
		order.PaidBank = &requstBody.Bank
		order.PaidAt = &current
		order.Passcode = nil

		c.JSON(200, order)
	}
}
func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		passcode := c.Query("passcode")

		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Something went wrong when get product %v\n", err)
			c.JSON(500, gin.H{"error": "Something went wrong with server"})
			return
		}

		if order.Passcode == nil {
			log.Println("Passcode is not valid")
			c.JSON(401, gin.H{"error": "not autorize to order"})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			log.Println("Passcode is not match")
			c.JSON(401, gin.H{"error": "not autorize to order"})
			return
		}

		order.Passcode = nil
		c.JSON(200, order)
	}
}
