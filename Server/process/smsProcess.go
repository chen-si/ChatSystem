package process2

import (
	"ChatSystem/Client/utils"
	"ChatSystem/Server/dao"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
	//暂时不需要字段
	Conn net.Conn
}

var GroupKey string

//转发消息的方法
func (smsProcess *SmsProcess) SendGroupMes(mes *message.Message) (err error) {
	//遍历服务器端的onlineUsers map[int]*userprocess
	//将消息转发出去
	var smsResMes message.SmsResMes
	var groupSmsMes message.GroupSmsMes
	err = json.Unmarshal([]byte(mes.Data), &groupSmsMes)
	if err != nil {
		smsResMes = message.SmsResMes{
			Code:  500,
			Error: err.Error(),
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}
	if groupSmsMes.GroupKey == GroupKey {
		smsRecieverMes := message.SmsReceiverMes{
			Content: groupSmsMes.Content,
			User: message.User{
				UserID:     groupSmsMes.UserID,
				UserName:   groupSmsMes.UserName,
				UserStatus: groupSmsMes.UserStatus,
			},
		}
		data, err := json.Marshal(smsRecieverMes)
		if err != nil {
			smsResMes = message.SmsResMes{
				Code:  500,
				Error: err.Error(),
			}
			smsProcess.ReturnSmsResMes(smsResMes)
			return err
		}
		resMes := message.Message{
			Type: message.SmsReceiverMesType,
		}
		resMes.Data = string(data)
		data, err = json.Marshal(resMes)
		if err != nil {
			smsResMes = message.SmsResMes{
				Code:  500,
				Error: err.Error(),
			}
			smsProcess.ReturnSmsResMes(smsResMes)
			return err
		}
		for id, up := range userMgr.onlineUsers {
			//这里我们需要过滤掉自己
			if id == groupSmsMes.UserID {
				continue
			}
			chatRecord := message.ChatRecord{
				Sender:   groupSmsMes.UserID,
				Receiver: id,
				Content:  groupSmsMes.Content,
			}
			err = dao.MyChatRecordDao.InsertChatRecord(&chatRecord)
			if err != nil {
				smsResMes = message.SmsResMes{
					Code:  500,
					Error: err.Error(),
				}
				smsProcess.ReturnSmsResMes(smsResMes)
				return err
			}
			smsProcess.SendMesToOnlineUser(data, up.Conn)
		}
		//return err
	} else {
		//返回错误信息
		smsResMes = message.SmsResMes{
			Code:  300,
			Error: "无效的GroupKey",
		}
		smsProcess.ReturnSmsResMes(smsResMes)
	}
	smsResMes = message.SmsResMes{
		Code:  200,
		Error: "",
	}
	smsProcess.ReturnSmsResMes(smsResMes)
	return err
}

func (smsProcess *SmsProcess) SendPrivateMes(mes *message.Message) (err error) {
	privateSmsMes := message.PrivateSmsMes{}
	var smsResMes message.SmsResMes

	err = json.Unmarshal([]byte(mes.Data), &privateSmsMes)
	if err != nil {
		smsResMes = message.SmsResMes{
			Code:  500,
			Error: err.Error(),
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}
	up, ok := userMgr.onlineUsers[privateSmsMes.FriendID]
	if !ok {
		smsResMes = message.SmsResMes{
			Code:  400,
			Error: "用户不存在或者不在线",
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}

	smsRecieverMes := message.SmsReceiverMes{
		Content: privateSmsMes.Content,
		User: message.User{
			UserID:     privateSmsMes.UserID,
			UserName:   privateSmsMes.UserName,
			UserStatus: privateSmsMes.UserStatus,
		},
	}
	data, err := json.Marshal(smsRecieverMes)
	if err != nil {
		smsResMes = message.SmsResMes{
			Code:  500,
			Error: err.Error(),
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}
	resMes := message.Message{
		Type: message.SmsReceiverMesType,
	}
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		smsResMes = message.SmsResMes{
			Code:  500,
			Error: err.Error(),
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}

	chatRecord := message.ChatRecord{
		Sender:   privateSmsMes.UserID,
		Receiver: privateSmsMes.FriendID,
		Content:  privateSmsMes.Content,
	}
	err = dao.MyChatRecordDao.InsertChatRecord(&chatRecord)
	if err != nil {
		smsResMes = message.SmsResMes{
			Code:  500,
			Error: err.Error(),
		}
		smsProcess.ReturnSmsResMes(smsResMes)
		return
	}
	smsProcess.SendMesToOnlineUser([]byte(data), up.Conn)
	smsResMes = message.SmsResMes{
		Code:  200,
		Error: "",
	}
	smsProcess.ReturnSmsResMes(smsResMes)
	return
}

func (smsProcess *SmsProcess) ReturnSmsResMes(smsResMes message.SmsResMes) {
	data, err := json.Marshal(smsResMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(smsResMes) error:", err)
		return
	}
	resMes := message.Message{
		Type: message.SmsResMesType,
	}
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal(resMes) error:", err)
		return
	}

	tf := &utils.Transfer{
		Conn: smsProcess.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendMesToEachOnlineUser error:", err)
		return
	}
}

func (smsProcess *SmsProcess) SendMesToOnlineUser(data []byte, conn net.Conn) {
	//创建一个transfer实例 发送data
	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendMesToEachOnlineUser error:", err)
		return
	}
}
