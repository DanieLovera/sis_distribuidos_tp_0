package common

import (
	"net"
)

const procotol = "tcp"

type SocketTcp struct {
	socket net.Conn
}

func NewSocketTcp() SocketTcp {
	return SocketTcp{socket: nil}
}

func sendRecv(streamBuff []byte, streamBuffSize int, callback func([]byte) (int, error)) error {
	for nbytes, err := callback(streamBuff); nbytes < streamBuffSize; {
		if err != nil {
			return err
		}
		n := 0
		n, err = callback(streamBuff[nbytes:])
		nbytes += n
	}
	return nil
}

func (s *SocketTcp) Connect(address string) error {
	conn, err := net.Dial(procotol, address)
	if err != nil {
		return err
	}
	s.socket = conn
	return nil
}

func (s *SocketTcp) Send(stream []byte) error {
	nBytesToWrite := len(stream)
	return sendRecv(stream, nBytesToWrite, s.socket.Write)
}

func (s *SocketTcp) Recv(buffer []byte) error {
	nBytesToRead := cap(buffer)
	return sendRecv(buffer, nBytesToRead, s.socket.Read)
}

func (s *SocketTcp) Close() error {
	if s.socket == nil {
		return nil
	}
	return s.socket.Close()
}
