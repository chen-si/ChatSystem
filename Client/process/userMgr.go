package process

import (
	"ChatSystem/Client/model"
	"ChatSystem/common/message"
	"fmt"
)

//客户端药品维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //我们在用户登录成功后完成对CurUser的初始化

//在客户端显示当前在线的用户
func outputOnlineUsers() {
	defer lock.Unlock()
	//遍历map
	fmt.Println("当前用户在线列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户ID = ", id)
	}
}

//编写一个方法处理返回的NotifyUserStatusMes
func UpdateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if notifyUserStatusMes.Status == message.UserOnline{
		if !ok { //原来不存在
			user = &message.User{
				UserID:     notifyUserStatusMes.UserID,
				UserStatus: notifyUserStatusMes.Status,
			}
		}
		user.UserStatus = notifyUserStatusMes.Status
		onlineUsers[notifyUserStatusMes.UserID] = user
	}else{
		//从在线用户列表删除logout的用户
		delete(onlineUsers,notifyUserStatusMes.UserID)
	}
}
