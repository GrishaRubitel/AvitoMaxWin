package api

import (
	cl "avitomaxwin/curloger"
	models "avitomaxwin/models"
	"errors"
	"net/http"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// API - GET - извлечение информации о пользователе
func GetInfo(db *gorm.DB, username string) (code int, resp string, err error) {
	var user models.User

	// Установка факта существования пользователя в системе
	// Если его нет - возвращаем 401
	result := db.First(&user, "login = ?", username)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		cl.Log(logrus.InfoLevel, "Username not in db", map[string]interface{}{
			"username": username,
			"error":    result.Error,
		})
		return http.StatusUnauthorized, "", errors.New("your username not exist")
	} else if result.Error != nil {
		cl.Log(logrus.ErrorLevel, "Internal server error", map[string]interface{}{
			"error": result.Error,
		})
		return http.StatusInternalServerError, "", errors.New("error while searching user")
	}

	var response models.InfoResponse

	// Извлечени информации о финансовом состоянии пользователя
	err = db.Table("users_cash").Select("cash").Where("login = ?", username).Scan(&response.Coins).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cl.Log(logrus.InfoLevel, "Username's cash not in db", map[string]interface{}{
			"username": username,
			"error":    err,
		})
		return http.StatusInternalServerError, "", errors.New("error while searching user")
	}

	// Извлечение информации об исходящих транзакциях
	response.CoinHistory.Sent, err = selectTransactions[models.SendCoinRequest](db, username, "t.recipient", "t.sender = ? and t.transaction_type = 'transfer'")
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	// Извлечени информации о входящих транзакциях
	response.CoinHistory.Received, err = selectTransactions[models.RecieveCoinRequest](db, username, "t.sender", "t.recipient = ?")
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	// Извлечение информации об инвентаре пользователя
	response.Inventory, err = selectInventory(db, username)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	// Приведение всей собранной информации в единый JSON документ
	jsonedData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		cl.Log(logrus.ErrorLevel, "Error while parsing transaction in JSON", map[string]interface{}{
			"error": err,
			"data":  response,
		})
		return http.StatusInternalServerError, "", errors.New("internal server error")
	}

	return http.StatusOK, string(jsonedData), nil
}

// Извлечение информации о транзакциях, в зависимости от её типа (вход./исход.)
func selectTransactions[T any](db *gorm.DB, username, who, where string) ([]T, error) {
	var transactions []T

	err := db.Select(who, "t.amount").Table("transactions t").Where(where, username).Scan(&transactions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		cl.Log(logrus.WarnLevel, "Error while searching for transactions in db", map[string]interface{}{
			"error":    err,
			"username": username,
		})
		return nil, errors.New("internal server error")
	}

	return transactions, nil
}

// Извлечение информации об инвентаре пользователя
func selectInventory(db *gorm.DB, username string) (inventory []models.ItemInfo, err error) {
	err = db.Select("i.item, i.quantity").Table("inventory i").Where("i.login = ?", username).Scan(&inventory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		cl.Log(logrus.InfoLevel, "Error while searching for inventory in db", map[string]interface{}{
			"error":    err,
			"username": username,
		})
		return nil, errors.New("internal server error")
	}

	return inventory, nil
}
