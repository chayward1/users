package main

import (
	"flag"
	"fmt"

	"chrishayward.xyz/users/handlers"
	"chrishayward.xyz/users/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	days   = flag.Int("days", 1, "-days=1")
	port   = flag.Uint("port", 8080, "-port=8080")
	cost   = flag.Int("cost", bcrypt.DefaultCost, "-cost=14")
	file   = flag.String("file", "users.db", "-file=users.db")
	secret = flag.String("secret", uuid.NewString(), "-secret=?")
)

func main() {
	// Parse the flags.
	flag.Parse()

	// Initialize the database.
	db, err := gorm.Open(sqlite.Open(*file), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		panic(err)
	}

	// Setup request handlers.
	r := gin.Default()

	r.POST("/register", handlers.Register(db, *cost))
	r.POST("/authenticate", handlers.Authenticate(db, *days))
	r.GET("/authorize", handlers.Authorize(db, *secret))

	// Run the application.
	r.Run(fmt.Sprintf(":%d", *port))
}
