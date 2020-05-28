package dao

import (
	"ChatSystem/Server/model"
	"ChatSystem/common/message"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

//我们在服务器启动时就初始化一个UserDao实例
//把他做成全局变量，在需要使用时直接使用即可
var (
	MyUserDao *UserDao
)

//先定义一个UserDao 结构体 Dao:data access object
//完成对User 结构体的各种操作

type UserDao struct {
	db *sql.DB
}

//使用工厂模式，创建UserDao的实例
func NewUserDao(db *sql.DB) (userdao *UserDao) {
	userdao = &UserDao{
		db:db,
	}
	return
}

func (userDao *UserDao) getUserById(id int) (user *message.User, err error) {
	sqlStr := "select user_id,user_pwd,user_name,ifnull(friends_info,\"没有好友信息\") from users where user_id = ?"
	row :=  userDao.db.QueryRow(sqlStr,strconv.Itoa(id))
	if row == nil{
		return nil,model.ERROR_USER_NOTEXISTS
	}
	//初始化user
	user = &message.User{}
	friendsInfo := ""
	err = row.Scan(&user.UserID,&user.UserPWD,&user.UserName,&friendsInfo)
	if err != nil{
		return nil,err
	}
	if friendsInfo == "没有好友信息"{
		user.FriendNotExist = true
	}else{
		user.FriendNotExist = false
		err = json.Unmarshal([]byte(friendsInfo),&user.Friends)
		if err != nil{
			return nil,err
		}
	}
	return
}

//完成登录的校验
//1、Login 完成用户的验证
//2、如果id和pwd都正确，则返回一个user实例
//3、如果id或pwd错误，返回错误提示信息
func (userDao *UserDao) Login(UserID int, UserPWD string) (user *message.User, err error) {
	//先从UserDao的链接池中取出一个链接
	user, err = userDao.getUserById(UserID)
	if err != nil {
		return
	}

	//这时候验证密码的正确性
	if UserPWD != user.UserPWD {
		err = model.ERROR_USER_PWD
		return
	}
	return
}

func (userDao *UserDao) Register(user *message.User) (err error) {
	_, err = userDao.getUserById(user.UserID)
	if err == nil {
		err = model.ERROR_USER_EXISTS
		return
	}

	//说明id在redis中不存在 可以完成注册
	sqlStr := "insert into users (user_id,user_name,user_pwd) value(?,?,?)"
	_,err = userDao.db.Exec(sqlStr,user.UserID,user.UserName,user.UserPWD)
	if err != nil{
		fmt.Println("db.Exec(sqlStr,user.UserID,user.UserName,user.UserPWD) error:",err)
	}

	return
}

//好友管理模块
func (userDao *UserDao) UpdateOneFriend(id int,friend *message.Friend)(err error) {
	user, err := userDao.getUserById(id)
	if err != nil {
		return model.ERROR_USER_NOTEXISTS
	}
	if user.FriendNotExist {
		user.Friends = []message.Friend{}
	}
	user.Friends = append(user.Friends, *friend)

	err = userDao.UpdateFriends(id,user)
	return
}

func (userDao *UserDao) UpdateFriends(id int, user *message.User)(err error){
	sqlStr := "update users set friends_info = ? where user_id = ?"
	if user.FriendNotExist{
		_, err = userDao.db.Exec(sqlStr, "没有好友信息", id)
		if err != nil {
			fmt.Println("userDao.db.Exec(sqlStr,string(friends_info),id) error:", err)
		}
		return
	}
	friends_info, err := json.Marshal(user.Friends)
	if err != nil {
		fmt.Println("json.Marshal(user.Friends) error:", err)
		return
	}

	_, err = userDao.db.Exec(sqlStr, string(friends_info), id)
	if err != nil {
		fmt.Println("userDao.db.Exec(sqlStr,string(friends_info),id) error:", err)
	}
	return
}

func (userDao *UserDao)DeleteFriend(id int,friendId int)(err error){
	user, err := userDao.getUserById(id)
	if err != nil {
		return model.ERROR_USER_NOTEXISTS
	}
	if user.FriendNotExist{
		return model.ERROR_FRIEND_NOT_EXIST
	}
	for i,friend := range user.Friends{
		if friend.FriendID == friendId{
			if i == len(user.Friends) - 1{
				user.Friends = user.Friends[:i]
			}else{
				user.Friends = append(user.Friends[:i],user.Friends[i+1:]...)
			}
		}
	}
	if len(user.Friends) == 0{
		user.FriendNotExist = true
	}
	err = userDao.UpdateFriends(id,user)
	return
}

func (userDao *UserDao) ListFriends(id int)(friends []message.Friend,err error){
	user, err := userDao.getUserById(id)
	if err != nil {
		err = model.ERROR_USER_NOTEXISTS
		return
	}
	if user.FriendNotExist{
		return nil,model.ERROR_FRIEND_NOT_EXIST
	}else{
		return user.Friends,nil
	}
}