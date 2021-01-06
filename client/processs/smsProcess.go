/**
 * @Author haifengli
 * @Date 9:01 下午 2021/1/1
 * @Description
1.处理和信息相关逻辑
2.私聊
3.群发
 **/
package processs

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊的消息
func (sp *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType
	//2.创建一个smsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content               //内容.
	smsMes.UserId = CurUser.UserId         //
	smsMes.UserStatus = CurUser.UserStatus //
	//3.序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("sendGroupMes json.Marshal fail err=", err)
		return
	}
	mes.Data = string(data)
	//4.对mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("sendGroupMes json.Marshal fail err=", err)
		return
	}
	//5.将mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("sendGroupMes err=", err)
		return
	}
	return
}
