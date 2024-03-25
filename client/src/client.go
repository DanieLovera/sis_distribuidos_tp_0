package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms/betmsg"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	socket common.SocketTcp
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createSocket() error {
	socket := common.NewSocketTcp()
	err := socket.Connect(c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.socket = socket
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM)
	betReader := common.GetBetReaderInstance()

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		bet, err := betReader.Read()
		if err != nil {
			// log.Errorf("action: read_bet | result: fail | client_id: %v | bet_sequence: %v", c.config.ID, bet.Sequence)
			log.Errorf("action: read_bet | result: fail | client_id: %v", c.config.ID)
			continue
		}

		err = c.createSocket()
		if err != nil {
			log.Fatalf(
				"action: connect | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
		defer c.socket.Close()
		// Process the client in a goroutine to avoid blocking operations
		join := make(chan uint32, 1)
		go c.processClient(join, bet)
		// Wait until timeout, signal or join from the processClient
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		case <-signalChannel:
			log.Infof("action: sigterm_handler | result: received | client_id: %v", c.config.ID)
			break loop
		case <-join:
		}
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) processClient(join chan uint32, bet common.BetDto) {
	bms := betmsg.NewSendBetMsg(&bet)
	se, _ := bms.Serialize()

	stream := []byte(fmt.Sprintf("[CLIENT %v] Message N°%v\n", c.config.ID, se))
	c.socket.Send(stream)
	buff := make([]byte, len(stream))
	err := c.socket.Recv(buff)

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		string(buff),
	)
	// Wait a time between sending one message and the next one
	time.Sleep(c.config.LoopPeriod)
	join <- bet.Document
}
