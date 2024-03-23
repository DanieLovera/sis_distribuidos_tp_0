package common

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type BetReader struct {
}

var instance *BetReader = nil
var once sync.Once = sync.Once{}
var sequence uint32 = 1
var v *viper.Viper = viper.New()

func GetInstance() *BetReader {
	once.Do(func() {
		v.AutomaticEnv()
		v.SetEnvPrefix("cli")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.BindEnv("document")
		v.BindEnv("name")
		v.BindEnv("lastname")
		v.BindEnv("betnumber")
		v.BindEnv("birthdate")
		instance = &BetReader{}
	})
	return instance
}

func (b *BetReader) Read() (BetDto, error) {
	defer func() { sequence++ }()
	return BetDto{
		Sequence:  sequence,
		Document:  v.GetUint32("document"),
		Name:      v.GetString("name"),
		Lastname:  v.GetString("lastname"),
		Betnumber: v.GetUint32("betnumber"),
		Birthdate: v.GetString("birthdate"),
	}, nil
}
