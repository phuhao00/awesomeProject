package Network

import (
	"../kcp"
	"fmt"
	"net"
)

type ProtocolType string

const (
	ProtocolType_TCP ProtocolType = "tcp"
	ProtocolType_KCP ProtocolType = "kcp"
	ProtocolType_UDP ProtocolType = "udp"
)

type ServerOptions struct {
	WriteBufferLen int
	ReadBufferLen  int
	StreamMode     bool
	Timeout        int
}

var ServerOptionsDefault ServerOptions = ServerOptions{
	WriteBufferLen: 8 * 1024,
	ReadBufferLen:  2 * 1024,
	StreamMode:     true,
	Timeout:        10,
}

type Server struct {
	id   string
	addr string
	//listener *StoppableListener
	listener        net.Listener
	protocol        ProtocolType
	options         ServerOptions
	Chan_Connection chan *Session
	Chan_Packet     chan *ClientPacket
	Chan_Close      chan *Session
}

func NewServer(id string, addr string, protocol ProtocolType, options *ServerOptions) *Server {
	if protocol == ProtocolType_UDP {
		panic("暂不支持ProtocolType_UDP!")
	}
	self := &Server{
		id:       id,
		addr:     addr,
		protocol: protocol,
	}
	if options == nil {
		self.options = ServerOptionsDefault
	} else {
		self.options = *options
	}
	return self
}

func (self *Server) Run() error {
	var err error = nil
	fmt.Println(self.protocol, "listen at :", self.addr)
	//

	switch self.protocol {
	case ProtocolType_TCP:
		self.listener, err = net.Listen("tcp", self.addr)
	case ProtocolType_KCP:
		self.listener, err = kcp.ListenWithOptions(self.addr, nil, 10, 3)
	default:
		panic("不支持的协议类型")
	}

	if err != nil {
		fmt.Println("listen error:", err)
		return err
	}

	switch self.protocol {
	case ProtocolType_KCP:
		kcplistener := self.listener.(*kcp.Listener)
		kcplistener.SetReadBuffer(2 * 1024)
		kcplistener.SetWriteBuffer(8 * 1024)
		kcplistener.SetDSCP(46)
	}

	for {
		conn, err := self.listener.Accept()

		if err != nil {
			netErr, ok := err.(net.Error)

			//超时，无视，继续接收下一个
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}

			//其他错误。停止
			fmt.Println("accept error:", err)
			return err
		}

		session := NewSession(&self.Chan_Packet, conn)
		session.HandleClose = self.OnSessionClose
		switch self.protocol {
		case ProtocolType_KCP:
			KCPConn := session.Conn.(*kcp.UDPSession)
			KCPConn.SetReadBuffer(self.options.ReadBufferLen)
			KCPConn.SetWriteBuffer(self.options.WriteBufferLen)
			KCPConn.SetWindowSize(4096, 4096)
			KCPConn.SetNoDelay(1, 10, 2, 1)
			KCPConn.SetDSCP(46)
			KCPConn.SetMtu(1400)
		}

		session.Send(NewPacket(0, 0))

		fmt.Println(self.protocol, "new connected", session.GetRemoteAddr())
		if self.Chan_Connection != nil {
			self.Chan_Connection <- session
		}
		go session.Run()
	}
	return nil
}

func (self *Server) Close() {
	self.listener.Close()
}

func (self *Server) OnSessionClose(session *Session) {
	self.Chan_Close <- session
}
