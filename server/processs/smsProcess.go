/**
 * @Author haifengli
 * @Date 12:57 下午 2021/1/1
 * @Description
1.处理和信息相关的请求
2.群聊
3.点对点聊天
 **/
package processs

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

func (sp *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的onlineUsers map[int]*UserProcess
	//将消息转发取出
	//取出mes的内容SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json unmarshal fail err=", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		sp.SendMesToEachOnlineUser(data, up.Conn)
	}

}
func (sp *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败err=", err)
	}
}
