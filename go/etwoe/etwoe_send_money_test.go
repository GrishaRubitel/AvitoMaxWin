package etwoe

import (
	"avitomaxwin/models"
	server "avitomaxwin/server"
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Отпрвка корректного запроса для перевода денег меж людьми
func TestSendMoney_Ok(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	// Получения балансов пользователей до транзакции
	joePeachWas, deadPWas := selectFriends(t, db)

	router := server.StartServer(envMap, db)

	// Формирование тела запроса
	requestBody := map[string]interface{}{
		"toUser": "deadp47",
		"amount": 1,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Генерация запроса
	token := envMap["MONTH_LONG_TOKEN"]
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusOK)

	// Получения балансов пользователей после транзакции
	joePeachNow, deadPNow := selectFriends(t, db)

	// Оценка корректности транзакции
	if joePeachWas-joePeachNow != 1 && deadPNow-deadPWas != 1 {
		t.Errorf("API not correct; jp (sender) was - %v, now - %v; dp (recip) was - %v, dp now - %v", joePeachWas, joePeachNow, deadPWas, deadPNow)
	}
}

// Отпрвка запроса для перевода денег меж людьми без нужного токена
func TestSendMoney_Unauthorized(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	// Некорректный токен
	req.Header.Set("Authorization", "Bearer amongu.amongus.amongus")
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusUnauthorized)
}

// Отпрвка запроса для перевода денег меж людьми с некорректным телом
func TestSendMoney_BadRequest(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	requestBody := map[string]interface{}{
		"toUser": "deadpool hitman 47", // Такого пользователя не существует
		"amount": 1,                    // Можно баланс в минус увести
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	token := envMap["MONTH_LONG_TOKEN"]
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusBadRequest)
}

// Отправка запроса без тела запроса
func TestSendMoney_500(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	token := envMap["MONTH_LONG_TOKEN"]
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusInternalServerError)
}

// Извлечение информации о пользователях из БД
func selectFriends(t *testing.T, db *gorm.DB) (int, int) {
	var joePeach, deadP models.UserCash

	rows, err := db.Raw("SELECT uc.login, uc.cash FROM users_cash uc WHERE uc.login IN ('joe_peach', 'deadp47')").Rows()
	if err != nil {
		t.Errorf("Error while executing SQL, error - %v", err)
		return 0, 0
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserCash
		if err := rows.Scan(&user.Login, &user.Cash); err != nil {
			t.Errorf("Error scanning row: %v", err)
			return 0, 0
		}

		switch user.Login {
		case "joe_peach":
			joePeach = user
		case "deadp47":
			deadP = user
		default:
			t.Errorf("Unexpected user in result: %s", user.Login)
		}
	}

	return joePeach.Cash, deadP.Cash
}
