/**
 * @Author haifengli
 * @Date 10:18 上午 2021/1/4
 * @Description
根据业务逻辑需要，自定义一些错误
 **/
package model

import "errors"

var (
	ERROR_USER_NOTEXISTS = errors.New("用户不存在")
	ERROR_USER_EXISTS    = errors.New("用户已存在")
	ERROR_USER_PWD       = errors.New("密码不正确")
)
