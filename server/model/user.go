/**
 * @Author haifengli
 * @Date 10:14 上午 2021/1/4
 * @Description
定义一个User结构体
 **/
package model

type User struct {
	//确定字段信息
	//为了序列化和反序列化成功，必须保证用户信息的json字符的key和结构体字段对应的tag名字一致
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
