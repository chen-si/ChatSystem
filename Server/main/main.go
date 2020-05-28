package main

import (
	"ChatSystem/Server/dao"
	process2 "ChatSystem/Server/process"
	"ChatSystem/Server/utils"
	"fmt"
	"net"
)

//func readPkg(conn net.Conn) (mes message.Message ,err error){
//	//defer conn.Close()
//	buf := make([]byte,8096)
//	fmt.Println("读取客户端发送的消息....")
//	_,err = conn.Read(buf[:4])
//	if err != nil{
//		//err = errors.New("conn.Read(buf[:4]) error")
//		return
//	}
//	//将读取到的buf[:4]转换成uint32类型
//	var PkgLen uint32
//	PkgLen = binary.BigEndian.Uint32(buf[:4])
//
//	//根据PkgLen读取数据和话剧
//	n,err := conn.Read(buf[:PkgLen])
//	if err != nil || n != int(PkgLen){
//		err = errors.New("conn.Read(buf[:PkgLen]) error")
//		return
//	}
//
//	//将读取到的数据反序列化
//	err = json.Unmarshal(buf[:PkgLen],&mes)
//	if err != nil{
//		err = errors.New("json.Unmarshal(buf[:PkgLen],&mes) error")
//		return
//	}
//
//	return
//}

//func writePkg(conn net.Conn,data []byte)(err error){
//	//先发送一个长度给对方
//	var pkglen uint32
//	pkglen = uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[0:4],pkglen)
//	//发送长度
//	n,err :=conn.Write(buf[0:4])
//	if n != 4 || err != nil{
//		fmt.Println("conn.Write(buf) error:",err)
//		return
//	}
//	//发送消息本身
//	//发送消息本身
//	n,err = conn.Write(data)
//	if n != 4 || err != nil{
//		fmt.Println("conn.Write(data) error:",err)
//		return
//	}
//	return
//}

//编写一个serverProcessLogin函数，专门处理登录请求
//func serverProcessLogin(conn net.Conn,mes *message.Message)(err error){
//	//处理登录
//	//先从mes.Data反序列化成LoginMes
//	var loginMes message.LoginMes
//	err = json.Unmarshal([]byte(mes.Data),&loginMes)
//	if err !=nil{
//		fmt.Println("json.Unmarshal([]byte(mes.Data),loginMes) error:",err)
//		return
//	}
//
//	var resMes message.Message
//	resMes.Type =message.LoginResMesType
//	var loginResMes message.LoginResMes
//
//	//如果用户名为100 密码为123456表示登录成功
//	if loginMes.UserID == 100 && loginMes.UserPWD == "123456"{
//		//登陆成功
//		loginResMes.Code=200 //200表示登录成功
//
//	}else{
//		//登录失败
//		loginResMes.Code=500 //500表示用户未注册
//		loginResMes.Error="该用户不存在，请注册再使用"
//	}
//	//反序列化loginResMes并复制给resMes.Data
//	data,err := json.Marshal(loginResMes)
//	if err != nil{
//		fmt.Println("json.Marshal(loginResMes) error",err)
//		return
//	}
//	resMes.Data = string(data)
//
//	//反序列化resMes
//	data,err = json.Marshal(resMes)
//	if err != nil{
//		fmt.Println("json.Marshal(resMes) error",err)
//		return
//	}
//	//发送data 将发送函数封装成writePkg（）
//	err = writePkg(conn,data)
//
//	return
//}

////编写一个serverProcessMes函数处理不同的消息
////根据消息类型的而不同调用不同的函数
//func serverProcessMes(conn net.Conn,mes *message.Message)(err error){
//	switch mes.Type{
//		case message.LoginMesType :
//			//处理用户登录
//			err = serverProcessLogin(conn,mes)
//		case message.RegisterMesType :
//			//处理注册
//		default:
//			fmt.Println("不支持的消息类型...")
//
//
//	}
//	return
//}

func process(conn net.Conn) {
	defer conn.Close()

	//这里调用总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.Process3()
	if err != nil {
		fmt.Println("processor.process3() error:", err)
		return
	}
}

//这里我们写一个函数完成对userdao的初始实例化
func initUserDao() {
	dao.MyUserDao = dao.NewUserDao(utils.Db)
}

func initChatRecordDao(){
	dao.MyChatRecordDao = dao.NewChatRecordDao(utils.Db)
}


func main() {
	//当服务器启动时，我们就去初始化redis的链接池
	initUserDao()
	initChatRecordDao()
	fmt.Println("服务器开始监听9999端口......")
	fmt.Println("请设定GroupKey:")
	_, _ = fmt.Scanf("%s\n", &process2.GroupKey)
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("Linten error:", err)
		return
	}
	defer listener.Close()
	for {
		fmt.Println("等待客户端连接服务器......")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener Accept error:", err)
		}
		//goroutine处理不同服务器的请求
		go process(conn)
	}
}
