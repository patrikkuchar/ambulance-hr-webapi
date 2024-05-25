package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrikkuchar/ambulance-hr-webapi/internal/ambulance_hr"
	"github.com/patrikkuchar/ambulance-hr-webapi/internal/db_service"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update  middleware
	dbService := db_service.NewMongoService[ambulance_hr.UserDto](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		log.Println("Middleware is being executed")
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// Apply the middleware to the router (allow all origins)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	engine.Use(cors.New(config))

	// request routings
	ambulance_hr.AddRoutes(engine)
	engine.Run(":" + port)
}
