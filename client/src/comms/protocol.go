package comms

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/betting/dto"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/betmsg"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/interfaces"
	log "github.com/sirupsen/logrus"
)

type Protocol struct {
	sendRecv interfaces.SendRecv
}

func NewProtocol(sendRecv interfaces.SendRecv) Protocol {
	return Protocol{sendRecv: sendRecv}
}

func (p *Protocol) SendBet(bet dto.BettingDto) (dto.BettingStatusDto, error) {
	err := p.sendBet(bet)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return dto.BettingStatusDto{}, err
	}

	betStatus, err := p.recvBetStatus()
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return dto.BettingStatusDto{}, err
	}
	return betStatus, nil
}

func (p *Protocol) sendBet(bet dto.BettingDto) error {
	sendBetMsg := betmsg.NewSendBetMsg(&bet)
	stream, _ := sendBetMsg.Serialize()
	return p.sendRecv.Send(stream)
}

func (p *Protocol) recvBetStatus() (dto.BettingStatusDto, error) {
	recvBetStatusMsg := betmsg.NewRecvBetStatusMsg()
	buff := make([]byte, recvBetStatusMsg.SizeOfPayloadSize())
	err := p.sendRecv.Recv(buff)
	if err != nil {
		return dto.BettingStatusDto{}, err
	}

	payloadSize := binary.BigEndian.Uint32(buff)
	buff = make([]byte, payloadSize)
	err = p.sendRecv.Recv(buff)
	if err != nil {
		return dto.BettingStatusDto{}, err
	}
	return recvBetStatusMsg.Deserialize(buff)
}
