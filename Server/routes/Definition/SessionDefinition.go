package Definition

import (
	"crypto/rc4"
	"net"
	"sync"
)

type Session struct {
	Conn            net.Conn //连接
	Userid          int64    //用户Id
	User            interface{}
	State           int32
	Account         string
	Packet          *Packet //包
	HandleConnected func()
	ProcessPacket   func(session *Session, packet *Packet)
	HandleClose     func(session *Session)
	chan_packet     *chan *ClientPacket
	IsClient        bool
	Connected       bool
	syncRequests    sync.Map
	cipher          *rc4.Cipher
	reqIdNext       uint32
}
