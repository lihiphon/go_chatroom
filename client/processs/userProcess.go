/**
 * @Author haifengli
 * @Date 9:00 下午 2021/1/1
 * @Description
1.处理和用户相关的业务
2.登录
3.注册
等等
 **/
package processs

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要字段。。。
}

//用户登录函数
func (this UserProcess) Login(userId int, userPwd string) (err error) {
	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net.dial fail err=", err)
		return
	}
	//延时关闭
	defer func() { conn.Close() }()
	//2.准备通过conn发送的消息
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.marshal fail err=", err)
		return
	}
	//5.把data赋给mes.data字段
	mes.Data = string(data)
	//6.将mes进行序列化. 返回值data就是将要发送的数据
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//7.先把data的长度发给服务器 （将data 的长度->转成一个表示长度的切片）
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn Write fail err=", err)
		return
	}
	fmt.Printf("客户端发送消息长度=%d,内容=%s\n", len(data), string(data))
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("conn write fail err=", err)
		return
	}

	//处理服务器端返回的数据

	//1.创建一个transfer 实例
	tf := utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg fail err=", err)
		return
	}
	//将mes的data反序列化成loginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		//fmt.Println("登录成功")
		//不显示自己在线
		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			if v == userId {
				continue
			}
			fmt.Println("用户id：\t", v)
			//完成客户端的onlineUsers初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println("\n\n")
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯.如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端.
		go serverProcessMes(conn)

		//1. 显示我们的登录成功的菜单[循环]..
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
func (this UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net dial err=", err)
		return
	}
	//延时关闭
	defer func() { conn.Close() }()
	//2.准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建一个LoginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	//4.将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
	}
	//5.把data赋给mes.data
	mes.Data = string(data)
	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
	}
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}
	//读取服务器发回的信息
	mes, err = tf.ReadPkg() //mes就是registerResMes
	if err != nil {
		fmt.Println("readPkg conn fail err=", err)
		return
	}
	//将mes的data部分反序列化成registerresmes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json unmarshal fail err=", err)
		return
	}
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
