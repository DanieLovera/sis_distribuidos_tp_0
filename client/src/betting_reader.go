package main

import (
	"strings"
	"sync"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/src/dto"
	"github.com/spf13/viper"
)

type BetingReader struct {
	vpr *viper.Viper
}

var instance *BetingReader = nil
var syncOnce sync.Once = sync.Once{}

func GetBettingReaderInstance() *BetingReader {
	syncOnce.Do(func() {
		var vpr *viper.Viper = viper.New()
		vpr.AutomaticEnv()
		vpr.SetEnvPrefix("cli")
		vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		vpr.BindEnv("document")
		vpr.BindEnv("name")
		vpr.BindEnv("lastname")
		vpr.BindEnv("betnumber")
		vpr.BindEnv("birthdate")
		instance = &BetingReader{
			vpr: vpr,
		}
	})
	return instance
}

func (b *BetingReader) Read() (dto.BettingDto, error) {
	return dto.BettingDto{
		Document:  b.vpr.GetUint32("document"),
		Name:      b.vpr.GetString("name"),
		Lastname:  b.vpr.GetString("lastname"),
		Betnumber: b.vpr.GetUint32("betnumber"),
		Birthdate: b.vpr.GetString("birthdate"),
	}, nil
}
