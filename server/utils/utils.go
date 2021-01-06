/**
 * @Author haifengli
 * @Date 12:58 下午 2021/1/1
 * @Description
1.常用的工具的函数，结构体
2.提供常用的方法和函数
 **/
package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //传输时，使用的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据-------")
	//conn.read 在conn没有关闭的情况下，才会阻塞
	//如果客户端关闭了conn则就不会阻塞
	_, err = this.Conn.Read(buf[:4]) //读取发送消息长度
	if err != nil {
		return
	}
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	//根据pkgLen 读取消息内容
	n, err := this.Conn.Read(buf[0:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	//把pkgLen 反序列化成  message.Message
	err = json.Unmarshal(
		buf[:pkgLen],
		&mes,
	)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return

}
func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长处给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen)
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail err=", err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.write(bytes) fail err=", err)
		return
	}
	return
}
