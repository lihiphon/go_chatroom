/**
 * @Author haifengli
 * @Date 12:58 下午 2021/1/1
 * @Description
总的处理器
根据客户端的请求，调用对应的处理器，完成相应的任务
 **/
package main

import (
	"chatroom/common/message"
	"chatroom/server/processs"
	"chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

//先创建一个Processor 的结构体
type Processor struct {
	Conn net.Conn
}

//编写一个serverProcessMes函数
//功能：根据客户端发送消息种类不同，决定调用哪个函数处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//创建一个UserProcess实例
		up := processs.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &processs.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//创建一个smsprocss实例完成转发群聊消息
		smsProcss := &processs.SmsProcess{}
		smsProcss.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理----------")

	}
	return
}

//处理客户端发送的消息
func (this *Processor) processs() (err error) {
	//循环的读取客户端发送消息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg（），返回Message，err
		//创建一个Transfer 实现完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端退出------")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
