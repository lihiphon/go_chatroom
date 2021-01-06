package main

import (
	"chatroom/client/processs"
	"fmt"
	"os"
)

//定义两个变量，一个表示用户id，一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {
	//接收用户输入的选择
	var key int
	//判断是否还继续显示菜单
	var loop = true
	for loop {
		fmt.Println("--------------欢迎登陆聊天系统----------------")
		fmt.Println("\t\t  1 登陆聊天室")
		fmt.Println("\t\t  2 注册用户")
		fmt.Println("\t\t  3 退出系统")
		fmt.Println("\t\t  请选择（1-3）：")

		_, _ = fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			//fmt.Println("登陆聊天室")
			//loop = false
			//说明用户要登陆
			fmt.Println("请输入用户的id")
			_, _ = fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			_, _ = fmt.Scanf("%s\n", &userPwd)
			//创建一个userProcess实例
			up := &processs.UserProcess{}
			_ = up.Login(userId, userPwd)
		case 2:
			//fmt.Println("注册用户")
			//loop = false
			fmt.Println("请输入用户的id")
			_, _ = fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			_, _ = fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户的名字")
			_, _ = fmt.Scanf("%s\n", &userName)
			up := &processs.UserProcess{}
			err := up.Register(userId, userPwd, userName)
			if err != nil {
				fmt.Println("用户注册失败 err=", err)
			}
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你输入有误，请重新输入")
		}

	}
	//增加用户的输入，显示新的提示信息
	//	if key == 1 {
	//		//说明用户要登陆
	//		fmt.Println("请输入用户的id")
	//		_, _ = fmt.Scanf("%d\n", &userId)
	//		fmt.Println("请输入用户的密码")
	//		_, _ = fmt.Scanf("%s\n", &userPWD)
	//		//先把登录的函数写到另外一个文件，login.go
	//		//err := processs.UserProcess.Login(userId,userPWD)
	//		//if err != nil {
	//		//	fmt.Println("登录失败")
	//		//
	//		//} else {
	//		//	fmt.Println("登录成功")
	//		//}
	//		//创建一个userProcess实例
	//		up:=&processs.UserProcess{}
	//		_ = up.Login(userId, userPWD)
	//	} else if key == 2 {
	//		fmt.Println("进行用户注册的逻辑")
	//	}
}
