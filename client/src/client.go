package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms"

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
	c.socket = common.NewSocketTcp()
	err := c.socket.Connect(c.config.ServerAddress)
	return err
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM)
	betReader := common.GetBetReaderInstance()

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		bet, err := betReader.Read()
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v", c.config.ID)
			continue
		}

		err = c.createSocket()
		if err != nil {
			log.Fatalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
			continue
		}

		defer c.socket.Close()
		join := make(chan uint8, 1)
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

func (c *Client) processClient(join chan uint8, bet common.BetDto) {
	defer func() {
		time.Sleep(c.config.LoopPeriod)
		join <- 0
	}()

	id, err := strconv.ParseUint(c.config.ID, 10, 16)
	if err != nil {
		log.Errorf("action: parse_id | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	protocol := comms.NewProtocol(&c.socket)
	bet.BettingHouseId = uint16(id)
	betStatus, err := protocol.SendBet(bet)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return
	}
	log.Infof("action: send_message | result: success | client_id: %v | msg: %v", bet.BettingHouseId, betStatus)
}
