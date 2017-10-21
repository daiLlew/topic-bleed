package model

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
)

type Topic struct {
	Name          string   `yaml:"name"`
	ConsumerGroup string   `yaml:"consumer_group"`
	Brokers       []string `yaml:"brokers"`
}

func (t Topic) Info(message string, data log.Data) {
	log.Info(fmt.Sprintf("[Topic: %s] %s", t.Name, message), data)
}

func (t Topic) ErrorC(message string, err error, data log.Data) {
	log.ErrorC(fmt.Sprintf("[Topic: %s] %s", t.Name, message), err, data)
}

type Config struct {
	Topics []Topic `yaml:"topics"`
}
