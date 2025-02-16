package main

import (
	api "avitomaxwin/api"
	assistants "avitomaxwin/api/assistants"
	validator "avitomaxwin/api/validator"
	cl "avitomaxwin/curloger"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cl.InitCurloger("./../logs/" + time.Now().Format("02-01-2006") + "/" + time.Now().Format("15-04"))

	envMap, err := godotenv.Read("./../.env")
	if err != nil {
		cl.Log(logrus.FatalLevel, "error while reading .env file", map[string]interface{}{
			"error": err,
		})
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(envMap["POSTGRES_CONN"]), &gorm.Config{})
	if err != nil {
		cl.Log(logrus.FatalLevel, "error while establishing db connection", map[string]interface{}{
			"error": err,
		})
		log.Fatal(err)
	}

	api.GenerateSecret(envMap["JWT_SECRET"])
	validator.GenerateSecret(envMap["JWT_SECRET"])

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "POST"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		AllowAllOrigins: true,
	}))

	apis := router.Group("/api")

	apis.GET("/info", validator.ValidateToken, func(ctx *gin.Context) {
		code, resp, err := api.ApiGetInfo(db, ctx.Keys["username"].(string))
		assistants.ResponseReturner(code, resp, err, ctx)
	})

	apis.POST("/sendCoin", validator.ValidateToken, func(ctx *gin.Context) {
		code, params, err := assistants.ReadBodyData(ctx)
		if err != nil {
			assistants.ResponseReturner(code, "", err, ctx)
		}

		code, err = api.ApiPostSendCoin(db, params["toUser"], ctx.Keys["username"].(string), params["amount"])
		assistants.ResponseReturner(code, "", err, ctx)
	})

	apis.GET("/buy/:item", validator.ValidateToken, func(ctx *gin.Context) {
		item := ctx.Param("item")
		code, err := api.ApiGetBuy(db, item, ctx.Keys["username"].(string))
		assistants.ResponseReturner(code, "", err, ctx)
	})

	apis.POST("/auth", func(ctx *gin.Context) {
		code, params, err := assistants.ReadBodyData(ctx)
		if err != nil {
			assistants.ResponseReturner(code, "", err, ctx)
		}

		code, resp, err := api.ApiPostAuth(db, params["username"], params["password"])
		assistants.ResponseReturner(code, resp, err, ctx)
	})

	router.Run(envMap["APPLICATION_URL"])
}
