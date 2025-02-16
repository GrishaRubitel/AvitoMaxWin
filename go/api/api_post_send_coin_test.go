package api

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostSendCoin(t *testing.T) {
	var recipient string
	var sender string
	var amount string

	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	recipient = "deadp47"
	sender = "joe_peach"
	amount = "1"
	_, err = PostSendCoin(db, recipient, sender, amount)
	if err != nil {
		t.Errorf("Unsuccessful attempt to send %v coins from %v to %v, error - %v", amount, sender, recipient, err)
	}

	recipient = "deadp47"
	sender = "joe_peach"
	amount = "2147483647"
	_, err = PostSendCoin(db, recipient, sender, amount)
	if err == nil {
		t.Errorf("Successful attempt to send all %v coins from %v to %v", amount, sender, recipient)
	}

	recipient = "deadp47"
	sender = "joe_peach"
	amount = "ohe billion gazillion bucks"
	_, err = PostSendCoin(db, recipient, sender, amount)
	if err == nil {
		t.Errorf("Wrong money format, what is this %v???", amount)
	}

	recipient = "deadp47"
	sender = "ozon671games"
	amount = "1"
	_, err = PostSendCoin(db, recipient, sender, amount)
	if err == nil {
		t.Errorf("Sender %v is not in service", sender)
	}

	recipient = "recipient"
	sender = "joe_peach"
	amount = "1"
	_, err = PostSendCoin(db, recipient, sender, amount)
	if err == nil {
		t.Errorf("Recipient %v is not in service", recipient)
	}
}
