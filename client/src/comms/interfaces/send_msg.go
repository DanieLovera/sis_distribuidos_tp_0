package interfaces

type SendMsg interface {
	Serialize() ([]byte, error)
}
