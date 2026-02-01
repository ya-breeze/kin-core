package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ya-breeze/kin-core/auth"
	"github.com/ya-breeze/kin-core/db"
	"github.com/ya-breeze/kin-core/middleware"
	"github.com/ya-breeze/kin-core/models"
	"gorm.io/gorm"
)

// List is an application-specific model
type List struct {
	models.TenantModel
	Title string `json:"title"`
}

var gormDB *gorm.DB
var jwtSecret = []byte("my-secret")

func main() {
	// 1. Initialize DB
	var err error
	gormDB, err = db.Init(":memory:")
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	if err := gormDB.AutoMigrate(&models.Family{}, &models.User{}, &List{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 2. Setup Router
	r := gin.Default()

	// Public: Login/Token generation
	r.POST("/login", func(c *gin.Context) {
		// Mock login: always returns a token for User 1 in Family 100
		token, _ := auth.GenerateToken(1, 100, jwtSecret, time.Hour)
		c.JSON(200, gin.H{"token": token})
	})

	// Protected: Uses shared middleware
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/lists", func(c *gin.Context) {
			familyID := c.MustGet("family_id").(uint)
			var lists []List
			// Enforce isolation via scope
			gormDB.Scopes(db.Scope(familyID)).Find(&lists)
			c.JSON(200, lists)
		})

		api.POST("/lists", func(c *gin.Context) {
			familyID := c.MustGet("family_id").(uint)
			var list List
			if err := c.ShouldBindJSON(&list); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			list.FamilyID = familyID // Set the tenant ID
			gormDB.Create(&list)
			c.JSON(201, list)
		})
	}

	log.Println("Example app running on :8080")
	// r.Run(":8080") // Uncomment to run
}
