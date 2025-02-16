package server

import (
	api "avitomaxwin/api"
	assistants "avitomaxwin/api/assistants"
	validator "avitomaxwin/api/validator"
	"errors"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(envMap map[string]string, db *gorm.DB) *gin.Engine {
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
		username, ok := assistants.ExtractUsername(ctx)
		if ok {
			code, resp, err := api.GetInfo(db, username)
			assistants.ResponseReturner(code, resp, err, ctx)
		}
	})

	apis.POST("/sendCoin", validator.ValidateToken, func(ctx *gin.Context) {
		code, params, err := assistants.ReadBodyData(ctx)
		if err != nil {
			assistants.ResponseReturner(code, "", err, ctx)
		}

		toUser, ok := params["toUser"]
		if !ok {
			assistants.ResponseReturner(http.StatusBadRequest, "", errors.New("no 'toUser' parameter"), ctx)
			return
		}

		amount, ok := params["amount"]
		if !ok {
			assistants.ResponseReturner(http.StatusBadRequest, "", errors.New("no 'amount' parameter"), ctx)
			return
		}

		username, ok := assistants.ExtractUsername(ctx)
		if ok {
			code, err = api.PostSendCoin(db, toUser, username, amount)
			assistants.ResponseReturner(code, "", err, ctx)
		}
	})

	apis.GET("/buy/:item", validator.ValidateToken, func(ctx *gin.Context) {
		username, ok := assistants.ExtractUsername(ctx)
		if ok {
			item := ctx.Param("item")
			code, err := api.GetBuy(db, item, username)
			assistants.ResponseReturner(code, "", err, ctx)
		}
	})

	apis.POST("/auth", func(ctx *gin.Context) {
		code, params, err := assistants.ReadBodyData(ctx)
		if err != nil {
			assistants.ResponseReturner(code, "", err, ctx)
			return
		}

		code, resp, err := api.PostAuth(db, params["username"], params["password"])
		assistants.ResponseReturner(code, resp, err, ctx)
	})

	return router
}
