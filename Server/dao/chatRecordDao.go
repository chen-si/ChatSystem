package dao

import (
	"ChatSystem/common/message"
	"database/sql"
	"fmt"
)

var (
	MyChatRecordDao ChatRecordDao
)

type ChatRecordDao struct{
	db *sql.DB
}

func NewChatRecordDao(db *sql.DB)ChatRecordDao{
	return ChatRecordDao{
		db : db,
	}
}

func (chatRecordDao *ChatRecordDao)GetChatRecordBySenderID(senderId int)(chatRecords []message.ChatRecord,err error){
	sqlStr := "select sender_id,receiver_id,record,chat_time from chat_record where sender_id = ?"
	rows,err := chatRecordDao.db.Query(sqlStr,senderId)
	if err != nil{
		fmt.Println("chatRecordDao.db.Query(sqlStr,senderId) error:",err)
		return
	}

	for rows.Next(){
		chatRecord := message.ChatRecord{}
		err = rows.Scan(&chatRecord.Sender,&chatRecord.Receiver,&chatRecord.Content,&chatRecord.ChatTime)
		if err != nil{
			fmt.Println("rows.Scan(&chatRecord.Sender,&chatRecord.Receiver,&chatRecord.Content,&chatRecord.ChatTime) error:",err)
		}
		chatRecords = append(chatRecords,chatRecord)
	}

	return
}

func (chatRecordDao *ChatRecordDao)GetChatRecordByReceiverID(receiverId int)(chatRecords []message.ChatRecord,err error){
	sqlStr := "select sender_id,receiver_id,record,chat_time from chat_record where receiver_id = ?"
	rows,err := chatRecordDao.db.Query(sqlStr,receiverId)
	if err != nil{
		fmt.Println("chatRecordDao.db.Query(sqlStr,receiverId) error:",err)
		return
	}

	for rows.Next(){
		chatRecord := message.ChatRecord{}
		err = rows.Scan(&chatRecord.Sender,&chatRecord.Receiver,&chatRecord.Content,&chatRecord.ChatTime)
		if err != nil{
			fmt.Println("rows.Scan(&chatRecord.Sender,&chatRecord.Receiver,&chatRecord.Content,&chatRecord.ChatTime) error:",err)
		}
		chatRecords = append(chatRecords,chatRecord)
	}

	return
}

func (chatRecordDao *ChatRecordDao)InsertChatRecord(chatRecord *message.ChatRecord)(err error){
	sqlStr := "insert into chat_record(sender_id,receiver_id,record) values(?,?,?)"
	_,err = chatRecordDao.db.Exec(sqlStr,chatRecord.Sender,chatRecord.Receiver,chatRecord.Content)
	if err != nil{
		fmt.Println("chatRecordDao.db.Exec(sqlStr,chatRecord.Sender,chatRecord.Reciever,chatRecord.Content) error:",err)
	}
	return
}