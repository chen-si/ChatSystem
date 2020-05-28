package main

import (
	"ChatSystem/Client/process"
	"fmt"
	"os"
)

var UserID int
var UserPWD string
var UserName string

func main() {
	//接收用户的输入的选择
	var key int
	//判断循环是否继续进行
	//var loop = true
	for {
		fmt.Println("----------------------欢迎登录聊天系统----------------------")
		fmt.Printf("\t\t\t 1 登录系统\n")
		fmt.Printf("\t\t\t 2 注册账号\n")
		fmt.Printf("\t\t\t 3 退出\n")
		fmt.Printf("\t\t\t 请选择（1-3）\n")

		_, _ = fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("开始登录")
			fmt.Println("请输入账户名：")
			_, _ = fmt.Scanf("%d\n", &UserID)
			fmt.Println("请输入密码：")
			_, _ = fmt.Scanf("%s\n", &UserPWD)
			//创建一个UserProcess实例
			up := &process.UserProcess{}
			err := up.Login(UserID, UserPWD)
			if err != nil {
				fmt.Println("up.Login(UserID,UserPWD) error:", err)
			}
		case 2:
			fmt.Println("开始注册")
			//loop = false
			fmt.Println("请输入账户ID：")
			_, _ = fmt.Scanf("%d\n", &UserID)
			fmt.Println("请输入密码：")
			_, _ = fmt.Scanf("%s\n", &UserPWD)
			fmt.Println("请输入用户名：")
			_, _ = fmt.Scanf("%s\n",&UserName)
			//接下来创建一个UserProcess实例，完成注册的请求
			up := &process.UserProcess{}
			err := up.Register(UserID,UserPWD,UserName)
			if err != nil{
				fmt.Println("up.Register(UserID,UserPWD,UserName) error:",err)
			}
		case 3:
			fmt.Println("退出成功...")
			os.Exit(0)
			//loop = false
		default:
			fmt.Println("输入有误，请重新选择功能")
		}
	}
	//fmt.Println("---------------------------------------------------------")

	//if key ==1 {
	//	fmt.Println("请输入账户名：")
	//	_, _ = fmt.Scanf("%d\n",&UserID)
	//	fmt.Println("请输入密码：")
	//	_, _ = fmt.Scanf("%s\n",&UserPWD)
	//
	//	//因为使用了分层的结构
	//	fmt.Println(err)
	//}
}
