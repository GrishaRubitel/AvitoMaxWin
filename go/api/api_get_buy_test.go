package api

import (
	"log"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetBuy(t *testing.T) {
	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	item := "pen"
	user := "joe_peach"
	code, err := GetBuy(db, item, user)
	if code != http.StatusOK && err != nil {
		t.Errorf("Unsuccessful atempt to buy item %v by %v, error - %v", item, user, err)
	}

	item = "pencil"
	user = "joe_peach"
	code, err = GetBuy(db, item, user)
	if code == http.StatusOK && err == nil {
		t.Errorf("Successful atempt to buy unexisting item %v by %v", item, user)
	}

	item = "cheesborg"
	user = "joe_peach"
	code, err = GetBuy(db, item, user)
	if code == http.StatusOK && err == nil {
		t.Errorf("Successful atempt to buy expensive item %v by %v", item, user)
	}

	item = "pen"
	user = "ozon671games"
	code, err = GetBuy(db, item, user)
	if code == http.StatusOK && err == nil {
		t.Errorf("Unsuccessful atempt to buy item %v by unexisting %v", item, user)
	}
}
