package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	port   = flag.Uint("port", 8080, "-port=8080")
	days   = flag.Int("days", 1, "-days=1")
	cost   = flag.Int("cost", bcrypt.DefaultCost, "-cost=14")
	debug  = flag.Bool("debug", true, "-debug=true")
	secret = flag.String("secret", uuid.NewString(), "-secret=?")
)

type User struct {
	gorm.Model
	Name     string
	Hash     string
	Sessions []Session
}

type Session struct {
	gorm.Model
	Token   string
	Expires int64
	UserID  uint
}

type Request struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pass string `form:"pass" json:"pass" binding:"required"`
}

func main() {
	flag.Parse()

	if *debug {
		fmt.Println(*secret)
	}

	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Session{})

	r := gin.Default()

	r.POST("/register", func(c *gin.Context) {
		var r Request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var u User
		tx := db.First(u, "name = ?", r.Name)
		if tx.RowsAffected != 0 {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(r.Pass), *cost)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		u.Name = r.Name
		u.Hash = string(bytes)

		tx = db.Create(&u)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	r.POST("/authenticate", func(c *gin.Context) {
		var r Request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var u User
		tx := db.First(&u, "name = ?", r.Name)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(r.Pass)); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s := &Session{
			Token:   uuid.NewString(),
			Expires: time.Now().AddDate(0, 0, *days).UnixNano(),
			UserID:  u.ID,
		}

		tx = db.Create(&s)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   s.Token,
			"expires": s.Expires,
		})
	})

	r.Run(fmt.Sprintf(":%d", *port))
}
