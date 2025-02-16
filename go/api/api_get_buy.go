package api

import (
	cl "avitomaxwin/curloger"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetBuy(db *gorm.DB, item, username string) (code int, err error) {
	if err := db.Exec("SELECT buy_item(?, ?)", username, item).Error; err != nil {
		cl.Log(logrus.InfoLevel, "error while buying item", map[string]interface{}{
			"error":    err.Error,
			"username": username,
			"item":     item,
		})
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
