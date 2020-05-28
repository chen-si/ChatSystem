package message

type ChatRecord struct{
	Sender 		int		`json:"sender"`
	Receiver 	int		`json:"reciever"`
	Content 	string 	`json:"content"`
	ChatTime	string  `json:"chattime"`
}

type FileRecord struct{
	Sender 		int		`json:"sender"`
	Receiver 	int		`json:"reciever"`
	FileName 	string 	`json:"filename"`
	MD5 		string 	`json:"md5"`
}