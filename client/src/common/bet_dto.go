package common

type BetDto struct {
	Document  uint32
	Name      string
	Lastname  string
	Betnumber uint32
	Birthdate string
}

func (b *BetDto) SizeOfDocument() int {
	return SizeOfField(BetDto{}, "Document")
}

func (b *BetDto) SizeOfName() int {
	return len(b.Name)
}

func (b *BetDto) SizeOfLastname() int {
	return len(b.Lastname)
}

func (b *BetDto) SizeOfBetnumber() int {
	return SizeOfField(BetDto{}, "Betnumber")
}

func (b *BetDto) SizeOfBirthdate() int {
	return len(b.Birthdate)
}

func (b *BetDto) SizeOf() int {
	return b.SizeOfDocument() + b.SizeOfName() + b.SizeOfLastname() + b.SizeOfBetnumber() + b.SizeOfBirthdate()
}
