package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"chrishayward.xyz/users/messages"
	"github.com/gin-gonic/gin"
)

func Authorize(endpoint, secret string) gin.HandlerFunc {
	type response struct {
		UserID int64 `json:"userID"`
	}

	return func(c *gin.Context) {
		// Read the token header.
		t := c.GetHeader("TOKEN")
		if len(t) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Create the request.
		m := &messages.AuthInfo{
			Secret: secret,
			Token:  t,
		}
		b, _ := json.Marshal(m)
		r, _ := http.NewRequest("GET", fmt.Sprintf(endpoint), bytes.NewBuffer(b))

		// Handle the response.
		p, _ := http.DefaultClient.Do(r)
		if p.StatusCode != http.StatusOK {
			c.AbortWithStatus(p.StatusCode)
			return
		}

		// Read the user id.
		var s response
		defer p.Body.Close()
		json.NewDecoder(p.Body).Decode(&s)

		// Set the user id in the context.
		c.Set("userID", s.UserID)

		// Call the next middleware.
		c.Next()
	}
}
