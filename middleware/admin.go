package middleware

import "github.com/gin-gonic/gin"

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "secret" // you can put this on env variable for better security approach

		// TODO: ambil header authorization
		auth := c.Request.Header.Get("Authorization")
		// TODO: Validate header
		if auth == "" {
			c.JSON(401, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		if auth != key {
			c.JSON(401, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		//TODO: Next to subsequent handler
		c.Next()
	}
}
