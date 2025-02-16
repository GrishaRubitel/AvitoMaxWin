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

// Валидатор получаемого токена
func ValidateToken(c *gin.Context) {
	// Извлечение нужного заголовка
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortValidation(c, "no authorization header")
		cl.Log(logrus.InfoLevel, "No authorization header", map[string]interface{}{})
		return
	}

	//Удаление префикса
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		abortValidation(c, "invalid token format")
		cl.Log(logrus.InfoLevel, "Invalid token format", map[string]interface{}{})
		return
	}

	// Конкретно валидация
	token, err := jwt.Parse(tokenString, extractSecret(secret))
	if err != nil || !token.Valid {
		abortValidation(c, "token validation error")
		time, err := token.Claims.GetExpirationTime()
		if err != nil {
			cl.Log(logrus.ErrorLevel, "Failed to get expiration time", map[string]interface{}{
				"error": err,
			})
			return
		}

		cl.Log(logrus.WarnLevel, "Wrong token", map[string]interface{}{
			"expiration": time,
			"error":      err,
		})
		return
	}

	// Извлечения составляющих токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		abortValidation(c, "token validation error")
		cl.Log(logrus.ErrorLevel, "Failed to parse token claims", map[string]interface{}{
			"error": err,
		})
		return
	}

	// Извлечение юзернейма из токена
	username, ok := claims["username"].(string)
	if !ok {
		abortValidation(c, "username claim is missing")
		cl.Log(logrus.ErrorLevel, "Failed to parse token claims", map[string]interface{}{
			"error": err,
		})
		return
	}

	// Занесение юзернейма в контекст операции
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("username", username)
	}

	// Переход к основной логике
	c.Next()
}

// "Занесение в логику" секретного слова
func GenerateSecret(keyword string) {
	secret = []byte(keyword)
}

// Извлечения секрета
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

// Отказ в обслуживании в случае неудачи
func abortValidation(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}
