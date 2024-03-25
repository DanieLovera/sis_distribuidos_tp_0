package comms

type SendMsg interface {
	Serialize() ([]byte, error)
}
