package interfaces

type RecvMsg[T interface{}] interface {
	Deserialize([]byte) (T, error)
}
