package api

import (
	cl "avitomaxwin/curloger"
	"errors"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// API - POST - отправка денег пользователю
func PostSendCoin(db *gorm.DB, recipient, sender, amount string) (code int, err error) {
	// Попытка конверсии денег из строкового представления в числовое
	coins, err := strconv.Atoi(amount)
	if err != nil {
		cl.Log(logrus.ErrorLevel, "money conversion error", map[string]interface{}{
			"error":  err,
			"amount": coins,
		})
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	// Вызов хранимой процедуры и передача в неё аргументов
	if err := db.Exec("SELECT send_coins(?, ?, ?)", sender, recipient, coins).Error; err != nil {
		cl.Log(logrus.ErrorLevel, "error while transfering item", map[string]interface{}{
			"error":     err,
			"sender":    sender,
			"recipient": recipient,
			"amount":    coins,
		})
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
