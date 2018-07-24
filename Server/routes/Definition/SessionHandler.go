package Definition

import "net"
import (
	"../../PB"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"time"
)

const (
	newSession int32 = 1
)
const (
	clienttimeout int64 = 2000
)

func (PSSM *PacketSessionServerManager) InitSessionHandler() {
	PSSM.Handlers = make(map[int32]interface{})
	PSSM.Handlers[newSession] = PacketSessionServerManager.NewSession
}
func (PSSM *PacketSessionServerManager) NewSession(chanpacket *chan *ClientPacket, conn net.Conn) *Session {
	self := &Session{
		nil,
		0,
		nil,
		0,
		"",
		nil,
		nil,
		nil,
		nil,
		nil,
		false,
		false,
		nil,
		nil,
		0,
	}
	return self
}
func (self *Session) Run() {
	self.ReadPacketCoroutine()
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
func (self *Session) Close() {
	self.connected = false
	//log.Printf("Session %p Close!", self)
	self.Conn.Close()
}
func (self *Session) ReadPacketCoroutine() {
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
					//self.packetsGet.Put(packet)
					*self.chan_packet <- &ClientPacket{Session: self, Packet: packet}
					//self.ProcessPacket(self, packet)
				}
			}
		}
	}
}
func (self *Session) GetPacket() (*Packet, error) {
	self.Packet = NewPacket(0, 0)
	// 读header
	_, err := self.readData(self.Packet.Data, PACK_HEAD_SIZE)
	if err != nil {
		return nil, err
	}

	//log.Print("Session get header")
	// 读够了数据，解析头（同时会扩大data到必须的大小）
	self.Packet.parseHeader()
	//log.Print("Session parseHeader finish")
	if self.Packet.MsgLength > 0 {
		// 读取包体
		_, err = self.readData(self.Packet.Data[PACK_HEAD_SIZE:], self.Packet.MsgLength)
		if err != nil {
			return nil, err
		}
	}
	//log.Print("Session get body")
	if self.Packet.HasReqId {
		self.Packet.ReqId = binary.LittleEndian.Uint32(self.Packet.Data[len(self.Packet.Data)-4:])
		self.Packet.Data = self.Packet.Data[:len(self.Packet.Data)-4]
	}
	return self.Packet, nil
}
func (self *Session) readData(to []byte, need int) (int, error) {
	n, err := io.ReadAtLeast(self.Conn, to, need)
	if err != nil {
		//log.Printf("Session %p get data error:%s", self, err)
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
func (self *Session) Send(packet *Packet) (int, error) {
	//log.Printf("Session %p Sending..., conn=%p", self, self.conn)
	packet.Pack()
	self.Conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	i, err := self.Conn.Write(packet.Data)
	if err != nil {
		self.Conn.Close()
	}
	return i, err
}
func (self *Packet) parseHeader() {
	self.MessageId = binary.LittleEndian.Uint32(self.Data[OFFSET_MESSAGE_ID : OFFSET_MESSAGE_ID+4])
	self.MsgLength = int(binary.LittleEndian.Uint32(self.Data[OFFSET_MESSAGE_LEN : OFFSET_MESSAGE_LEN+4]))
	self.ErrCode = binary.LittleEndian.Uint32(self.Data[OFFSET_ERRCODE : OFFSET_ERRCODE+4])
	self.Version = binary.LittleEndian.Uint16(self.Data[OFFSET_VERSION : OFFSET_VERSION+2])
	if self.Version>>15 == 1 {
		self.HasReqId = true
		self.Version = self.Version << 1 >> 1
	}
	if self.MsgLength > 0 {
		self.CheckSize(self.MsgLength + PACK_HEAD_SIZE)
	}
}
func (self *Session) SendPbmsg(messageid PB.Message, errcode uint32, pbmsg proto.Message) (int, error) {
	pkt, err := NewPbPacket(messageid, errcode, pbmsg)
	if err != nil {
		fmt.Println("Marshal SendPbmsg err :", messageid)
		return 0, err
	}
	return self.Send(pkt)
}
func (self *Session) GetRemoteAddr() string {
	return self.Conn.RemoteAddr().String()
}
