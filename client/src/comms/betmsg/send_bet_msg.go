package betmsg

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/betting/dto"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/util"
)

const Codop uint8 = 0x00

var sequence uint32 = 0x00000000

/* Communication protocol to send a bet message is as follows:
 *
 * HEADER:
 * codop (1 byte) | payloadSize (big endian 2 bytes)
 *
 * PAYLOAD:
 * bettingHouseId (big endian 2 bytes) sequence (big endian 4 bytes) | document (big endian 4 bytes) | betnumber (big endian 4 bytes) |
 * nameSize (1 byte) | name (stream n bytes) | lastnameSize (1 byte) | lastname (stream n bytes) | birthdateSize (1 byte) | birthdate (stream n bytes)
 */
type SendBetMsg struct {
	bet           *dto.BettingDto
	sequence      uint32
	payloadSize   uint16
	nameSize      uint8
	lastnameSize  uint8
	birthdateSize uint8
}

func NewSendBetMsg(bet *dto.BettingDto) SendBetMsg {
	defer func() { sequence++ }()
	return SendBetMsg{
		bet:           bet,
		sequence:      sequence,
		payloadSize:   uint16(sizeOfPaylod(bet)),
		nameSize:      uint8(sizeOfBetName(bet)),
		lastnameSize:  uint8(sizeOfBetLastname(bet)),
		birthdateSize: uint8(sizeOfBetBirthdate(bet)),
	}
}

func sizeOfCodop() int {
	return util.SizeOfType(Codop)
}

func sizeOfSequence() int {
	return util.SizeOfField(SendBetMsg{}, "sequence")
}

func sizeOfPaylodSize() int {
	return util.SizeOfField(SendBetMsg{}, "payloadSize")
}

func sizeOfNameSize() int {
	return util.SizeOfField(SendBetMsg{}, "nameSize")
}

func sizeOfLastnameSize() int {
	return util.SizeOfField(SendBetMsg{}, "lastnameSize")
}

func sizeOfBirthdateSize() int {
	return util.SizeOfField(SendBetMsg{}, "birthdateSize")
}

func sizeOfBetName(bet *dto.BettingDto) int {
	return bet.SizeOfName()
}

func sizeOfBetLastname(bet *dto.BettingDto) int {
	return bet.SizeOfLastname()
}

func sizeOfBetBirthdate(bet *dto.BettingDto) int {
	return bet.SizeOfBirthdate()
}

func sizeOfBet(bet *dto.BettingDto) int {
	return bet.SizeOf()
}

func sizeOfPaylod(bet *dto.BettingDto) int {
	return sizeOfSequence() + sizeOfNameSize() + sizeOfLastnameSize() + sizeOfBirthdateSize() + sizeOfBet(bet)
}

func sizeOfSendBetMsg(bet *dto.BettingDto) int {
	return sizeOfCodop() + sizeOfPaylodSize() + sizeOfPaylod(bet)
}

func (s *SendBetMsg) Serialize() ([]byte, error) {
	result := make([]byte, 0, sizeOfSendBetMsg(s.bet))

	/** Header **/
	// codop
	codopBuf := []byte{byte(Codop)}
	result = append(result, codopBuf...)
	// payloadSize
	payloadSizeBuf := make([]byte, sizeOfPaylodSize())
	binary.BigEndian.PutUint16(payloadSizeBuf, s.payloadSize)
	result = append(result, payloadSizeBuf...)

	/** Payload **/
	/* Fixed Stream */
	// bettingHouseId
	bettingHouseIdBuf := make([]byte, s.bet.SizeOfBettingHouseId())
	binary.BigEndian.PutUint16(bettingHouseIdBuf, s.bet.BettingHouseId)
	result = append(result, bettingHouseIdBuf...)
	// sequence
	sequenceBuf := make([]byte, sizeOfSequence())
	binary.BigEndian.PutUint32(sequenceBuf, s.sequence)
	result = append(result, sequenceBuf...)
	// document
	documentBuf := make([]byte, s.bet.SizeOfDocument())
	binary.BigEndian.PutUint32(documentBuf, s.bet.Document)
	result = append(result, documentBuf...)
	// betnumber
	betnumberBuf := make([]byte, s.bet.SizeOfBetnumber())
	binary.BigEndian.PutUint32(betnumberBuf, s.bet.Betnumber)
	result = append(result, betnumberBuf...)

	/* Variable Stream */
	// name
	nameSize := []byte{byte(s.nameSize)}
	result = append(result, nameSize...)
	result = append(result, []byte(s.bet.Name)...)
	// lastname
	lastnameSize := []byte{byte(s.lastnameSize)}
	result = append(result, lastnameSize...)
	result = append(result, []byte(s.bet.Lastname)...)
	// birthdate
	birthdateSize := []byte{byte(s.birthdateSize)}
	result = append(result, birthdateSize...)
	result = append(result, []byte(s.bet.Birthdate)...)

	return result, nil
}
