package api

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostAuth(t *testing.T) {
	var user string
	var pass string

	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		return
	}

	secret := envMap["JWT_SECRET"]
	fmt.Println(secret)

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	user = "joe_peach"
	pass = "1233211331"
	_, _, err = PostAuth(db, user, pass)
	if err != nil {
		t.Errorf("Unsuccessful attempt to generate token for %v (pass - %v), error - %v", user, pass, err)
	}

	user = "user_" + time.Now().Format("15-04")
	pass = "1111"
	_, _, err = PostAuth(db, user, pass)
	if err != nil {
		t.Errorf("Unsuccessful attempt to sign in new user %v (pass - %v), error - %v", user, pass, err)
	}

	user = "joe_peach"
	pass = "123321133"
	_, _, err = PostAuth(db, user, pass)
	if err == nil {
		t.Errorf("Wrong password %v for user %v accepted", pass, user)
	}
}
