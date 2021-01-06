/**
 * @Author haifengli
 * @Date 12:58 下午 2021/1/1
 * @Description
1.监听
2.等待客户端的连接
3.初始化的工作
 **/
package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

//func readPkg(conn net.Conn) (mes message.Message, err error) {
//	buf := make([]byte, 8096)
//	fmt.Println("读取客户端发送的数据-------")
//	//conn.read 在conn没有关闭的情况下，才会阻塞
//	//如果客户端关闭了conn则就不会阻塞
//	_, err = conn.Read(buf[:4]) //读取发送消息长度
//	if err != nil {
//		return
//	}
//	//根据buf[:4]转成一个uint32类型
//	var pkgLen uint32
//	pkgLen = binary.BigEndian.Uint32(buf[0:4])
//	//根据pkgLen 读取消息内容
//	n, err := conn.Read(buf[0:pkgLen])
//	if n != int(pkgLen) || err != nil {
//		return
//	}
//	//把pkgLen 反序列化成  message.Message
//	err = json.Unmarshal(
//		buf[:pkgLen],
//		&mes,
//	)
//	if err != nil {
//		fmt.Println("json.Unmarshal err=", err)
//		return
//	}
//	return
//
//}
//func writePkg(conn net.Conn, data []byte) (err error) {
//	//先发送一个长度给对方
//	var pkgLen uint32
//	pkgLen = uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
//	//发送长度
//	n, err := conn.Write(buf[0:4])
//	if n != 4 || err != nil {
//		fmt.Println("conn.write(byte) fail err=", err)
//		return
//	}
//	//发送data数据
//	n, err = conn.Write(data)
//	if n != int(pkgLen) || err != nil {
//		fmt.Println("conn.wriet(data) fail err=", err)
//		return
//	}
//	return
//}

//处理和客户端的通讯
func process(conn net.Conn) {
	//延迟关闭conn
	//defer conn.Close()
	defer func() { _ = conn.Close() }()
	//循环的读取客户端发送的信息
	//for {
	//	//v1.0
	//	//buf := make([]byte, 8096)
	//	//fmt.Println("读取客户端发送的数据-------")
	//	//n, err := conn.Read(buf[:4])
	//	//if n != 4 || err != nil {
	//	//	fmt.Println("conn.Read err=", err)
	//	//	return
	//	//}
	//	//fmt.Println("读到的buf=", buf[:4])
	//	//v1.1
	//	//读取数据包，直接封装成一个函数readPkg()，返回Message，err
	//	mes, err := readPkg(conn)
	//	if err != nil {
	//		if err == io.EOF {
	//			fmt.Println("客户端已退出--------")
	//			return
	//		} else {
	//			fmt.Println("readPkg err=", err)
	//			return
	//		}
	//	}
	//	err =serverProcessMes(conn, &mes)
	//	if err!=nil{
	//		return
	//	}
	//}
	//v1.2这里调用总控，创建一个总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.processs()
	if err != nil {
		fmt.Println("客户端和服务器通讯错误，err=", err)
		return
	}
}

//编写一个ServerProcessMes 函数
//功能： 根据客户端发送消息种类不同，决定调用哪个函数来处理
//func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
//	switch mes.Type {
//	case message.LoginMesType:
//		//处理登录
//		err = serverProcessLogin(conn, mes)
//	case message.LoginResMesType:
//	//处理用户注册
//	default:
//		fmt.Println("消息类型不存在，无法处理-------")
//
//	}
//
//	return
//}
//func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
//	//核心代码，验证用户登录
//	//先从mes中取出mes.data,并直接反序列化成LoginMes 实例
//	var loginMes message.LoginMes
//	err = json.Unmarshal([]byte(mes.Data), &loginMes)
//	if err != nil {
//		fmt.Println("json.Unmarshal fail err = ", err)
//		return
//	}
//	//连接响应数据，返回客户端信息
//	//1.先声明一个resMes
//	var resMes message.Message
//	resMes.Type = message.LoginResMesType
//	//2.声明一个LoginResMes,并完成赋值
//	var loginResMes message.LoginResMes
//	//如果用户id=100 ，密码=123 ，认为合法
//	if loginMes.UserId == 100 && loginMes.UserPwd == "123" {
//		loginResMes.Code = 200
//	} else {
//		loginResMes.Code = 500
//		loginResMes.Error = "该用户不存在，请注册在使用--------"
//	}
//	//3.将loginResMes 序列化
//	data, err := json.Marshal(loginResMes)
//	if err != nil {
//		fmt.Println("json.Marshal fail err=", err)
//		return
//	}
//	//4.将data 赋值给resMes
//	resMes.Data = string(data)
//	//5.对resMes 进行序列化，准备发送
//	data, err = json.Marshal(resMes)
//	if err != nil {
//		fmt.Println("json.Marshal fail err=", err)
//		return
//	}
//	//6.发送data，将其封装到writePkg函数
//	err = writePkg(conn, data)
//
//	return
//}
//完成对UserDao的初始化工作
func initUserDao() {
	//这里的pool本身就是一个全局的变量
	//这里需要注意一个初始化顺序问题 initPool 在initUsrDao
	model.MyUserDao = model.NewUserDao(pool)
}
func main() {
	//当服务器启动时，初始化redis连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
	//提示信息
	fmt.Println(" 服务器在8888端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	//defer listen.Close()
	defer func() { _ = listen.Close() }()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//如果监听成功，就等待客户端来连接服务器
	for {
		fmt.Println(" 等待客户端来连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		go process(conn)
	}
}
