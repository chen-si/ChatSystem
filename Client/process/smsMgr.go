package process

import (
	"ChatSystem/common/message"
	"encoding/json"
	"fmt"
)

func outputMes(mes *message.Message) {
	//显示即可
	//1、反序列化mes
	var smsReceiverMes message.SmsReceiverMes
	err := json.Unmarshal([]byte(mes.Data), &smsReceiverMes)
	if err != nil {
		fmt.Println("outputMes error:", err)
		return
	}

	//显示信息
	info := fmt.Sprintf("用户:%s id:%d 对你说：%s", smsReceiverMes.UserName ,smsReceiverMes.UserID, smsReceiverMes.Content)
	fmt.Println(info)
}
