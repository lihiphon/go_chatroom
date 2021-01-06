/**
 * @Author haifengli
 * @Date 10:57 上午 2021/1/5
 * @Description
用户登录后，可以得到当前在线用户列表
完成对onlineUsers这个map增删改查
 **/
package processs

import (
	"chatroom/client/model"
	"chatroom/common/message"
	"fmt"
)

//客户端要维护的map
var onlineUsers = make(map[int]*message.User, 10)
var CurUser model.CurUser //我们在用户登录成功后，完成对CurUser初始化
//在客户端显示当前在线的用户
func outputOnlineUser() {
	//遍历onlineUsers
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id：\t", id)
	}
}

//编写一个方法，处理返回的notifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	//适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { //原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}
