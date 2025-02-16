package etwoe

import (
	server "avitomaxwin/server"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetInfoGood_Ok(t *testing.T) {
	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	token := envMap["MONTH_LONG_TOKEN"]
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/info", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusOK)

	var response map[string]interface{}
	err = json.NewDecoder(terver.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	coinsFromResp, ok := response["coins"]
	if !ok {
		t.Error("No coins field in response")
	}

	var cash int
	err = db.Raw("SELECT cash FROM users_cash WHERE login = ?", "joe_peach").Scan(&cash).Error
	if err != nil {
		t.Error("Error while executing SQL, error - ", err)
	}

	if int(coinsFromResp.(float64)) != cash {
		t.Errorf("No coins field in response, cash from resp - %v, cash from db - %v", coinsFromResp, cash)
	}
}

func TestGetInfo_Unauthorized(t *testing.T) {
	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/info", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer amongus.amongus.amongus")
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusUnauthorized)

	var errorResponse map[string]interface{}
	err = json.NewDecoder(terver.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}
