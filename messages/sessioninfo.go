package messages

type SessionInfo struct {
	Token  string `form:"token" json:"token" binding:"required"`
}
