package interfaces

type SendRecv interface {
	Send([]byte) error
	Recv([]byte) error
}
