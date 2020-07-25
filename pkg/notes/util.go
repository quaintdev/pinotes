package notes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	configFileName = "./conf.json"
	config         *Conf
	once           sync.Once
)

type Conf struct {
	Path         string
	WebNotesFile string
	HttpPort     string
	InterfaceIP  string
	GrpcPort     string
}

func GetConfig() *Conf {
	once.Do(func() {
		log.Println("loading config..")
		configBytes, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Printf("config not present. err: %s exiting.", err)
			panic(err)
		}
		config = new(Conf)
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			log.Println("invalid config present. err", err)
			panic(err)
		}
	})
	return config
}
