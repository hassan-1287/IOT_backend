package main

import (
	gconfig "ginframework/intarnel/auth/config"
	"ginframework/intarnel/auth/handlers"
	"ginframework/intarnel/config"
	"ginframework/intarnel/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Println("ok")
	}
	gconfig.Intauth()

	var pool *pgxpool.Pool

	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Println("Failed to connect to database:", err)
	}
	defer pool.Close()

	r := gin.Default()
	r.GET("/pong", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"msg": "ping"}) })
	r.GET("/login", handlers.Login)
	r.GET("/callbackfromgoogle", handlers.CallBackFromGoogle)
	r.Run(":" + cfg.Port)

}
