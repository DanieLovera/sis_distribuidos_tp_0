package common

type SendRecv interface {
	Send([]byte) error
	Recv([]byte) error
}
