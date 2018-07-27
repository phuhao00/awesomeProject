package Definition

import "github.com/golang/protobuf/proto"
import "../../PB"

type PacketSessionServerManager struct {
	Handlers map[int32]func()
}

var PSSM = new(PacketSessionServerManager)

func (self *PacketSessionServerManager) initPSSM() {
	PSSM.Handlers = make(map[int32]func())
	PSSM.InitServerHandler()
	PSSM.initPacketHandler()
	PSSM.InitSessionHandler()
}

type SyncPBRequest struct {
	Id              uint32
	MsgId           PB.Message
	Errcode         uint32
	Request         proto.Message
	Response        proto.Message
	ResponseErrCode int32
	Ch_response     chan *Packet
	Time            int64
}

type AsyncPBRequest struct {
	Id      uint32
	MsgId   PB.Message
	Errcode uint32
	Request proto.Message
}

type ClientPacketHandler func(clientpacket *ClientPacket)
