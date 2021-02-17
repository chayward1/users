package handlers

import (
	"net/http"
	"time"

	"chrishayward.xyz/users/messages"
	"chrishayward.xyz/users/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Authorize(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request.
		var r messages.AuthInfo
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Compare the secret.
		if r.Secret != secret {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the session.
		var s models.Session
		tx := db.First(&s, "token = ? AND expires >= ?",
			r.Token, time.Now().UnixNano())
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Return the user id.
		c.JSON(http.StatusOK, gin.H{
			"userID": s.UserID,
		})
	}
}
