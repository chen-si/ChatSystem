package utils

import (
	"ChatSystem/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Transfer struct {
	//分析他应该有那些字段
	Conn net.Conn
	Buf  [8096]byte //传输时使用的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//defer conn.Close()
	//buf := make([]byte,8096)
	fmt.Println("读取客户端发送的消息....")
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("conn.Read(buf[:4]) error")
		return
	}
	//将读取到的buf[:4]转换成uint32类型
	var PkgLen uint32
	PkgLen = binary.BigEndian.Uint32(this.Buf[:4])

	//根据PkgLen读取数据和话剧
	n, err := this.Conn.Read(this.Buf[:PkgLen])
	if err != nil || n != int(PkgLen) {
		err = errors.New("conn.Read(buf[:PkgLen]) error")
		return
	}

	//将读取到的数据反序列化
	err = json.Unmarshal(this.Buf[:PkgLen], &mes)
	if err != nil {
		err = errors.New("json.Unmarshal(buf[:PkgLen],&mes) error")
		return
	}

	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkglen uint32
	pkglen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkglen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) error:", err)
		return
	}
	//发送消息本身
	//发送消息本身
	n, err = this.Conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) error:", err)
		return
	}
	return
}
