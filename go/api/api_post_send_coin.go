package api

import (
	cl "avitomaxwin/curloger"
	"errors"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ApiPostSendCoin(db *gorm.DB, recipient, sender, amount string) (code int, err error) {
	coins, err := strconv.Atoi(amount)
	if err != nil {
		cl.Log(logrus.ErrorLevel, "money conversion error", map[string]interface{}{
			"error":  err.Error,
			"amount": coins,
		})
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	if err := db.Exec("SELECT send_coins(?, ?, ?)", sender, recipient, coins).Error; err != nil {
		cl.Log(logrus.InfoLevel, "error while transfering item", map[string]interface{}{
			"error":     err.Error,
			"sender":    sender,
			"recipient": recipient,
			"amount":    coins,
		})
		return http.StatusBadRequest, err
	} else {
		return http.StatusOK, nil
	}
}
