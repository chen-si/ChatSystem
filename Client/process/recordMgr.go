package process

import (
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
)

func outputRecords(mes *message.Message) {
	defer lock.Unlock()
	var queryRecordsResMes message.QueryRecordResMes
	err := json.Unmarshal([]byte(mes.Data), &queryRecordsResMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data),&queryRecordsResMes) error:", err)
		return
	}
	fmt.Println(queryRecordsResMes)
	if queryRecordsResMes.ChatRecords != nil{
		for _,chatRecord := range queryRecordsResMes.ChatRecords{
			fmt.Println("Sender:",chatRecord.Sender,"   Receiver:",chatRecord.Receiver,"   ChatTime:",chatRecord.ChatTime)
			fmt.Println("Content:",chatRecord.Content)
			fmt.Println()
		}
	}else{
		fmt.Println("没有相关记录！")
	}
}
