package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenValidator abstracts validating access tokens.
type TokenValidator interface {
	ValidateAccess(string) (int64, error)
}

func AuthMiddleware(v TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if ah == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		// expect Bearer <token>
		var tokenStr string
		if len(ah) > 7 && ah[:7] == "Bearer " {
			tokenStr = ah[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		uid, err := v.ValidateAccess(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", uid)
		c.Next()
	}
}
