package process

import (
	"ChatSystem/Client/utils"
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
)

type RecordProcess struct {
	//暂时不需要字段
}

func (rp *RecordProcess) QueryRecords(id int, mod int) {
	queryRecord := message.QueryRecordMes{
		Mod: mod,
		ID: id,
	}

	mes := message.Message{
		Type: message.QueryRecordMesType,
		Data: "",
	}

	data,err := json.Marshal(queryRecord)
	if err != nil{
		fmt.Println("json.Marshal(queryRecord) error:",err)
		return
	}
	mes.Data = string(data)

	data,err = json.Marshal(mes)
	if err != nil{
		fmt.Println("QueryRecordsBySenderId json.Marshal(mes) error:",err)
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(data)
	//fmt.Println(data)
	if err != nil {
		fmt.Println("SendGroupMes writePkg(data) error", err)
		return
	}

	return
}

