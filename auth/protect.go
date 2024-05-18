package auth

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Protect(signature []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(s, "Bearer ")
		// use jwt.Parse
		// logging
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// check it same jwt.SigningMethodHMSC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return signature, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			aud := claims["aud"]
			c.Set("aud", aud)
		}

		c.Next()
	}
}
