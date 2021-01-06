/**
 * @Author haifengli
 * @Date 10:22 上午 2021/1/4
 * @Description
对user对象操作的各种方法。增删改查
 **/
package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//在服务器启动后，就初始化一个userDao实例
//把它做成全局的变量，在需要和redis操作时，就直接可以使用
var MyUserDao *UserDao

//定义一个UserDao结构体
//完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//userDao 可以提供的方法
//1.根据用户id 返回一个User实例和err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//通过给定的id去redis中查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在users哈希中，没有对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	//这里我们需要把res反序列化成user实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

//完成登录的校验
//1.Login 完成对应的验证
//2.如果用户的id和pwd都正确，则返回一个user实例
//3.如果用的id和pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先充UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer func() { conn.Close() }()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}
	//这时证明这个用户是获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

//用户注册
func (this *UserDao) Register(user *message.User) (err error) {
	//先从userdao的连接池中取出一个连接
	conn := this.pool.Get()
	defer func() { conn.Close() }()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	//这时，说明id在redis还没有，可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json marshal fail err=", err)
		return
	}
	//入库
	_, err = conn.Do("Hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
