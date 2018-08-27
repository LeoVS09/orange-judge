package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"orange-judge/log"
	"sync/atomic"
	"time"
)

func Read(fileName string, delay time.Duration) *atomic.Value {
	var store atomic.Value
	var config ConfigFile
	err := readSourcePermanently(fileName, &store, delay, func(data []byte) (interface{}, error) {
		err := json.Unmarshal(data, &config)
		if err != nil {
			log.Error("Error parse source file", err)
		}
		return &config, err
	})
	log.Check("Cannot read config", err)
	return &store
}

func ToConfigData(configStore *atomic.Value) (*ConfigFile, error) {
	configData, ok := configStore.Load().(*ConfigFile)
	if ok == true {
		return configData, nil
	}

	log.Error("Cannot recognise config in memory")
	return nil, errors.New("Cannot recognise config in memory")
}

func readSourcePermanently(fileName string, store *atomic.Value, delay time.Duration, unmarshal func([]byte) (interface{}, error)) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error("Error read source file", err)
		return err
	}

	value, err := unmarshal(data)
	if err != nil {
		log.Error("Cannot unmarshal data", err)
		return err
	}
	store.Store(value)
	log.InfoFmt("Source had read from file: %s", fileName)

	go func() {
		for {
			time.Sleep(delay)
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				log.WarningFmt("Error read source file\n%s", err.Error())
				continue
			}
			value, err := unmarshal(data)
			if err != nil {
				log.WarningFmt("Cannot unmarshal data\n%s", err)
				continue
			}
			store.Store(value)
		}
	}()

	return nil
}
