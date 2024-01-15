package middleware

import (
	"net/http"
	"os"
	"strings"
	"ta-kasir/model/request"
	"ta-kasir/model/response"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		godotenv.Load()
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized, (token kosong)",
			})
			return
		}

		split := strings.Split(tokenString, " ")
		if len(split) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized (token tidak sesuai format)",
			})
			return
		}
		if split[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized (token tidak sesuai formar Bearer)",
			})
			return
		}

		var claims request.JwtClaim
		_, err := jwt.ParseWithClaims(split[1], &claims, func(token *jwt.Token) (interface{}, error)  {
			return []byte(os.Getenv("SECRET")), nil
		})
		c.Set("jwt_claims", claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
				Status: http.StatusUnauthorized,
				Message: err.Error(),
			})
			return	
		}
		c.Next()
	}
}



