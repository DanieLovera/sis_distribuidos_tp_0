package comms

import (
	"encoding/binary"

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

func (p *Protocol) SendBet(bet common.BetDto) (common.BetStatusDto, error) {
	err := p.sendBet(bet)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return common.BetStatusDto{}, err
	}

	betStatus, err := p.recvBetStatus()
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return common.BetStatusDto{}, err
	}
	return betStatus, nil
}

func (p *Protocol) sendBet(bet common.BetDto) error {
	sendBetMsg := betmsg.NewSendBetMsg(&bet)
	stream, _ := sendBetMsg.Serialize()
	return p.sendRecv.Send(stream)
}

func (p *Protocol) recvBetStatus() (common.BetStatusDto, error) {
	recvBetStatusMsg := betmsg.NewRecvBetStatusMsg()
	buff := make([]byte, recvBetStatusMsg.SizeOfPayloadSize())
	err := p.sendRecv.Recv(buff)
	if err != nil {
		return common.BetStatusDto{}, err
	}

	payloadSize := binary.BigEndian.Uint32(buff)
	buff = make([]byte, payloadSize)
	err = p.sendRecv.Recv(buff)
	if err != nil {
		return common.BetStatusDto{}, err
	}
	return recvBetStatusMsg.Deserialize(buff)
}
