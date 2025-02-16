package etwoe

import (
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Тест соединения с базой данных
func TestDBConnection(t *testing.T) {
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		t.Errorf("Failed to read .env file, error - %v", err)
	}

	// Установка соединения
	dbTest, err := gorm.Open(postgres.Open(envMap["POSTGRES_LOCAL_CONN"]), &gorm.Config{})
	if err != nil {
		t.Errorf("Error while establishing connection to DB, error - %v", err)
	}

	// Выполнение простого запроса к базе данных
	if err := dbTest.Exec("SELECT NOW()").Error; err != nil {
		t.Errorf("Failed to execute simple SELECT operation")
	}
}
