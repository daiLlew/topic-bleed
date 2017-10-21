package model

import (
	"github.com/ONSdigital/go-ns/log"
	"encoding/json"
	"errors"
)

type Topic struct {
	Name          string   `yaml:"name"`
	ConsumerGroup string   `yaml:"consumer_group"`
	Brokers       []string `yaml:"brokers"`
}

func (t Topic) Validate() error {
	if len(t.Name) == 0 {
		return errors.New("topic name is empty")
	}
	if len(t.ConsumerGroup) == 0 {
		return errors.New("consumer group is empty")
	}
	if len(t.Brokers) == 0 {
		return errors.New("brokers is empty")
	}
	return nil
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

func (c Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c Config) Validate() error {
	if len(c.Topics) == 0 {
		return errors.New("no topics found")
	}
	for i, topic := range c.Topics {
		if err := topic.Validate(); err != nil {
			log.ErrorC("invalid topic", err, log.Data{"index": i, "topic": topic})
			return err
		}
	}
	return nil
}
