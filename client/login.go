package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {

	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据~~~...")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭了 conn 则，就不会阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4] 转成一个 uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	//根据 pkgLen 读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes！！
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}
func login(userId int, userPwd string) (err error) {
	//fmt.Printf("userId=%d userPwd=%s\n", userId, userPwd)

	//return nil
	//	1、客户端连接到服务器
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("net Dial err=", err)
		return
	}
	//	延迟关闭连接
	// defer conn.Close()
	defer func() { _ = conn.Close() }()
	//2、准备通过conn发送信息到服务端
	var mes message.Message
	mes.Type = message.LoginMesType
	//3、loginMes 创建实体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//4.将loginMes实例 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//5.把data赋给 mes.data 字段  把loginMes 再次包装
	mes.Data = string(data)
	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//7.此时data 就是要发到服务端的消息
	//首先先把data的长度发给服务端
	//先获取data的长度-->转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//8.发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn,write bytes fail err=", err)
	}
	fmt.Printf("客户端发送的消息长度=%d,内容=%s", len(data), string(data))
	//9.发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.write(data) fail err=", err)
		return
	}
	//time.Sleep(time.Second * 20) //休眠20秒
	//处理服务器返回来的消息
	mes, err = readPkg(conn)
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	//将mes的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")

	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
