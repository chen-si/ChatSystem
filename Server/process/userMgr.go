package process2

import "fmt"

//因为UserMgr在服务器端有且仅有一个
//在很多地方都有用到
//所以定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//增加在线用户列表
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserID] = up
}

//删除
func (this *UserMgr) DelOnlineUser(UserID int) {
	delete(this.onlineUsers, UserID)
}

//返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据ID返回对应的值
func (this *UserMgr) GetOnlineUserByID(UserID int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[UserID]
	if !ok { //说明你要找的用户当前不在线
		err = fmt.Errorf("用户 %d 不在线或不存在", UserID)
		return
	}
	return
}
