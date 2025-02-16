package etwoe

import (
	server "avitomaxwin/server"
	"bytes"
	"encoding/json"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Выполнение корректного запроса с целью получения токена
func TestAuth_Ok(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	// Заполнение тела запроса данными о несуществующем пользователе
	requestBody := map[string]interface{}{
		"username": "user_" + time.Now().Format("15-04"),
		"password": "123",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Подготовока запроса к генератору токена
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	// Валидация кода ответа
	assert.Equal(t, terver.Code, http.StatusOK)

	var responseMap map[string]string

	err = json.Unmarshal(terver.Body.Bytes(), &responseMap)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Получение нового токена
	newToken, ok := responseMap["token"]
	if !ok {
		t.Fatal("Token not found in response")
	}

	// Подготовка запроса для получения информации о пользователе по новому токену
	req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/api/info", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+newToken)
	req.Header.Set("Accept", "application/json")

	terver = httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	// Валидация кода ответа
	assert.Equal(t, terver.Code, http.StatusOK)
}

// Выполнение некорректного запроса с пустыми параметрами тела запроса
func TestAuth_BadRequest(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	// Создание "пустого" тела запроса
	requestBody := map[string]interface{}{
		"username": "",
		"password": "",
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Генерация запроса с "пустым" телом
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	// Валидация овтета
	assert.Equal(t, terver.Code, http.StatusBadRequest)
}

// Выполнение запроса без тела запроса
func TestAuth_500(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusInternalServerError)
}
