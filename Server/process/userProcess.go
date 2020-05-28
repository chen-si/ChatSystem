package process2

import (
	"ChatSystem/Server/dao"
	"ChatSystem/Server/model"
	"ChatSystem/Server/utils"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段 表示是哪一个用户的conn
	UserID int
}

//这里我们编写通知用户上线的方法
//userid 通知其他在线用户 我上线了
func (up *UserProcess) NotifyOtherOnlineUser(UserID int) {
	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == UserID {
			continue
		}
		//开始通知【单独的一个函数】
		up.NotifyMeStatus(UserID, message.UserOnline)
	}
}

func (up *UserProcess) NotifyMeStatus(UserID int, status int) {
	//组装我们的消息NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = UserID
	notifyUserStatusMes.Status = status

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal(notifyUserStatusMes) error:", err)
		return
	}

	//把序列化的notifyUserStatusMes赋值给 mes.Data
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal(mes) error:", err)
		return
	}

	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline tf.WritePkg(data) error:", err)
		return
	}

}

//编写一个serverProcessLogin函数，专门处理登录请求
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//处理登录
	//先从mes.Data反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),loginMes) error:", err)
		return err
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	//到redis数据库验证用户
	//使用model.MyUserDao 到mysql去验证
	user, err := dao.MyUserDao.Login(loginMes.UserID, loginMes.UserPWD)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}

	} else {
		loginResMes.Code = 200
		loginResMes.UserName = user.UserName
		//将登录成功的用户id赋给this
		up.UserID = loginMes.UserID
		//把登录成功的用户放入userMgr中
		userMgr.AddOnlineUser(up)
		//通知其他在线的用户 我上线了
		up.NotifyOtherOnlineUser(loginMes.UserID)

		//遍历userid
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersID = append(loginResMes.UsersID, id)
		}
		fmt.Println(user, "登陆成功...")
	}

	// //如果用户名为100 密码为123456表示登录成功
	// if loginMes.UserID == 100 && loginMes.UserPWD == "123456"{
	// 	//登陆成功
	// 	loginResMes.Code=200 //200表示登录成功

	// }else{
	// 	//登录失败
	// 	loginResMes.Code=500 //500表示用户未注册
	// 	loginResMes.Error="该用户不存在，请注册再使用"
	// }
	//反序列化loginResMes并复制给resMes.Data
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) error", err)
		return err
	}
	resMes.Data = string(data)

	//反序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) error", err)
		return err
	}
	//发送data 将发送函数封装成writePkg（）
	//因为使用了分层模式（mvs）我们先创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)

	return err
}

func (up *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&registerMes) error:", err)
		return err
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//到redis数据库注册用户
	//使用model.MyUserDao 到redis去注册
	err = dao.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
		fmt.Println("注册成功...")
	}

	//反序列化registerResMes并复制给resMes.Data
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(registerResMes) error", err)
		return err
	}
	resMes.Data = string(data)

	//反序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) error", err)
		return err
	}
	//发送data 将发送函数封装成writePkg（）
	//因为使用了分层模式（mvs）我们先创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (up *UserProcess) ServerProcessAddFriend(mes *message.Message) (err error) {
	addFriendMes := message.AddFriendMes{}
	err = json.Unmarshal([]byte(mes.Data), &addFriendMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&addFriendMes) error:", err)
		return err
	}

	addFriendResMes := message.AddFriendResMes{}

	err = dao.MyUserDao.UpdateOneFriend(addFriendMes.UserID, &addFriendMes.Friend)
	if err != nil {
		addFriendResMes.Code = 500
		addFriendResMes.Error = err.Error()
	} else {
		addFriendResMes.Code = 200
		fmt.Println("添加好友成功。。。")
	}
	resMes := message.Message{
		Type: message.AddFriendResMesType,
	}
	data, err := json.Marshal(addFriendResMes)
	if err != nil {
		fmt.Println("json.Marshal(json.Marshal(addFriendResMes)) error", err)
		return err
	}
	resMes.Data = string(data)

	//反序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) error", err)
		return err
	}
	//发送data 将发送函数封装成writePkg（）
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (up *UserProcess) ServerProcessDeleteFriend(mes *message.Message) (err error) {
	deleteFriendMes := message.DeleteFriendMes{}
	err = json.Unmarshal([]byte(mes.Data), &deleteFriendMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&addFriendMes) error:", err)
		return err
	}

	deleteFriendResMes := message.DeleteFriendResMes{}

	err = dao.MyUserDao.DeleteFriend(deleteFriendMes.UserID, deleteFriendMes.FriendID)
	if err != nil {
		if err == model.ERROR_FRIEND_NOT_EXIST {
			deleteFriendResMes.Code = 500
		} else {
			deleteFriendResMes.Code = 400
		}
		deleteFriendResMes.Error = err.Error()
	} else {
		deleteFriendResMes.Code = 200
		fmt.Println("删除好友成功。。。")
	}
	resMes := message.Message{
		Type: message.DeleteFriendResMesType,
	}
	data, err := json.Marshal(deleteFriendResMes)
	if err != nil {
		fmt.Println("json.Marshal(json.Marshal(deleteFriendResMes)) error", err)
		return err
	}
	resMes.Data = string(data)

	//反序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) error", err)
		return err
	}
	//发送data 将发送函数封装成writePkg（）
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (up *UserProcess) ServerProcessLogOut(mes *message.Message) (err error) {
	logOutMes := message.LogOutMes{}
	err = json.Unmarshal([]byte(mes.Data), &logOutMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&logOutMes) error:", err)
		return
	}
	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == logOutMes.UserID {
			continue
		}
		//开始通知
		up.NotifyMeStatus(logOutMes.UserID, message.UserOffline)
	}
	userMgr.DelOnlineUser(logOutMes.UserID)
	return
}

func (up *UserProcess) ServerProcessListFriends(mes *message.Message) (err error) {
	listFriendsMes := message.ListFriendsMes{}
	err = json.Unmarshal([]byte(mes.Data), &listFriendsMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&listFriendsMes) error:", err)
		return
	}
	listFriendsResMes := message.ListFriendsResMes{}
	friends, err := dao.MyUserDao.ListFriends(listFriendsMes.UserID)
	if err != nil {
		if err == model.ERROR_FRIEND_NOT_EXIST {
			listFriendsResMes.Error = err.Error()
		}
	} else {
		listFriendsResMes.Friends = friends
		listFriendsResMes.Error = ""
	}

	resMes := message.Message{
		Type: message.ListFriendsResMesType,
	}
	data, err := json.Marshal(listFriendsResMes)
	if err != nil {
		fmt.Println("json.Marshal(json.Marshal(listFriendsResMes)) error", err)
		return err
	}
	resMes.Data = string(data)

	//反序列化resMes
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) error", err)
		return err
	}
	//发送data 将发送函数封装成writePkg（）
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)

	return
}
