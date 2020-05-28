package process2

import (
	"ChatSystem/Server/dao"
	"ChatSystem/Server/utils"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
	"net"
)

type RecordProcess struct {
	//暂时不需要字段
	Conn net.Conn
}

func (rp *RecordProcess) ServerProcessQueryRecords(mes *message.Message)(err error){
	queryRecordMes := message.QueryRecordMes{}

	err = json.Unmarshal([]byte(mes.Data),&queryRecordMes)
	if err != nil{
		fmt.Println("json.Unmarshal([]byte(mes.Data),&queryRecordMes) error:",err)
		return
	}

	var chatRecords []message.ChatRecord

	switch queryRecordMes.Mod{
	case 1:
		chatRecords,err = dao.MyChatRecordDao.GetChatRecordBySenderID(queryRecordMes.ID)
		if err != nil{
			fmt.Println("dao.MyChatRecordDao.GetChatRecordBySenderID(queryRecordMes.ID) error:",err)
			return
		}
	case 2:
		chatRecords,err = dao.MyChatRecordDao.GetChatRecordByReceiverID(queryRecordMes.ID)
		if err != nil{
			fmt.Println("dao.MyChatRecordDao.GetChatRecordByReceiverID(queryRecordMes.ID) error:",err)
			return
		}
	default:
		panic("UnKnow Mod!")
	}

	resMes := message.Message{
		Type: message.QueryRecordResMesType,
	}

	queryRecordResMes := message.QueryRecordResMes{
		ChatRecords: chatRecords,
	}

	data,err := json.Marshal(queryRecordResMes)
	if err != nil{
		fmt.Println("data,err := json.Marshal(queryRecordResMes) error:",err)
		return
	}

	resMes.Data = string(data)

	data,err = json.Marshal(resMes)
	if err != nil{
		fmt.Println(" json.Marshal(resMes) error",err)
		return
	}
	tf := &utils.Transfer{
		Conn: rp.Conn,
	}
	err = tf.WritePkg(data)

	return err
}
