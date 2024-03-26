package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/comms"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/dto"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/util"

	log "github.com/sirupsen/logrus"
)

// BettingHouseConfig Configuration used by the client
type BettingHouseConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// BettingHouse Entity that encapsulates how
type BettingHouse struct {
	config BettingHouseConfig
	socket util.SocketTcp
}

// NewBettingHouse Initializes a new client receiving the configuration
// as a parameter
func NewBettingHouse(config BettingHouseConfig) *BettingHouse {
	bettingHouse := &BettingHouse{
		config: config,
	}
	return bettingHouse
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (b *BettingHouse) createSocket() error {
	b.socket = util.NewSocketTcp()
	err := b.socket.Connect(b.config.ServerAddress)
	return err
}

// Start Send messages to the client until some time threshold is met
func (b *BettingHouse) Start() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM)
	betReader := GetBettingReaderInstance()

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(b.config.LoopLapse); ; {
		bet, err := betReader.Read()
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v", b.config.ID)
			continue
		}

		err = b.createSocket()
		if err != nil {
			log.Fatalf("action: connect | result: fail | client_id: %v | error: %v", b.config.ID, err)
			continue
		}

		defer b.socket.Close()
		join := make(chan uint8, 1)
		go b.processClient(join, bet)

		// Wait until timeout, signal or join from the processClient
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				b.config.ID,
			)
			break loop
		case <-signalChannel:
			log.Infof("action: sigterm_handler | result: received | client_id: %v", b.config.ID)
			break loop
		case <-join:
		}
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", b.config.ID)
}

func (b *BettingHouse) processClient(join chan uint8, bet dto.BettingDto) {
	defer func() {
		time.Sleep(b.config.LoopPeriod)
		join <- 0
	}()

	id, err := strconv.ParseUint(b.config.ID, 10, 16)
	if err != nil {
		log.Errorf("action: parse_id | result: fail | client_id: %v | error: %v", b.config.ID, err)
		return
	}

	protocol := comms.NewProtocol(&b.socket)
	bet.BettingHouseId = uint16(id)
	betStatus, err := protocol.SendBet(bet)
	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v", bet.BettingHouseId, err)
		return
	}
	log.Infof("action: send_message | result: success | client_id: %v | msg: %v", bet.BettingHouseId, betStatus)
}
