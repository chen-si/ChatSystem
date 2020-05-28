package message

//消息类型常量表示
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	GroupSmsMesType         = "GroupSmsMes"
	PrivateSmsMesType       = "PrivateSmsMes"
	SmsResMesType           = "SmsResMes"
	SmsReceiverMesType      = "SmsReceiverMes"
	AddFriendMesType        = "AddFriendMes"
	AddFriendResMesType     = "AddFriendResMes"
	DeleteFriendMesType     = "DeleteFriendMes"
	DeleteFriendResMesType  = "DeleteFriendResMes"
	LogOutMesType           = "LogOutMes"
	ListFriendsMesType      = "ListFriendsMes"
	ListFriendsResMesType   = "ListFriendsResMes"
	QueryRecordMesType      = "QueryRecordMes"
	QueryRecordResMesType   = "QueryRecordResMes"
)

//Message 的一般类型
type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
}

//登录消息
type LoginMes struct {
	UserID   int    `json:"userid"`   //用户ID
	UserPWD  string `json:"userpwd"`  //用户密码
	UserName string `json:"username"` //用户名
}

type LoginResMes struct {
	//登录结果消息
	//code：
	//500 ：未注册
	//200 ：登陆成功
	//403 ：密码错误
	//505 : 服务器内部错误
	Code     int    `json:"code"`    //返回状态码
	UsersID  []int  `json:"usersid"` //保存用户id的切片
	UserName string `json:"username"`
	Error    string `json:"error"` //错误消息
}

type RegisterMes struct {
	//注册消息
	User User `json:"user"`
}

type RegisterResMes struct {
	//400表示用户已占用
	//200表示注册成功
	Code  int    `json:"code"`
	Error string `json:"error"`
}

//为了配合服务器推送用户状态变化消息的类型
type NotifyUserStatusMes struct {
	UserID int `json:"userid"` //用户id
	Status int `json:"status"` //用户状态
}

//增加一个smsMes 发送消息
type GroupSmsMes struct {
	Content  string `json:"content"`
	GroupKey string `json:"groupkey"`
	User
}

type PrivateSmsMes struct {
	Content string `json:"content"`
	User
	FriendID int `json:"friendid"`
}

type SmsReceiverMes struct {
	User
	Content string `json:"content"`
}

type SmsResMes struct {
	//200 发送成功
	//500 发送失败
	//400 用户不存在或不在线
	//300 GroupKey错误
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type AddFriendMes struct {
	UserID int    `json:"userid"`
	Friend Friend `json:"friend"`
}

type AddFriendResMes struct {
	//200 添加成功
	//500 添加失败
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type DeleteFriendMes struct {
	UserID   int `json:"userid"`
	FriendID int `json:"friendid"`
}

type DeleteFriendResMes struct {
	//200 删除成功
	//500 好友不存在
	//400 其他原因
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type LogOutMes struct {
	UserID int `json:"userid"`
}

type ListFriendsMes struct {
	UserID int `json:"userid"`
}

type ListFriendsResMes struct {
	Friends []Friend `json:"friends"`
	Error   string   `json:"error"`
}

type QueryRecordMes struct {
	//1 bySender
	//2 byReceiver
	Mod int `json:"mod"`
	ID  int `json:"id"`
}

type QueryRecordResMes struct {
	ChatRecords []ChatRecord `json:"chatrecords"`
	Error       string       `json:"error"`
}
