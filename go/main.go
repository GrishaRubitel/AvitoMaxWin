package main

import (
	cl "avitomaxwin/curloger"
	server "avitomaxwin/server"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Запуск потока и канала для логирования
	cl.InitCurloger("./../logs/" + time.Now().Format("02-01-2006") + "/" + time.Now().Format("15-04"))

	// Чтение переменных окружения
	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		cl.Log(logrus.FatalLevel, "error while reading .env file", map[string]interface{}{
			"error": err,
		})
	}

	// Установка соединения с БД, которое будет передаваться дальше в API
	// Это нужно для уменьшения затрат повторную установку соединений
	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		cl.Log(logrus.FatalLevel, "error while establishing db connection", map[string]interface{}{
			"error": err,
		})
	}

	// Сборка и конфигурация сервера API
	router := server.StartServer(envMap, db)

	// Запуск сервера
	if err := router.Run(envMap["APPLICATION_URL"]); err != nil {
		cl.Log(logrus.ErrorLevel, "Failed to start the server", map[string]interface{}{
			"error": err,
		})
	}
}
