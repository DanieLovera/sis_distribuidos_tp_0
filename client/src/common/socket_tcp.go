package common

import (
	"fmt"
	"net"
)

const procotol = "tcp"

type Socket struct {
	socket net.Conn
}

func NewSocket() Socket {
	return Socket{socket: nil}
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

func (s *Socket) Connect(host string, service string) error {
	address := fmt.Sprintf("%s:%s", host, service)
	conn, err := net.Dial(procotol, address)
	if err != nil {
		return err
	}
	s.socket = conn
	return nil
}

func (s *Socket) Send(stream []byte) error {
	nBytesToWrite := len(stream)
	return sendRecv(stream, nBytesToWrite, s.socket.Write)
}

func (s *Socket) Recv(buffer []byte) error {
	nBytesToRead := cap(buffer)
	return sendRecv(buffer, nBytesToRead, s.socket.Read)
}

func (s *Socket) Close() error {
	return s.socket.Close()
}
