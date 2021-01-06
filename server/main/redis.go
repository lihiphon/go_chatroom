/**
 * @Author haifengli
 * @Date 11:00 上午 2021/1/4
 * @Description
redis 连接池
 **/
package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空闲链接数
		MaxActive:   maxActive,   //表示和数据库的最大链接数，0表示	没有限制
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) { //初始化连接代码，链接按个地址的redis
			return redis.Dial("tcp", address)
		},
	}
}
