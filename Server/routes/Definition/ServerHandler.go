package Definition

import (
	"../../kcp"
	"fmt"
	"net"
)

const (
	newServer int32 = 1
)
const (
	ProtocolType_TCP ProtocolType = "tcp"
	ProtocolType_KCP ProtocolType = "kcp"
	ProtocolType_UDP ProtocolType = "udp"
)

func (PSSM *PacketSessionServerManager) InitServerHandler() {
}
func (PSSM *PacketSessionServerManager) NewServer(id string, addr string, protocol ProtocolType, options *ServerOptions) *Server {
	if protocol == ProtocolType_UDP {
		panic("暂不支持ProtocolType_UDP!")
	}
	newServer := &Server{
		id,
		addr,
		nil,
		protocol,
		ServerOptionsDefault,
		nil,
		nil,
		nil,
	}
	if options != nil {
		newServer.Options = *options
	}
	return newServer
}
func (self *Server) Run() error {
	var err error = nil
	fmt.Println(self.Protocol, "listen at :", self.Addr)
	switch self.Protocol {
	case ProtocolType_TCP:
		self.Listener, err = net.Listen("tcp", self.Addr)
	case ProtocolType_KCP:
		self.Listener, err = kcp.ListenWithOptions(self.Addr, nil, 10, 3)
	default:
		panic("不支持的协议类型")
	}
	if err != nil {
		fmt.Println("listen error:", err)
		return err
	}

	switch self.Protocol {
	case ProtocolType_KCP:
		kcplistener := self.Listener.(*kcp.Listener)
		kcplistener.SetReadBuffer(2 * 1024)
		kcplistener.SetWriteBuffer(8 * 1024)
		kcplistener.SetDSCP(46)
	}

	for {
		conn, err := self.Listener.Accept()

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

		session := PSSM.NewSession(&self.Chan_Packet, conn)
		session.HandleClose = self.OnSessionClose
		switch self.Protocol {
		case ProtocolType_KCP:
			KCPConn := session.Conn.(*kcp.UDPSession)
			KCPConn.SetReadBuffer(self.Options.ReadBufferLen)
			KCPConn.SetWriteBuffer(self.Options.WriteBufferLen)
			KCPConn.SetWindowSize(4096, 4096)
			KCPConn.SetNoDelay(1, 10, 2, 1)
			KCPConn.SetDSCP(46)
			KCPConn.SetMtu(1400)
		}

		session.Send(PSSM.NewPacket(0, 0))

		fmt.Println(self.Protocol, "new connected", session.GetRemoteAddr())
		if self.Chan_Connection != nil {
			self.Chan_Connection <- session
		}
		go session.Run()
	}
	return nil
}
func (self *Server) Close() {
	self.Listener.Close()
}
func (self *Server) OnSessionClose(session *Session) {
	self.Chan_Close <- session
}
