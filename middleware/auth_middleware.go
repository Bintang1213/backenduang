package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Silahkan login terlebih dahulu"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Mengambil secret langsung dari .env
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Memastikan method signing sesuai (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				// Konversi claims ke string secara aman
				role := fmt.Sprintf("%v", claims["role"])
				c.Set("role", role) 
				c.Next()
				return
			}
		}

		// Jika token salah atau secret tidak cocok, ngrok tidak akan crash lagi
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sesi habis"})
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Role tidak ditemukan"})
			return
		}

		userRole := val.(string)

		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		// Response JSON yang rapi jika Kasir mencoba Create/Update/Delete
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Akses ditolak: Role %s tidak diizinkan melakukan tindakan ini", userRole),
		})
	}
}