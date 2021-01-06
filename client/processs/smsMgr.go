/**
 * @Author haifengli
 * @Date 3:48 下午 2021/1/6
 * @Description

 **/
package processs

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) { //这个地方mes一定smsmes
	//显示即可
	//1.反序列化mes.data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json unmarshal err=", err)
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说：\t%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()

}
