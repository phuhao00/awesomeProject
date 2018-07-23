package Network

import (
	"../PB"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	//"time"
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	SessionState_Logout = 1
	SessionState_Login  = 1 << 1
	SessionState_InGame = 1 << 2
)

type Session struct {
	Conn    net.Conn
	Userid  int64
	User    interface{}
	State   int32
	Account string

	packet          *Packet
	HandleConnected func()
	ProcessPacket   func(session *Session, packet *Packet)
	HandleClose     func(session *Session)
	chan_packet     *chan *ClientPacket
	IsClient        bool
	connected       bool
	syncRequests    sync.Map

	cipher    *rc4.Cipher
	reqIdNext uint32
}

func NewSession(chanpacket *chan *ClientPacket, conn net.Conn) *Session {
	self := &Session{
		Userid:      0,
		Conn:        conn,
		chan_packet: chanpacket,
		packet:      NewPacket(0, 0),
		State:       SessionState_Logout,
	}
	return self
}

func (self *Session) Send(packet *Packet) (int, error) {
	//log.Printf("Session %p Sending..., conn=%p", self, self.conn)
	packet.Pack()
	self.Conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	i, err := self.Conn.Write(packet.data)
	if err != nil {
		self.Conn.Close()
	}
	return i, err
}

func (self *Session) InitRC4(key string) {
	self.cipher, _ = rc4.NewCipher([]byte(key))
}

func (self *Session) RC4EnCode(data []byte) []byte {
	m_data := []byte{}
	m_data = data[0:]
	if self.cipher != nil {
		m_data = make([]byte, len(data))
		self.cipher.XORKeyStream(m_data, data)
	}
	return m_data
}

func (self *Session) SendPbmsg(messageid pbID.Message, errcode uint32, pbmsg proto.Message) (int, error) {
	pkt, err := NewPbPacket(messageid, errcode, pbmsg)
	if err != nil {
		fmt.Println("Marshal SendPbmsg err :", messageid)
		return 0, err
	}
	return self.Send(pkt)
}

func (self *Session) SendPbmsgWithRequestId(messageid pbID.Message, errcode uint32, pbmsg proto.Message, reqId uint32) (int, error) {
	pkt, err := NewPbPacket(messageid, errcode, pbmsg)
	if err != nil {
		fmt.Println("Marshal SendPbmsg err :", messageid)
		return 0, err
	}
	if reqId > 0 {
		pkt.hasReqId = true
		pkt.ReqId = reqId
	}
	return self.Send(pkt)
}

func (self *Session) SendAsyncRequest(request *AsyncPBRequest) (int, error) {
	return self.SendPbmsgWithRequestId(request.MsgId, request.Errcode, request.Request, request.Id)
}

//func (self *Session) Cor_SyncSendPbmsg(messageid pbID.Message, errcode uint32, pbmsg proto.Message, respMsg proto.Message) *SyncRequest {
//	req := &SyncRequest{
//		0,
//		messageid,
//		errcode,
//		pbmsg,
//		respMsg,
//		0,
//		make(chan *Packet),
//		time.Now().Unix(),
//	}
//	reqid := atomic.AddUint32(&self.reqIdNext, 1)
//	self.SendPbmsgWithRequestId(messageid,errcode,pbmsg, reqid)
//	packet := <- req.ch_response
//	if packet == nil {
//		req.responseErrCode = -1
//	} else {
//		req.responseErrCode = int32(packet.ErrCode)
//		err := proto.Unmarshal(packet.MessageData(), req.response)
//		if err != nil {
//			req.responseErrCode = -2
//		}
//	}
//	return req
//}

func (self *Session) Cor_SyncSendPbmsg(request *SyncPBRequest) {
	reqid := atomic.AddUint32(&self.reqIdNext, 1)
	request.Id = reqid
	if request.Ch_response == nil {
		request.Ch_response = make(chan *Packet)
	}
	self.syncRequests.Store(request.Id, request)
	self.SendPbmsgWithRequestId(request.MsgId, request.Errcode, request.Request, request.Id)
	packet := <-request.Ch_response
	if packet == nil {
		request.ResponseErrCode = -1
	} else {
		request.ResponseErrCode = int32(packet.ErrCode)
		if request.Response != nil {
			err := proto.Unmarshal(packet.MessageData(), request.Response)
			if err != nil {
				request.ResponseErrCode = -2
			}
		}
	}
	return
}

func (self *Session) GetUserID() int64 {
	return self.Userid
}

func (self *Session) GetPacket() (*Packet, error) {
	self.packet = NewPacket(0, 0)
	// 读header
	_, err := self.readData(self.packet.data, PACK_HEAD_SIZE)
	if err != nil {
		return nil, err
	}

	//log.Print("Session get header")
	// 读够了数据，解析头（同时会扩大data到必须的大小）
	self.packet.parseHeader()
	//log.Print("Session parseHeader finish")
	if self.packet.msgLength > 0 {
		// 读取包体
		_, err = self.readData(self.packet.data[PACK_HEAD_SIZE:], self.packet.msgLength)
		if err != nil {
			return nil, err
		}
	}
	//log.Print("Session get body")
	if self.packet.hasReqId {
		self.packet.ReqId = binary.LittleEndian.Uint32(self.packet.data[len(self.packet.data)-4:])
		self.packet.data = self.packet.data[:len(self.packet.data)-4]
	}
	return self.packet, nil
}

var clienttimeout int64 = 2000

//go
func (self *Session) Run() {
	for {
		packet, _ := self.GetPacket()
		if packet == nil {
			break
		}
		self.connected = true
		if !self.IsClient {
			self.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(clienttimeout)))
		}
		if packet.MessageId == 0 {
			if !self.IsClient {
				self.SendPbmsg(0, 0, nil)
			}
		} else {
			if packet.ReqId > 0 {
				v, _ := self.syncRequests.Load(packet.ReqId)
				if v == nil {
					if self.chan_packet != nil {
						*self.chan_packet <- &ClientPacket{Session: self, Packet: packet}
						//self.ProcessPacket(self, packet)
					}
				} else {
					v.(*SyncPBRequest).Ch_response <- packet
					self.syncRequests.Delete(packet.ReqId)
				}
			} else {
				if self.chan_packet != nil {
					*self.chan_packet <- &ClientPacket{Session: self, Packet: packet}
					//self.ProcessPacket(self, packet)
				}
			}
		}
	}
	self.Close()
	if self.HandleClose != nil {
		self.HandleClose(self)
	}
	self.syncRequests.Range(func(k, v interface{}) bool {
		v.(*SyncPBRequest).Ch_response <- nil
		self.syncRequests.Delete(k)
		return true
	})
}

func (self *Session) readData(to []byte, need int) (int, error) {
	n, err := io.ReadAtLeast(self.Conn, to, need)
	if err != nil {
		log.Printf("Session %p get data error:%s", self, err)
	}
	return n, err

	//geted, n := 0, 0
	//var err error
	//for {
	//	n, err = self.conn.Read(self.packet.data)
	//
	//	log.Printf("Session %p get data,len=%d,err=%s", self, len(self.packet.data), err)
	//	if err != nil {
	//		return n, err
	//	}
	//	geted += n
	//	if(geted >= need) {
	//		break
	//	}
	//}
	//return geted,err
}

func (self *Session) Close() {
	self.connected = false
	log.Printf("Session %p Close!", self)
	self.Conn.Close()
}

func (self *Session) GetRemoteAddr() string {
	return self.Conn.RemoteAddr().String()
}

type ClientPacket struct {
	Session *Session
	Packet  *Packet
}

type SyncPBRequest struct {
	Id              uint32
	MsgId           pbID.Message
	Errcode         uint32
	Request         proto.Message
	Response        proto.Message
	ResponseErrCode int32
	Ch_response     chan *Packet
	Time            int64
}

type AsyncPBRequest struct {
	Id      uint32
	MsgId   pbID.Message
	Errcode uint32
	Request proto.Message
}

type ClientPacketHandler func(clientpacket *ClientPacket)
