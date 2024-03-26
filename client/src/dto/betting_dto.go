package dto

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/util"

type BettingDto struct {
	BettingHouseId uint16
	Document       uint32
	Name           string
	Lastname       string
	Betnumber      uint32
	Birthdate      string
}

func (b *BettingDto) SizeOfBettingHouseId() int {
	return util.SizeOfField(BettingDto{}, "BettingHouseId")
}

func (b *BettingDto) SizeOfDocument() int {
	return util.SizeOfField(BettingDto{}, "Document")
}

func (b *BettingDto) SizeOfBetnumber() int {
	return util.SizeOfField(BettingDto{}, "Betnumber")
}

func (b *BettingDto) SizeOfName() int {
	return len(b.Name)
}

func (b *BettingDto) SizeOfLastname() int {
	return len(b.Lastname)
}

func (b *BettingDto) SizeOfBirthdate() int {
	return len(b.Birthdate)
}

func (b *BettingDto) SizeOf() int {
	return b.SizeOfBettingHouseId() + b.SizeOfDocument() + b.SizeOfName() + b.SizeOfLastname() + b.SizeOfBetnumber() + b.SizeOfBirthdate()
}
