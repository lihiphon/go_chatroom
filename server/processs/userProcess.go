/**
 * @Author haifengli
 * @Date 12:51 下午 2021/1/1
 * @Description
1.处理和用户相关的请求
2.登录
3.注册
4.注销
5.用户列表管理
 **/
package processs

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段，表示该conn是哪个用户的
	UserId int
}

//编写一个函数serverProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码。。
	//1.先从mes中取出mes.data,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//2.声明一个resMes
	var resMes message.Message

	resMes.Type = message.LoginResMesType
	//3.再声明一个LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	//到redis数据库中取完成验证
	//使用model.myuserDao 到redis中验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		//这里登录成功，把改用户放到userMgr中
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其他的在线用户，我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用的id放入到loginResMes.Usersid
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}
	//如果用户id= 100， 密码=123456, 认为合法，否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	//合法
	//	loginResMes.Code = 200
	//
	//} else {
	//	//不合法
	//	loginResMes.Code = 500 // 500 状态码，表示该用户不存在
	//	loginResMes.Error = "该用户不存在, 请注册再使用..."
	//}
	//4.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//5.将data赋值给resMes
	resMes.Data = string(data)
	//6.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//7.发送data，将其封装到writePkg函数。因为使用分层模式，我们先创建一个transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//处理用户注册
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1.先从mes中取出mes.data，并直接反序列化成registermes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json umarshal fail err= ", err)
		return
	}
	//2.声明一个resmes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	//3.使用modle.myuserdao到redis中取验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误。。。。。"
		}
	} else {
		registerResMes.Code = 200
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//4.将data赋值给resmes
	resMes.Data = string(data)
	//5.对resmes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//6.发送data 将其封装到writePkg函数中
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}
func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	//将notifyUeserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给mes.data
	mes.Data = string(data)
	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//发送，创建Transfer实例
	var tf = &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("notifyMeOnline writerPkg fail err=", err)
		return
	}
}

//通知其他的在线用户，我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userid int) {
	//遍历onlineUers,然后一个一个的发送notifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤到自己
		if id == userid {
			continue
		}
		//开始通知其他在线用户
		up.NotifyMeOnline(userid)
	}

}
