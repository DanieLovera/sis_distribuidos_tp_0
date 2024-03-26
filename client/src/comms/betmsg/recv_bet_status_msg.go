package betmsg

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/betting/dto"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/util"
)

/* Communication protocol to recv a bet message status is as follows:
 *
 * payloadSize(big endian 2 bytes) | status (1 byte) | document (big endian 4 bytes) | betnumber (big endian 4 bytes)
 */
type RecvBetStatusMsg struct{}

func NewRecvBetStatusMsg() RecvBetStatusMsg {
	return RecvBetStatusMsg{}
}

func (r *RecvBetStatusMsg) SizeOfPayloadSize() int {
	return util.SizeOfType(uint16(0x0000))
}

func (r *RecvBetStatusMsg) Deserialize(stream []byte) (dto.BettingStatusDto, error) {
	betStatus := dto.BettingStatusDto{}
	firstPointer := 0
	secondPointer := firstPointer + betStatus.SizeOfStatus()
	status := stream[firstPointer]
	firstPointer = secondPointer
	secondPointer = firstPointer + betStatus.SizeOfDocument()
	document := binary.BigEndian.Uint32(stream[firstPointer:secondPointer])
	firstPointer = secondPointer
	secondPointer = firstPointer + betStatus.SizeOfBetnumber()
	betnumber := binary.BigEndian.Uint32(stream[firstPointer:secondPointer])
	betStatus.Status = status
	betStatus.Document = document
	betStatus.Betnumber = betnumber
	return betStatus, nil
}
