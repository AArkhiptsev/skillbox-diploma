package conf

import (
	"diploma/lib"
	"github.com/naoina/toml"
	"os"
)

const ConfigFileName = "config.toml"

type ConfigT struct {
	Files struct {
		SmsFileName,
		VoiceCallFileName,
		EmailFileName,
		BillingFileName string
	}
	FetchServers struct {
		MmsDataServer,
		SupportServer,
		AccidentServer string
	}

	Server struct {
		Serveraddr string
	}
}

var Config ConfigT

func ReadConfig(filename string) {

	lib.LogParseErr(1, "Чтение файла конфигурации: "+ConfigFileName)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := toml.NewDecoder(f).Decode(&Config); err != nil {
		panic(err)
	}
	lib.LogParseErr(1, "Параметры прочитаны успешно")

}
