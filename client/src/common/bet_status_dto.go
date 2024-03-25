package common

type BetStatusDto struct {
	Status    uint8
	Document  uint32
	Betnumber uint32
}

func (b *BetStatusDto) SizeOfStatus() int {
	return SizeOfField(BetStatusDto{}, "Status")
}

func (b *BetStatusDto) SizeOfDocument() int {
	return SizeOfField(BetStatusDto{}, "Document")
}

func (b *BetStatusDto) SizeOfBetnumber() int {
	return SizeOfField(BetStatusDto{}, "Betnumber")
}
