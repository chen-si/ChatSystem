package message

//这里我们定义几个用户在线状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type User struct {
	UserID         int      `json:"userid"`
	UserPWD        string   `json:"userpwd"`
	UserName       string   `json:"username"`
	UserStatus     int      `json:"userstatus"` //用户状态
	FriendNotExist bool     `json:"friendnotexist"`
	Friends        []Friend `json:"friends"`
}

type Friend struct {
	FriendID int    `json:"friendid"`
	Group    string `json:"group"`
}
