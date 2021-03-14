package main

import (
	"flag"
	"fmt"

	"github.com/chayward1/users/handlers"
	"github.com/chayward1/users/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	days   = flag.Int("days", 1, "-days=1")
	port   = flag.Uint("port", 8080, "-port=8080")
	cost   = flag.Int("cost", bcrypt.DefaultCost, "-cost=14")
	secret = flag.String("secret", uuid.NewString(), "-secret=?")
	dbHost = flag.String("dbHost", "localhost", "-dbHost=localhost")
	dbPort = flag.Uint("dbPort", 5432, "-dbPort=5432")
	dbName = flag.String("dbName", "postgres", "-dbName=postgres")
	dbUser = flag.String("dbUser", "postgres", "-dbUser=postgres")
	dbPass = flag.String("dbPass", "postgres", "-dbPass=postgres")
)

func main() {
	// Parse the flags.
	flag.Parse()
	
	// Initialize the database.
	db, _ := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			*dbHost, *dbUser, *dbPass, *dbName, *dbPort)), &gorm.Config{})
	db.AutoMigrate(&models.User{}, *&models.Session{})

	// Run the application.
	r := gin.Default()
	r.POST("/register", handlers.Register(db, *cost))
	r.POST("/authenticate", handlers.Authenticate(db, *days))
	r.GET("/authorize", handlers.Authorize(db, *secret))
	r.Run(fmt.Sprintf(":%d", *port))
}
