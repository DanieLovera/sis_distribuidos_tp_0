package comms

type RecvMsg[T interface{}] interface {
	Deserialize([]byte) (T, error)
}
