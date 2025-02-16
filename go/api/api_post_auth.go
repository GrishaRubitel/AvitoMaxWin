package api

import (
	cl "avitomaxwin/curloger"
	models "avitomaxwin/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret []byte

// API - POST - получение JWT токена
func PostAuth(db *gorm.DB, username, password string) (code int, resp string, err error) {
	var new bool

	if username == "" || password == "" {
		cl.Log(logrus.InfoLevel, "Wrong username or password", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return http.StatusBadRequest, "", errors.New("invalid user data")
	}
	var user models.User

	// Установка факта существования пользователя в системе
	result := db.First(&user, "login = ?", username)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		new = true
		// Если пользователя нет в системе - он в ней автоматически регистрируется
		user, err = signInUser(db, username, password)
		if err != nil {
			cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
				"error": err,
			})
			return http.StatusInternalServerError, "", err
		}
	} else if result.Error != nil {
		cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
			"error": result.Error,
		})
		return http.StatusInternalServerError, "", errors.New("error while searching user")
	}

	// Если пользователя в эту итерацию не регистрировали, сравниваем хеши паролей, полученные
	// из бд и от пользователя
	if new || comparePasswords(user.PassHash, password) {
		// Вызов генератора токена
		token, err := generateToken(username)
		if err != nil {
			cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
				"error": err,
			})
			return http.StatusInternalServerError, "", errors.New("error while generating jwt token")
		}

		respMap := map[string]string{
			"token": token,
		}

		resp, err := json.Marshal(respMap)
		if err != nil {
			cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
				"error": err,
			})
			return http.StatusInternalServerError, "", errors.New("error while generating jwt token")
		}

		return http.StatusOK, string(resp), nil
	}

	cl.Log(logrus.WarnLevel, "Error while comparing password hashs", map[string]interface{}{
		"password from db":     user.PassHash,
		"password from client": password,
	})
	return http.StatusUnauthorized, "", errors.New("wrong password")
}

// Генератор и подпись токена
func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 720).Unix(),
	})

	return token.SignedString(secret)
}

// Регистрация пользователя
func signInUser(db *gorm.DB, username, password string) (models.User, error) {
	var newUser models.User
	newUser.Login = username

	hashedPass, err := hashPassword(password)
	if err != nil {
		cl.Log(logrus.ErrorLevel, "error while hashing password", map[string]interface{}{
			"error":    err,
			"username": username,
			"password": password,
		})
		return newUser, errors.New("error while creating user in databse")
	}
	newUser.PassHash = hashedPass

	if err := db.Create(&newUser).Error; err != nil {
		cl.Log(logrus.ErrorLevel, "error adding user in user databse", map[string]interface{}{
			"error":    err,
			"username": username,
		})
		return newUser, errors.New("error while creating users in database")
	}

	// Кладём косарик в кошелёк нового пользователя
	if err := db.Table("users_cash").Create(models.UserCash{
		Login: username,
		Cash:  1000,
	}).Error; err != nil {
		cl.Log(logrus.ErrorLevel, "error adding user in user_cash databse", map[string]interface{}{
			"error":    err,
			"username": username,
		})
		return newUser, errors.New("error while creating user in database")
	}

	return newUser, nil
}

// Хеширование пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Сравнение хешей
func comparePasswords(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Чтение секретного слова из переменных окружения
func GenerateSecret(keyword string) {
	secret = []byte(keyword)
}
