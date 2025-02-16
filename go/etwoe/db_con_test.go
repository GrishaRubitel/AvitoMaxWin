package etwoe

import (
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDBConnection(t *testing.T) {
	envMap, err := godotenv.Read("./../../.env")
	if err != nil {
		t.Errorf("Failed to read .env file, error - %v", err)
	}

	dbTest, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		t.Errorf("Error while establishing connection to DB, error - %v", err)
	}

	if err := dbTest.Exec("SELECT NOW()").Error; err != nil {
		t.Errorf("Failed to execute simple SELECT operation")
	}
}
