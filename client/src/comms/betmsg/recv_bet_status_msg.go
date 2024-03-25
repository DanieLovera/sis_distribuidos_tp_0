package betmsg

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/common"
)

/* Communication protocol to recv a bet message status is as follows:
 *
 * status (1 byte) | document (big endian 4 bytes) | betnumber (big endian 4 bytes)
 */
type RecvBetStatusMsg struct {
	stream []byte
}

func NewRecvBetStatusMsg(stream []byte) RecvBetStatusMsg {
	return RecvBetStatusMsg{
		stream: stream,
	}
}

func (r *RecvBetStatusMsg) Deserialize() (common.BetStatusDto, error) {
	betStatus := common.BetStatusDto{}
	firstPointer := 0
	secondPointer := firstPointer + betStatus.SizeOfStatus()
	status := r.stream[firstPointer]
	firstPointer = secondPointer
	secondPointer = firstPointer + betStatus.SizeOfDocument()
	document := binary.BigEndian.Uint32(r.stream[firstPointer:secondPointer])
	firstPointer = secondPointer
	secondPointer = firstPointer + betStatus.SizeOfBetnumber()
	betnumber := binary.BigEndian.Uint32(r.stream[firstPointer:secondPointer])
	betStatus.Status = status
	betStatus.Document = document
	betStatus.Betnumber = betnumber
	return betStatus, nil
}
