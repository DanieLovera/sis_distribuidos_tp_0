package betmsg

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/common"
)

const Codop uint8 = 0x00

var sequence uint32 = 0x00000000

/* Communication protocol to send a bet message is as follows:
 *
 * HEADER:
 * codop (1 byte) | payloadSize (big endian 2 bytes)
 *
 * PAYLOAD:
 * sequence (big endian4 bytes) | document (big endian 4 bytes) | betnumber (big endian 4 bytes) |
 * nameSize (1 byte) | name (stream n bytes) | lastnameSize (1 byte) | lastname (stream n bytes) | birthdateSize (1 byte) | birthdate (stream n bytes)
 */
type SendBetMsg struct {
	payloadSize   uint16
	nameSize      uint8
	lastnameSize  uint8
	birthdateSize uint8
	sequence      uint32
	bet           *common.BetDto
}

func NewSendBetMsg(bet *common.BetDto) SendBetMsg {
	defer func() { sequence++ }()
	return SendBetMsg{
		payloadSize:   uint16(sizeOfPaylod(bet)),
		nameSize:      uint8(sizeOfName(bet)),
		lastnameSize:  uint8(sizeOfLastname(bet)),
		birthdateSize: uint8(sizeOfBirthdate(bet)),
		sequence:      sequence,
		bet:           bet,
	}
}

func sizeOfCodop() int {
	return common.SizeOfType(Codop)
}

func sizeOfSequence() int {
	return common.SizeOfField(SendBetMsg{}, "sequence")
}

func sizeOfPaylodSize() int {
	return common.SizeOfField(SendBetMsg{}, "payloadSize")
}

func sizeOfNameSize() int {
	return common.SizeOfField(SendBetMsg{}, "nameSize")
}

func sizeOfLastnameSize() int {
	return common.SizeOfField(SendBetMsg{}, "lastnameSize")
}

func sizeOfBirthdateSize() int {
	return common.SizeOfField(SendBetMsg{}, "birthdateSize")
}

func sizeOfName(bet *common.BetDto) int {
	return bet.SizeOfName()
}

func sizeOfLastname(bet *common.BetDto) int {
	return bet.SizeOfLastname()
}

func sizeOfBirthdate(bet *common.BetDto) int {
	return bet.SizeOfBirthdate()
}

func sizeOfBet(bet *common.BetDto) int {
	return bet.SizeOf()
}

func sizeOfPaylod(bet *common.BetDto) int {
	return sizeOfSequence() + sizeOfNameSize() + sizeOfLastnameSize() + sizeOfBirthdateSize() + sizeOfBet(bet)
}

func sizeOfSendBetMsg(bet *common.BetDto) int {
	return sizeOfCodop() + sizeOfPaylodSize() + sizeOfPaylod(bet)
}

func (s *SendBetMsg) Serialize() ([]byte, error) {
	result := make([]byte, 0, sizeOfSendBetMsg(s.bet))
	codopBuf := []byte{byte(Codop)}
	payloadSizeBuf := make([]byte, sizeOfPaylodSize())
	sequenceBuf := make([]byte, sizeOfSequence())
	documentBuf := make([]byte, s.bet.SizeOfDocument())
	betnumberBuf := make([]byte, s.bet.SizeOfBetnumber())
	nameSize := []byte{byte(s.nameSize)}
	lastnameSize := []byte{byte(s.lastnameSize)}
	birthdateSize := []byte{byte(s.birthdateSize)}
	binary.BigEndian.PutUint16(payloadSizeBuf, uint16(s.payloadSize))
	binary.BigEndian.PutUint32(sequenceBuf, uint32(s.sequence))
	binary.BigEndian.PutUint32(documentBuf, uint32(s.bet.Document))
	binary.BigEndian.PutUint32(betnumberBuf, uint32(s.bet.Betnumber))

	result = append(result, codopBuf...)
	result = append(result, payloadSizeBuf...)
	result = append(result, sequenceBuf...)
	result = append(result, documentBuf...)
	result = append(result, betnumberBuf...)
	result = append(result, nameSize...)
	result = append(result, []byte(s.bet.Name)...)
	result = append(result, lastnameSize...)
	result = append(result, []byte(s.bet.Lastname)...)
	result = append(result, birthdateSize...)
	result = append(result, []byte(s.bet.Birthdate)...)
	return result, nil
}
