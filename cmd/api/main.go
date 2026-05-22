package main

import (
	gconfig "ginframework/intarnel/auth/config"
	"ginframework/intarnel/auth/handlers"
	"ginframework/intarnel/auth/middleware"
	"ginframework/intarnel/config"
	"ginframework/intarnel/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

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
		log.Fatal("Faild to connect to databse:", err)
	}
	defer pool.Close()

	r := gin.Default()
	r.GET("/login", handlers.Login)
	r.GET("/callbackfromgoogle", handlers.CallBackFromGoogle(pool))

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profiel", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "authorized"})
		})
	}
	r.Run(":" + cfg.Port)

}
