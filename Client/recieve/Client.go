package recieve

import (
	"../../Server/kcp"
	DF "../../Server/routes/Definition"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	id                 string
	server             string
	Session            *DF.Session
	protocol           DF.ProtocolType
	AutoReconnect      bool
	Chan_Packet        chan *DF.ClientPacket
	Chan_Close         chan *DF.Session
	Chan_Connection    chan *DF.Session
	Connected          bool
	chan_nextconnect   <-chan time.Time
	chan_nextheartbeat <-chan time.Time
}

func NewClient(id string, server string, protocol DF.ProtocolType) *Client {
	if protocol == DF.ProtocolType_UDP {
		panic("暂不支持ProtocolType_UDP")
	}
	self := &Client{
		id:            id,
		server:        server,
		AutoReconnect: true,
		protocol:      protocol,
	}
	return self
}

func (self *Client) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", self.server)
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
		return err
	}

	self.Connected = false
	self.Session = nil

	switch self.protocol {
	case DF.ProtocolType_TCP:
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Printf("Fatal error: %s", err.Error())
			if self.AutoReconnect {
				self.chan_nextconnect = time.After(time.Second * 5)
			}
			return err
		}
		self.Session = DF.PSSM.NewSession(&self.Chan_Packet, conn)
	case DF.ProtocolType_KCP:
		conn, err := kcp.DialWithOptions(self.server, nil, 10, 3)
		if err != nil {
			log.Printf("Fatal error: %s", err.Error())
			return err
		}
		self.Session = DF.PSSM.NewSession(&self.Chan_Packet, conn)
	}
	self.Session.IsClient = true
	self.Session.HandleClose = self.OnSessionClose

	switch self.protocol {
	case DF.ProtocolType_KCP:
		KCPConn := self.Session.Conn.(*kcp.UDPSession)
		KCPConn.SetReadBuffer(2 * 1024)
		KCPConn.SetWriteBuffer(8 * 1024)
		KCPConn.SetWindowSize(4096, 4096)
		KCPConn.SetNoDelay(1, 10, 2, 1)
		KCPConn.SetDSCP(46)
		KCPConn.SetMtu(1400)
	case DF.ProtocolType_TCP:
	}
	go self.Session.Run()
	self.Session.Send(DF.PSSM.NewPacket(0, 0))
	self.chan_nextheartbeat = time.After(time.Second * 10)
	return nil
}

func (self *Client) Send(packet *DF.Packet) (int, error) {
	//TODO:存储一个队列用于重连后重发
	//log.Printf("client=%p,id=%d,session=%p,pkt=%p", &self, self.id, self.session, packet)
	if self.Session == nil {
		return 0, errors.New(fmt.Sprintf("client %p 正在链接中，不能发包", self))
	}
	return self.Session.Send(packet)
}

func (self *Client) _ProcessPacket(packet *DF.Packet) error {
	//注意这个packet没有拷贝，是同步调用的，函数返回后下一个包会重用
	//log.Printf("%p,id=%s,Get packet:%d",self,self.Id,packet.MessageId)
	return nil
}

func (self *Client) Close() {
	if self.Session != nil {
		self.Session.Close()
	}
}

func (self *Client) OnSessionClose(session *DF.Session) {
	if self.Session != session {
		return
	}
	if self.Chan_Close != nil {
		self.Chan_Close <- session
	}
	if self.AutoReconnect {
		self.chan_nextconnect = time.After(time.Second * 5)
	}
}

func (self *Client) Run() {
	self.Connect()
	for {
		if !self.Connected && self.Session != nil && self.Session.Connected {
			self.Connected = true
			if self.Chan_Connection != nil {
				self.Chan_Connection <- self.Session
			}
		}
		select {
		case <-self.chan_nextconnect:
			self.chan_nextconnect = nil
			self.Connect()
		case <-self.chan_nextheartbeat:
			if self.Session != nil {
				self.Session.Send(DF.PSSM.NewPacket(0, 0))
			}
			self.chan_nextheartbeat = time.After(time.Second * 10)
		default:

		}
		time.Sleep(time.Millisecond * 200)
	}
}
