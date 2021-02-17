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

func main() {
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
		name, pass :=
			c.PostForm("name"),
			c.PostForm("pass")

		var u *User
		tx := db.First(&u, "name = ?", name)
		if tx.RowsAffected != 0 {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte(pass), *cost)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		u.Name = name
		u.Hash = string(bytes)

		tx = db.Save(u)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	r.POST("/login", func(c *gin.Context) {
		name, pass :=
			c.PostForm("name"),
			c.PostForm("pass")

		var u *User
		tx := db.First(&u, "name = ?", name)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(pass)); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		s := &Session{
			Token:   uuid.NewString(),
			Expires: time.Now().AddDate(0, 0, *days).UnixNano(),
			UserID:  u.ID,
		}

		tx = db.Save(&s)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   s.Token,
			"expires": s.Expires,
		})
	})

	r.GET("/authorize", func(c *gin.Context) {
		key, token :=
			c.Query("secret"),
			c.Query("token")

		if key != *secret {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var s *Session
		tx := db.First(&s, "token <> ? AND expires > ?", token, time.Now().UnixNano())
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userID": s.UserID,
		})
	})

	r.Run(fmt.Sprintf(":%d", *port))
}
