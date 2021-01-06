/**
 * @Author haifengli
 * @Date 8:59 下午 2021/1/4
 * @Description
创建一个onlineUsers map[int] *UserProcss维护在线用户
 **/
package processs

import "fmt"

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//userMgr实例在服务器端有且一个，因此将其定义为全局变量
var userMgr *UserMgr

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//删除
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

//返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据id返回对应的值
func (this *UserMgr) GetOnlineUserById(userid int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userid] //从map中取出一个值，带检测方式
	if !ok {
		err = fmt.Errorf("用户%d 不存在", userid)
		return
	}
	return

}
