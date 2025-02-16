package etwoe

import (
	server "avitomaxwin/server"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Отправка корректного запроса со всеми заголовками и параметрами для покупки предмета
func TestBuyItem_Ok(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	// Извлечение информации о предметах до запроса
	var quantityWas int
	err = db.Raw("SELECT quantity FROM inventory WHERE login = 'joe_peach' AND item = 'pen'").Scan(&quantityWas).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			quantityWas = 0
		} else {
			t.Errorf("Error while executing SQL, error - %v", err)
		}
	}

	router := server.StartServer(envMap, db)

	// Подготовка запроса с рабочим токеном
	token := envMap["MONTH_LONG_TOKEN"]
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/buy/pen", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	// Валидация кода ответа
	assert.Equal(t, terver.Code, http.StatusOK)

	// Извлечение информации о предметах после запроса
	var quantityNow int
	err = db.Raw("SELECT quantity FROM inventory WHERE login = 'joe_peach' AND item = 'pen'").Scan(&quantityNow).Error
	if err != nil {
		t.Error("Error while executing SQL, error - ", err)
	}

	// Сравнение того что было с тем что стало с целью проверки правильности операции
	if quantityNow-quantityWas != 1 {
		t.Errorf("API works incorrectly, new item wasn't added in joe_peach's inventory, was - %v, became - %v;", quantityWas, quantityNow)
	}
}

// Отправка запроса с некорректным токеном
func TestBuyItem_Unauthorized(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Fatal("Error while reading .env file")
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Fatal("Error while establishing db connection, error - ", err)
	}

	router := server.StartServer(envMap, db)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/buy/pen", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	//Неверный токен
	req.Header.Set("Authorization", "Bearer amongu.amongus.amongus")
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusUnauthorized)
}

// Отправка запроса с покупкой несуществующего/дорогого предмета
func TestBuyItem_BadRequest(t *testing.T) {
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
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/buy/lambargambor", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	terver := httptest.NewRecorder()
	router.ServeHTTP(terver, req)

	assert.Equal(t, terver.Code, http.StatusBadRequest)
}
