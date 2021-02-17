package messages

type AuthInfo struct {
	Secret string `form:"secret" json:"secret" binding:"required"`
	Token  string `form:"token" json:"token" binmding:"required"`
}
