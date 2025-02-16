package curloger

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

var logChannel = make(chan logrus.Entry, 100)
var log = logrus.New()
var debugLog = logrus.New()

// Создание log файлов и запуск пишущей горутины
func InitCurloger(logFilePath string) {
	logDir := filepath.Dir(logFilePath)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatal("error while creating log dir - ", err)
		}
	}

	file, err := os.OpenFile(logFilePath+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error while opening log file - ", err)
	}

	debugFile, err := os.OpenFile(logFilePath+"_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error while opening debug log file - ", err)
	}

	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC850,
	})

	debugLog.SetLevel(logrus.DebugLevel)
	debugLog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC850,
	})
	debugLog.SetOutput(debugFile)

	go processLogs(file, debugFile)
}

// Пишущий процесс
func processLogs(file *os.File, debugFile *os.File) {
	defer file.Close()
	defer debugFile.Close()

	for entry := range logChannel {
		log.WithFields(entry.Data).Log(entry.Level, entry.Message)

		if entry.Level >= logrus.ErrorLevel {
			entry.Data["stacktrace"] = string(debug.Stack())
			debugLog.WithFields(entry.Data).Log(entry.Level, entry.Message)
		}
	}
}

// Функция запуска процесса записи в файл
func Log(level logrus.Level, message string, fields map[string]interface{}) {
	entry := logrus.Entry{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Data:    fields,
	}

	select {
	case logChannel <- entry:
	default:
	}
}
