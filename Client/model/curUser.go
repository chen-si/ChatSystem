package model

import (
	"ChatSystem/common/message"
	"net"
)

//因为在客户端很多地方都会用到CurUser 所以全局\
type CurUser struct {
	Conn net.Conn
	message.User
}
