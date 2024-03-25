package common

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type BetReader struct {
	vpr *viper.Viper
}

var instance *BetReader = nil
var syncOnce sync.Once = sync.Once{}

func GetBetReaderInstance() *BetReader {
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
		instance = &BetReader{
			vpr: vpr,
		}
	})
	return instance
}

func (b *BetReader) Read() (BetDto, error) {
	return BetDto{
		Document:  b.vpr.GetUint32("document"),
		Name:      b.vpr.GetString("name"),
		Lastname:  b.vpr.GetString("lastname"),
		Betnumber: b.vpr.GetUint32("betnumber"),
		Birthdate: b.vpr.GetString("birthdate"),
	}, nil
}
