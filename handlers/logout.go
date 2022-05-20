package handlers

import (
	"net/http"
	"time"

	"github.com/chayward1/users/messages"
	"github.com/chayward1/users/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request.
		var r messages.SessionInfo
		if err := c.ShouldBindJSON(&r); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Find the session.
		var s models.Session
		now := time.Now().UnixNano()
		tx := db.First(&s, "token = ? AND expires >= ?",
			r.Token, now)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Force expiration.
		s.Expires = now
		tx = db.Save(&s)
		if tx.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
