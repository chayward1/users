package handlers

import (
	"net/http"
	"time"

	"github.com/chayward1/users/messages"
	"github.com/chayward1/users/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Authenticate(db *gorm.DB, days int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request.
		var r messages.UserInfo
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Find the user.
		var u models.User
		tx := db.First(&u, "name = ?", r.Name)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Compare the password.
		if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(r.Pass)); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Create the session.
		s := &models.Session{
			Token:   uuid.NewString(),
			Expires: time.Now().AddDate(0, 0, days).UnixNano(),
			UserID:  u.ID,
		}
		tx = db.Create(&s)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Return the data.
		c.JSON(http.StatusOK, gin.H{
			"token":   s.Token,
			"expires": s.Expires,
		})
	}
}
