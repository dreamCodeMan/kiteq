package handler

import (
	"github.com/blackbeans/kiteq-common/protocol"
	log "github.com/blackbeans/log4go"
	client "github.com/blackbeans/turbo/client"
	packet "github.com/blackbeans/turbo/packet"
	p "github.com/blackbeans/turbo/pipe"
)

//----------------鉴权handler
type ValidateHandler struct {
	p.BaseForwardHandler
	clientManager *client.ClientManager
}

//------创建鉴权handler
func NewValidateHandler(name string, clientManager *client.ClientManager) *ValidateHandler {
	ahandler := &ValidateHandler{}
	ahandler.BaseForwardHandler = p.NewBaseForwardHandler(name, ahandler)
	ahandler.clientManager = clientManager
	return ahandler
}

func (self *ValidateHandler) TypeAssert(event p.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *ValidateHandler) cast(event p.IEvent) (val iauth, ok bool) {
	val, ok = event.(iauth)
	return val, ok
}

func (self *ValidateHandler) Process(ctx *p.DefaultPipelineContext, event p.IEvent) error {

	aevent, ok := self.cast(event)
	if !ok {
		return p.ERROR_INVALID_EVENT_TYPE
	}

	remoteClient := aevent.getClient()
	//做权限校验.............
	isAuth := self.clientManager.Validate(remoteClient)
	// log.DebugLog("kite_handler",  "ValidateHandler|CONNETION|%s|%s\n", remoteClient.RemoteAddr(), isAuth)
	if isAuth {
		ctx.SendForward(event)
	} else {
		log.Warn("ValidateHandler|UnAuth CONNETION|%s\n", remoteClient.RemoteAddr())
		cmd := protocol.MarshalConnAuthAck(false, "未授权的访问,连接关闭!")
		//响应包
		p := packet.NewPacket(protocol.CMD_CONN_AUTH, cmd)

		//直接写出去授权失败
		remoteClient.Write(*p)
		//断开连接
		remoteClient.Shutdown()
	}

	return nil

}