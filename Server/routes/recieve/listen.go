package recieve

type recieve struct {
	clientServer *Server
}

func (self *recieve) listen() {
	self.clientServer = NewServer("world1", Config.Listen, ProtocolType_TCP, nil)
}
