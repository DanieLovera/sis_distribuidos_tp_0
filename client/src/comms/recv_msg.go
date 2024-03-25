package comms

type RecvMsg[T interface{}] interface {
	Deserialize() (T, error)
}
