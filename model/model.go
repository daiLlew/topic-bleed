package model

import (
	"github.com/ONSdigital/go-ns/log"
)

type Topic struct {
	Name          string   `yaml:"name"`
	ConsumerGroup string   `yaml:"consumer_group"`
	Brokers       []string `yaml:"brokers"`
}

func (t Topic) Info(message string, data log.Data) {
	if data == nil {
		data = log.Data{}
	}
	data["topic"] = t.Name
	log.Info(message, data)
}

func (t Topic) ErrorC(message string, err error, data log.Data) {
	if data == nil {
		data = log.Data{}
	}
	data["topic"] = t.Name
	log.ErrorC(message, err, data)
}

type Config struct {
	Topics []Topic `yaml:"topics"`
}
