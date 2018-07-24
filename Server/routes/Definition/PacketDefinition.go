package Definition

type Packet struct {
	MsgLength int
	Version   uint16 //暂时尚未使用。首位征用，用于表示是否有ReqId
	MessageId uint32
	ErrCode   uint32
	Data      []byte
	HasReqId  bool
	ReqId     uint32 //在包尾
}

const (
	PACK_HEAD_SIZE     = 14
	PACKVERSION        = 1
	OFFSET_MESSAGE_LEN = 0
	OFFSET_VERSION     = 4
	OFFSET_MESSAGE_ID  = 6
	OFFSET_ERRCODE     = 10
	MULTI_PACKET_ID    = 1
)

type ClientPacket struct {
	Session *Session
	Packet  *Packet
}

const (
	NewPacket_   int32 = 1
	NewPbPacket_ int32 = 2
)
