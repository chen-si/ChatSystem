package model

type User struct {
	UserID   int    `json:"userid"`
	UserPWD  string `json:"userpwd"`
	UserName string `json:"username"`
}
