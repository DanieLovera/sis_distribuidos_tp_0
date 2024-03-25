package comms

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/betmsg"
	log "github.com/sirupsen/logrus"
)

type Protocol struct {
	sendRecv common.SendRecv
}

func NewProtocol(sendRecv common.SendRecv) Protocol {
	return Protocol{sendRecv: sendRecv}
}

func (p *Protocol) SendBet(bet common.BetDto) {
	sendBetMsg := betmsg.NewSendBetMsg(&bet)
	stream, _ := sendBetMsg.Serialize()

	p.sendRecv.Send(stream)
	buff := make([]byte, len(stream))
	err := p.sendRecv.Recv(buff)

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			bet.BettingHouseId,
			err,
		)
		return
	}
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		bet.BettingHouseId,
		string(buff),
	)
}
