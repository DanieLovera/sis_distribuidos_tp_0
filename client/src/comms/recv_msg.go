package comms

type RecvMsg[T any] interface {
	Deserialize() (T, error)
}
