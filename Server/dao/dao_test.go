package dao

import (
	"ChatSystem/Server/utils"
	"ChatSystem/common/message"
	"fmt"
	"testing"
)

func TestUserDao(t *testing.T){
	//t.Run("测试添加一个好友", testUpdateOneFriend)
	//t.Run("测试删除好友",testDeleteFriend)
	//t.Run("测试显示好友列表",testListFriends)
	//t.Run("测试保存聊天记录",testInsertChatRecord)
	//t.Run("测试通过发送者查询",testGetChatRecordBySenderID)
	t.Run("测试通过接收者查询",testGetChatRecordByReceiverID)
}

func testUpdateOneFriend(t *testing.T) {
	userDao := NewUserDao(utils.Db)
	friend := &message.Friend{
		FriendID: 159,
		Group:    "1",
	}
	err := userDao.UpdateOneFriend(123,friend)
	fmt.Println(err)
}

func testDeleteFriend(t *testing.T){
	userDao := NewUserDao(utils.Db)
	err := userDao.DeleteFriend(123,159)
	fmt.Println(err)
}

func testListFriends(t *testing.T){
	userDao := NewUserDao(utils.Db)
	friends,err := userDao.ListFriends(123)
	fmt.Println(err)
	fmt.Println(friends)
}

func testInsertChatRecord(t *testing.T) {
	chatRecord := message.ChatRecord{
		Sender:   123,
		Receiver: 159,
		Content:  "hello world!",
	}
	chatRecordDao := NewChatRecordDao(utils.Db)

	err := chatRecordDao.InsertChatRecord(&chatRecord)
	fmt.Println(err)
}

func testGetChatRecordBySenderID(t *testing.T){
	chatRecordDao := NewChatRecordDao(utils.Db)

	chatRecords,err := chatRecordDao.GetChatRecordBySenderID(123)
	fmt.Println(err)
	fmt.Println(chatRecords)
}

func testGetChatRecordByReceiverID(t *testing.T){
	chatRecordDao := NewChatRecordDao(utils.Db)

	chatRecords,err := chatRecordDao.GetChatRecordByReceiverID(159)
	fmt.Println(err)
	fmt.Println(chatRecords)
}