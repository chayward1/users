package handlers

import (
	"net/http"

	"chrishayward.xyz/users/messages"
	"chrishayward.xyz/users/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, cost int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request.
		var r messages.UserInfo
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Check for an existing user.
		var u models.User
		tx := db.First(&u, "name = ?", r.Name)
		if tx.RowsAffected != 0 {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		// Generate a password hash.
		bytes, err := bcrypt.GenerateFromPassword([]byte(r.Pass), cost)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Create the user.
		u.Name = r.Name
		u.Hash = string(bytes)
		tx = db.Create(&u)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
