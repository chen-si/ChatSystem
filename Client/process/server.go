package process

import (
	"ChatSystem/Client/utils"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

var lock sync.Mutex

//ShowMenu 显示登录成功后的界面
func ShowMenu(name string) {
	lock.Lock()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Printf("---------恭喜%s登陆成功---------\n", name)
	fmt.Println("---------1、好友管理-------------")
	fmt.Println("---------2、发送消息-------------")
	fmt.Println("---------3、消息记录查询----------")
	fmt.Println("---------4、在线用户列表----------")
	fmt.Println("---------5、退出系统-------------")
	fmt.Println("请选择（1-5）：")
	var key int

	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		FriendManagerMenu()
	case 2:
		SmsMenu()
	case 3:
		//fmt.Println("消息记录查询")
		RecordMenu()
	case 4:
		outputOnlineUsers()
	case 5:
		fmt.Println("退出系统")
		up := &UserProcess{}
		_ = up.LogOut()
		os.Exit(0)
	default:
		fmt.Println("")

	}

}

func RecordMenu() {
	recordProcess := RecordProcess{}
	var key int
	fmt.Println("消息记录查询：")
	fmt.Println("	1、查看自己发送的消息记录")
	fmt.Println("	2、查看自己接收的消息记录")
	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		recordProcess.QueryRecords(CurUser.UserID,1)
	case 2:
		recordProcess.QueryRecords(CurUser.UserID,2)
	default:
		fmt.Println("未知选项，返回主菜单")
		defer lock.Unlock()
		return
	}
}

func SmsMenu() {
	var content string
	//因为我们总会使用到SmsProcess 因此将其定义在switch外部
	smsProcess := &SmsProcess{}
	var key int
	var friendId int
	var groupKey string
	fmt.Println("发送消息：")
	fmt.Println("	1、发送私聊消息")
	fmt.Println("	2、发送广播消息")
	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("请输入好友ID：")
		_, _ = fmt.Scanf("%d\n", &friendId)
		fmt.Println("请输入私聊发送的消息内容：")
		content = ScanLine()
		fmt.Println()
		fmt.Println(content)
		_ = smsProcess.SendPrivateMes(friendId, content)
	case 2:
		fmt.Println("请输入权限密码：")
		_, _ = fmt.Scanf("%s\n", &groupKey)
		fmt.Println("请输入广播发送的消息内容：")
		content = ScanLine()
		_ = smsProcess.SendGroupMes(groupKey, content)
	default:
		fmt.Println("未知选项，返回主菜单")
		defer lock.Unlock()
		return
	}
}

func FriendManagerMenu() {
	fmt.Println("好友管理：")
	fmt.Println("	1、添加好友")
	fmt.Println("	2、删除好友")
	fmt.Println("	3、好友列表")
	var key int
	var friendId int
	var friendGroup string
	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("添加好友")
		fmt.Println("请输入好友ID：")
		_, _ = fmt.Scanf("%d\n", &friendId)
		fmt.Println("请输入好友分组：")
		_, _ = fmt.Scanf("%s\n", &friendGroup)
		up := UserProcess{}
		up.AddFriend(friendId, friendGroup)
	case 2:
		fmt.Println("删除好友")
		fmt.Println("请输入待删除的好友ID：")
		_, _ = fmt.Scanf("%d\n", &friendId)
		up := UserProcess{}
		up.DeleteFriend(friendId)
	case 3:
		fmt.Println("好友列表")
		up := UserProcess{}
		_ = up.ListFriends()
	default:
		fmt.Println("未知选项，返回主菜单")
		defer lock.Unlock()
		return
	}
}

//ServerProcessMes 和服务器段保持通讯
func ServerProcessMes(Conn net.Conn) {
	//创建一个transfer实例 循环读取消息
	tf := &utils.Transfer{
		Conn: Conn,
	}
	for {
		//fmt.Printf("客户端：%s正在等待读取服务器发送的消息...\n", Conn.LocalAddr().String())
		mes, err := tf.ReadPkg()
		if err != nil {
			//fmt.Println("ServerProcessMes tf.ReadPkg() error:", err)
			return
		}

		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人状态改变了
			//1、取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2、更新用户状态
			UpdateUserStatus(&notifyUserStatusMes)
		case message.SmsReceiverMesType:
			outputMes(&mes)
		case message.SmsResMesType:
			sms := SmsProcess{}
			sms.ProcessSmsResMes(&mes)
		case message.AddFriendResMesType:
			up := UserProcess{}
			err = up.ProcessAddFriendResMes(&mes)
		case message.DeleteFriendResMesType:
			up := UserProcess{}
			err = up.ProcessDeleteFriendResMes(&mes)
		case message.ListFriendsResMesType:
			up := UserProcess{}
			err = up.ProcessListFriendsResMes(&mes)
		case message.QueryRecordResMesType:
			outputRecords(&mes)
		default:
			fmt.Println("服务器端返回了一个未知的消息类型")
		}

	}
}


func ScanLine() string {
	var c byte
	var err error
	var b []byte
	for ; err == nil; {
		_, err = fmt.Scanf("%c", &c)

		if c != '\n' {
			b = append(b, c)
		} else {
			break;
		}
	}

	return string(b)
}