package Definition

import (
	"../../PB"
	DF "./"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (PPSM *PacketSessionServerManager) initPacketHandler() {
	PPSM.Handlers = make(map[int32]interface{})
	PPSM.Handlers[NewPacket_] = PacketSessionServerManager.NewPacket
	PPSM.Handlers[NewPbPacket_] = PacketSessionServerManager.NewPbPacket
}
func (self *Packet) Init(messageId uint32, errCode uint32) {
	self.MessageId = messageId
	self.ErrCode = errCode
	Reset(self)
}
func Reset(self *Packet) {
	self.Data = self.Data[:PACK_HEAD_SIZE]
	self.HasReqId = false
}
func (self *Packet) Pack() *Packet {
	binary.LittleEndian.PutUint32(self.Data[OFFSET_MESSAGE_ID:OFFSET_MESSAGE_ID+4], self.MessageId)
	binary.LittleEndian.PutUint32(self.Data[OFFSET_ERRCODE:OFFSET_ERRCODE+4], self.ErrCode)
	if self.ReqId > 0 {
		ver := PACKVERSION
		ver |= (1 << 15)
		binary.LittleEndian.PutUint16(self.Data[OFFSET_VERSION:OFFSET_VERSION+2], uint16(ver))
		self.AppendUint32(self.ReqId)
	} else {
		binary.LittleEndian.PutUint16(self.Data[OFFSET_VERSION:OFFSET_VERSION+2], PACKVERSION)
	}
	self.MsgLength = len(self.Data) - PACK_HEAD_SIZE
	binary.LittleEndian.PutUint32(self.Data[OFFSET_MESSAGE_LEN:OFFSET_MESSAGE_LEN+4], uint32(self.MsgLength))
	return self
}
func (self *Packet) Append(data []byte) *Packet {
	return self.Write(data, len(self.Data))
}
func (self *Packet) AppendPacket(packet *Packet) *Packet {
	packet.Pack()
	return self.Write(packet.Data, len(self.Data))
}
func (self *Packet) AppendUint32(v uint32) *Packet {
	self.CheckSize(len(self.Data) + 4)
	binary.LittleEndian.PutUint32(self.Data[len(self.Data)-4:len(self.Data)], v)
	return self
}
func (self *Packet) AppendUint16(v uint16) *Packet {
	self.CheckSize(len(self.Data) + 2)
	binary.LittleEndian.PutUint16(self.Data[len(self.Data)-2:len(self.Data)], v)
	return self
}
func (self *Packet) Write(data []byte, pos int) *Packet {
	self.CheckSize(pos + len(data))
	copy(self.Data[pos:], data)
	return self
}
func (self *Packet) CheckSize(len int) {
	if len > cap(self.Data) {
		newBuf := make([]byte, len+1024) //因为要丢到pool中，所以这里扩展为比实际需要的大一些（1k），避免频繁扩展
		copy(newBuf, self.Data)
		self.Data = newBuf
	}
	self.Data = self.Data[:len] //裁切为确切的大小
}
func (self *Packet) GetMessageData() []byte {
	return self.Data[PACK_HEAD_SIZE:]
}
func (PPSM *PacketSessionServerManager) NewPacket(messageId uint32, errCode uint32) *Packet {
	//pkt := PacketPool.Get().(*Packet)
	pkt := &Packet{
		0,
		0, //暂时尚未使用。首位征用，用于表示是否有ReqId
		0,
		0,
		make([]byte, PACK_HEAD_SIZE, 128),
		false,
		0, //在包尾
		nil,
	}
	pkt.Init(messageId, errCode)
	return pkt
}
func (PPSM *PacketSessionServerManager) NewPbPacket(messageid PB.Message, errcode uint32, pbmsg proto.Message) (*Packet, error) {
	var mData []byte
	var err error
	if pbmsg != nil {
		mData, err = proto.Marshal(pbmsg)
		if err != nil {
			fmt.Println("Marshal SendPbmsg err :", messageid)
			return nil, err
		}
	}
	pkt := PPSM.NewPacket(uint32(messageid), errcode).Append(mData)
	return pkt, nil
}
