package messages

type UserInfo struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pass string `form:"pass" json:"pass" binding:"required"`
}
