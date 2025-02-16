package api

import (
	models "avitomaxwin/models"
	"log"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetInfo(t *testing.T) {
	var user string

	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	user = "joe_peach"
	_, _, err = GetInfo(db, user)
	if err != nil {
		t.Errorf("unsuccessful atempt to get info about %v", user)
	}

	user = "ozon671games"
	_, _, err = GetInfo(db, user)
	if err == nil {
		t.Errorf("successful atempt to get info about %v", user)
	}
}

/*
Важное уточнени - в этом тесте не тестируется поиск по несуществующим пользователям,
так как аутентификация пользователя происходит во время валидации токена
и во время обращения к вышепротестированной функции
*/
func TestSelectTransactions(t *testing.T) {
	var user string
	var who string
	var where string

	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	user = "joe_peach"
	who = "t.sender"
	where = "t.recipient"
	transRecieve, err := selectTransactions[models.RecieveCoinRequest](db, user, who, where)
	if err != nil {
		t.Errorf("Unsuccessful attempt to get recieve coin transaction info about %v", user)
	}

	expected := reflect.TypeOf([]models.RecieveCoinRequest{})
	actual := reflect.TypeOf(transRecieve)
	if actual != expected {
		t.Errorf("Expected type %v, instead we got %v", expected, actual)
	}

	user = "joe_peach"
	who = "t.recipient"
	where = "t.sender = ? and t.transaction_type = 'transfer'"
	transSend, err := selectTransactions[models.SendCoinRequest](db, user, who, where)
	if err != nil {
		t.Errorf("Unsuccessful attempt to get send coin transaction info about %v", user)
	}

	expected = reflect.TypeOf([]models.SendCoinRequest{})
	actual = reflect.TypeOf(transSend)
	if actual != expected {
		t.Errorf("Expected type %v, instead we got %v", expected, actual)
	}

	user = "joe_peach"
	who = "i.sender"
	where = "t.recipient"
	_, err = selectTransactions[models.SendCoinRequest](db, user, who, where)
	if err == nil {
		t.Errorf("Successful attempt to get transaction info about %v, but with wrong SQL", user)
	}
}

func TestSelectInventor(t *testing.T) {
	var user string

	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	user = "joe_peach"
	items, err := selectInventory(db, user)
	if err != nil {
		t.Errorf("Unsuccessful attempt to get inventory of %v, error - %v", user, err)
	}

	expected := reflect.TypeOf([]models.ItemInfo{})
	actual := reflect.TypeOf(items)

	if actual != expected {
		t.Errorf("Expected type %v, instead we got %v", expected, actual)
	}
}
