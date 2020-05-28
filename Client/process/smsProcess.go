package process

import (
	"ChatSystem/Client/utils"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
)

//SmsProcess ...
type SmsProcess struct {
}

//SendGroupMes 发送群聊消息
func (sms *SmsProcess) SendGroupMes(groupKey string,content string) (err error) {
	//1、创建一个Message
	var mes message.Message
	mes.Type = message.GroupSmsMesType

	//2、创建一个smsMes实例
	var groupSmsMes message.GroupSmsMes
	groupSmsMes.Content = content //内容
	groupSmsMes.GroupKey = groupKey
	groupSmsMes.UserID = CurUser.UserID
	groupSmsMes.UserStatus = CurUser.UserStatus
	groupSmsMes.UserName = CurUser.UserName


	//3、 序列化smsMes
	data, err := json.Marshal(groupSmsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(groupSmsMes) error:", err)
		return
	}
	mes.Data = string(data)

	//4、把 mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(mes) error:", err)
		return
	}

	//5、发送数据data
	//5.1 、先把data的长度发送给服务器
	//先获取到data的长度，然后转化为一个表示长度的切片
	//先创建一个Transfer实例

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes writePkg(data) error", err)
		return
	}

	return
}

func (sms *SmsProcess) SendPrivateMes(friendId int,content string)(err error){
	//1、创建一个Message
	var mes message.Message
	mes.Type = message.PrivateSmsMesType

	//2、创建一个smsMes实例
	var privateSmsMes message.PrivateSmsMes
	privateSmsMes.Content = content
	privateSmsMes.FriendID = friendId
	privateSmsMes.UserID = CurUser.UserID
	privateSmsMes.UserStatus = CurUser.UserStatus
	privateSmsMes.UserName = CurUser.UserName

	//3、 序列化smsMes
	data, err := json.Marshal(privateSmsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(privateSmsMes) error:", err)
		return
	}
	mes.Data = string(data)

	//4、把 mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(mes) error:", err)
		return
	}

	//5、发送数据data
	//5.1 、先把data的长度发送给服务器
	//先获取到data的长度，然后转化为一个表示长度的切片
	//先创建一个Transfer实例
	//fmt.Println(string(data))
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes writePkg(data) error", err)
		return
	}

	return
}

func (sms *SmsProcess) ProcessSmsResMes(mes *message.Message){
	defer lock.Unlock()
	smsResMes := message.SmsResMes{}
	err := json.Unmarshal([]byte(mes.Data),&smsResMes)
	if err != nil{
		fmt.Println("json.Unmarshal([]byte(mes.Data),&smsResMes) error:",err)
		return
	}
	switch smsResMes.Code{
	case 200:
		fmt.Println("消息发送成功")
	case 300:
		fmt.Println("GroupKey错误")
	case 400:
		fmt.Println("好友不在线或者不存在")
	case 500:
		fmt.Println("发送失败：",smsResMes.Error)
	default:
		fmt.Println(smsResMes.Error)
	}
}