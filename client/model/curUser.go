/**
 * @Author haifengli
 * @Date 9:41 下午 2021/1/5
 * @Description

 **/
package model

import (
	"chatroom/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
