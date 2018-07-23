package Network

import (
	"../PB"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
)

//14字节包头：4长度 2版本 4消息id 4errcode
const (
	PACK_HEAD_SIZE     = 14
	PACKVERSION        = 1
	OFFSET_MESSAGE_LEN = 0
	OFFSET_VERSION     = 4
	OFFSET_MESSAGE_ID  = 6
	OFFSET_ERRCODE     = 10
	MULTI_PACKET_ID    = 1
)

type Packet struct {
	msgLength int
	version   uint16 //暂时尚未使用。首位征用，用于表示是否有ReqId
	MessageId uint32
	ErrCode   uint32
	data      []byte
	hasReqId  bool
	ReqId     uint32 //在包尾
}

//var PacketPool = sync.Pool{
//	New: func() interface{} {
//		pkt := new(Packet)
//		pkt.data = make([]byte, PACK_HEAD_SIZE, 128)
//		pkt.data[OFFSET_VERSION] = PACKVERSION
//		pkt.data[OFFSET_ERRCODE] = 0
//		return pkt
//	},
//}

func NewPacket(messageId uint32, errCode uint32) *Packet {
	//pkt := PacketPool.Get().(*Packet)
	pkt := &Packet{
		data: make([]byte, PACK_HEAD_SIZE, 128),
	}

	pkt._init(messageId, errCode)
	return pkt
}

func NewMultiPacket() *Packet {
	return NewPacket(MULTI_PACKET_ID, 0)
}

func (self *Packet) _init(messageId uint32, errCode uint32) *Packet {
	self.MessageId = messageId
	self.ErrCode = errCode
	self.reset()
	return self
}

func (self *Packet) Append(data []byte) *Packet {
	self.checkSize(len(self.data) + len(data))
	copy(self.data[len(self.data):], data)
	return self
}

func (self *Packet) AppendPacket(packet *Packet) *Packet {
	packet.Pack()
	return self.Write(packet.data, len(self.data))
}

func (self *Packet) AppendUint32(v uint32) *Packet {
	self.checkSize(len(self.data) + 4)
	binary.LittleEndian.PutUint32(self.data[len(self.data)-4:len(self.data)], v)
	return self
}

func (self *Packet) AppendUint16(v uint16) *Packet {
	self.checkSize(len(self.data) + 2)
	binary.LittleEndian.PutUint16(self.data[len(self.data)-2:len(self.data)], v)
	return self
}

func (self *Packet) Pack() *Packet {
	binary.LittleEndian.PutUint32(self.data[OFFSET_MESSAGE_ID:OFFSET_MESSAGE_ID+4], self.MessageId)
	binary.LittleEndian.PutUint32(self.data[OFFSET_ERRCODE:OFFSET_ERRCODE+4], self.ErrCode)
	if self.ReqId > 0 {
		ver := PACKVERSION
		ver |= (1 << 15)
		binary.LittleEndian.PutUint16(self.data[OFFSET_VERSION:OFFSET_VERSION+2], uint16(ver))
		self.AppendUint32(self.ReqId)
	} else {
		binary.LittleEndian.PutUint16(self.data[OFFSET_VERSION:OFFSET_VERSION+2], PACKVERSION)
	}
	self.msgLength = len(self.data) - PACK_HEAD_SIZE
	binary.LittleEndian.PutUint32(self.data[OFFSET_MESSAGE_LEN:OFFSET_MESSAGE_LEN+4], uint32(self.msgLength))
	return self
}

func (self *Packet) reset() *Packet {
	self.data = self.data[:PACK_HEAD_SIZE]
	self.hasReqId = false
	return self
}

func (self *Packet) checkSize(len int) {
	if len > cap(self.data) {
		newBuf := make([]byte, len+1024) //因为要丢到pool中，所以这里扩展为比实际需要的大一些（1k），避免频繁扩展
		copy(newBuf, self.data)
		self.data = newBuf
	}
	self.data = self.data[:len] //裁切为确切的大小
}

func (self *Packet) parseHeader() {
	self.MessageId = binary.LittleEndian.Uint32(self.data[OFFSET_MESSAGE_ID : OFFSET_MESSAGE_ID+4])
	self.msgLength = int(binary.LittleEndian.Uint32(self.data[OFFSET_MESSAGE_LEN : OFFSET_MESSAGE_LEN+4]))
	self.ErrCode = binary.LittleEndian.Uint32(self.data[OFFSET_ERRCODE : OFFSET_ERRCODE+4])
	self.version = binary.LittleEndian.Uint16(self.data[OFFSET_VERSION : OFFSET_VERSION+2])
	if self.version>>15 == 1 {
		self.hasReqId = true
		self.version = self.version << 1 >> 1
	}
	if self.msgLength > 0 {
		self.checkSize(self.msgLength + PACK_HEAD_SIZE)
	}
}

func (self *Packet) MessageData() []byte {
	return self.data[PACK_HEAD_SIZE:]
}

/*
func (self *Packet) RemoveData(start int, leng int){
	self.data = append(self.data[:start], self.data[start + leng:]...)
	self.Pack()
}*/

func NewPbPacket(messageid pbID.Message, errcode uint32, pbmsg proto.Message) (*Packet, error) {
	var mData []byte
	var err error
	if pbmsg != nil {
		mData, err = proto.Marshal(pbmsg)
		if err != nil {
			fmt.Println("Marshal SendPbmsg err :", messageid)
			return nil, err
		}
	}
	pkt := NewPacket(uint32(messageid), errcode).Append(mData)
	return pkt, nil
}
