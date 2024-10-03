package middleware

import (
	"fmt"
	"golang/initializers"
	"golang/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Split Bearer and token
		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse the token
		tokenString := authToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		// Handle token parsing errors
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Get user from DB using ID from token claims (sub)
		var user models.User
		initializers.DB.Where("ID=?", claims["sub"]).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract the user's role from the claims
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role claim is missing"})
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Check if the user's role matches one of the allowed roles
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this resource"})
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set the user in the Gin context to be accessed by the next handlers
		c.Set("currentUser", user)

		// Continue to the next middleware/handler
		c.Next()
	}
}
