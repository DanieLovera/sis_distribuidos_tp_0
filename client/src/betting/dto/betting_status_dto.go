package dto

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/util"

type BettingStatusDto struct {
	Status    uint8
	Document  uint32
	Betnumber uint32
}

func (b *BettingStatusDto) SizeOfStatus() int {
	return util.SizeOfField(BettingStatusDto{}, "Status")
}

func (b *BettingStatusDto) SizeOfDocument() int {
	return util.SizeOfField(BettingStatusDto{}, "Document")
}

func (b *BettingStatusDto) SizeOfBetnumber() int {
	return util.SizeOfField(BettingStatusDto{}, "Betnumber")
}
