package Definition

import "net"

type ProtocolType string

const (
	SessionState_Logout = 1
	SessionState_Login  = 1 << 1
	SessionState_InGame = 1 << 2
)

type ServerOptions struct {
	WriteBufferLen int
	ReadBufferLen  int
	StreamMode     bool
	Timeout        int
}

type ttt struct {
	gg string
}

var ServerOptionsDefault ServerOptions = ServerOptions{
	WriteBufferLen: 8 * 1024,
	ReadBufferLen:  2 * 1024,
	StreamMode:     true,
	Timeout:        10,
}

type Server struct {
	Id              string
	Addr            string
	Listener        net.Listener
	Protocol        ProtocolType
	Options         ServerOptions
	Chan_Connection chan *Session
	Chan_Packet     chan *ClientPacket
	Chan_Close      chan *Session
}
