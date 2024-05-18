package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
)

// create func AccessToken
func AccessToken(signature string) gin.HandlerFunc {
	// jwt.NewWithClaims
	return func(c *gin.Context){
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Audience: "Pungping",
		})
		// put license use SignedString
		ss, err := token.SignedString([]byte(signature))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	
		// send token to frontend
		c.JSON(http.StatusOK, gin.H{
			"token": ss,
		})
	} 
}
