package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizationRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"forbidden": "role not found",
			})
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid role tyoe",
			})
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				log.Printf("User %s has permission to access this resource", userRole)
				log.Println("Successfully authorized")
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "permission denied",
		})
	}
}
