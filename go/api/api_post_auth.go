package api

import (
	cl "avitomaxwin/curloger"
	models "avitomaxwin/models"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret []byte

func PostAuth(db *gorm.DB, username, password string) (code int, resp string, err error) {
	if username == "" || password == "" {
		cl.Log(logrus.InfoLevel, "Wrong username or password", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return http.StatusBadRequest, "", errors.New("invalid user data")
	}
	var user models.User

	result := db.First(&user, "login = ?", username)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

	if comparePasswords(user.PassHash, password) {
		token, err := generateToken(username)
		if err != nil {
			cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
				"error": err,
			})
			return http.StatusInternalServerError, "", errors.New("error while generating jwt token")
		}

		return http.StatusOK, token, nil
	}

	cl.Log(logrus.WarnLevel, "Error while comparing password hashs", map[string]interface{}{
		"password from db":     user.PassHash,
		"password from client": password,
	})
	return http.StatusUnauthorized, "", errors.New("wrong password")
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(secret)
}

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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePasswords(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSecret(keyword string) {
	secret = []byte(keyword)
}
