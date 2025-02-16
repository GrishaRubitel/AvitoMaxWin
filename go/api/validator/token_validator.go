package validator

import (
	cl "avitomaxwin/curloger"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var secret []byte

func ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortValidation(c, "no authorization header")
		cl.Log(logrus.InfoLevel, "No authorization header", map[string]interface{}{})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		abortValidation(c, "invalid token format")
		cl.Log(logrus.InfoLevel, "Invalid token format", map[string]interface{}{})
		return
	}

	token, err := jwt.Parse(tokenString, extractSecret(secret))
	if err != nil || !token.Valid {
		abortValidation(c, "token validation error")
		time, _ := token.Claims.GetExpirationTime()
		cl.Log(logrus.WarnLevel, "Wrong token", map[string]interface{}{
			"user":       token.Claims.(jwt.MapClaims)["username"].(string),
			"expiration": time,
			"error":      err,
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		c.Set("username", username)
	}

	c.Next()
}

func GenerateSecret(keyword string) {
	secret = []byte(keyword)
}

func extractSecret(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			cl.Log(logrus.InfoLevel, "unexpected signing method", map[string]interface{}{
				"method": token.Header["alg"],
			})
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	}
}

func abortValidation(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}
