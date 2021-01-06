/**
 * @Author haifengli
 * @Date 9:01 下午 2021/1/1
 * @Description
1.显示登录成功界面
2.保持和服务器通讯（启动协程）
3.当读取服务器发送的消息后，就会显示在界面
 **/
package processs

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("-----恭喜xxx登录成功-----")
	fmt.Println("---1.显示在线用户列表---")
	fmt.Println("---2.发送消息---")
	fmt.Println("---3.信息列表---")
	fmt.Println("---4.退出系统---")
	fmt.Println("请选择（1-4）：")
	var key int
	var content string

	//因为，我们总会使用到SmsProcess实例，因此我们将其定义在swtich外部
	smsProcess := &SmsProcess{}
	_, _ = fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表---")
		outputOnlineUser()
	case 2:
		//fmt.Println("发送消息---")
		fmt.Println("你要说什么：")
		_, _ = fmt.Scanf("%s\n", &content)
		_ = smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表---")
	case 4:
		fmt.Println("退出系统---")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确---")
	}
}

//和服务器保持通讯
func serverProcessMes(conn net.Conn) {
	//创建一个transfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg fail err=", err)
			return
		}
		//如果读取到消息，将下一步处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线了
			//1.取出notifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			_ = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2.把这个用户的信息，状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}

		//fmt.Printf("mes=%v\n",mes)
	}

}