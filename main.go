package main

import (
	"github/billing-engine-main/config"
	"github/billing-engine-main/internal/controllers"
	"github/billing-engine-main/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Panic(err)
	}
}

func main() {

	// Initialize database
	db := config.InitDB()

	services := services.NewBillingService(db)

	controllers := controllers.NewController(services)

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := router.Group("/api/v1")
	{
		loans := v1.Group("/loans")
		{
			loans.POST("", controllers.CreateLoan)
			loans.POST("/:id/payments", controllers.MakePayment)
			loans.GET("/:id/outstanding", controllers.GetOutstanding)
			loans.GET("/:id/delinquency", controllers.GetDelinquencyStatus)
		}
	}

	log.Printf("Server starting on :8088")
	if err := router.Run(":8088"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
