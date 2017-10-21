package main

import (
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"os/signal"
	"syscall"
	"gopkg.in/yaml.v2"
	"github.com/daiLlew/topic-bleed/model"
	"io/ioutil"
	"sync"
	"context"
	"time"
	"flag"
	"path/filepath"
	"errors"
)

var wg sync.WaitGroup

func main() {
	filename := flag.String("config", "config.yml", "The YAML configuration specifying the topics to bleed")
	flag.Parse()

	if filepath.Ext(*filename) != ".yml" {
		log.Error(errors.New("cannot load configuration: config mut be a yml file"), nil)
		os.Exit(1)
	}

	config := loadConfig(*filename)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	bleed(config, shutdown, errorChan)

	select {
	case sig := <-sigChan:
		log.Info("os signal received shutting down", log.Data{"signal": sig.String()})
		close(shutdown)
	case err := <-errorChan:
		log.ErrorC("error chan received error shutting down", err, nil)
		close(shutdown)
	}

	wg.Wait()
	log.Info("shutdown complete", nil)
}

func bleed(config model.Config, shutdown chan struct{}, errorChan chan error) {
	for _, topic := range config.Topics {
		consumer, err := kafka.NewConsumerGroup(topic.Brokers, topic.Name, topic.ConsumerGroup, kafka.OffsetNewest)
		if err != nil {
			topic.ErrorC("failed to create kafka consumer", err, nil)
			wg.Done()
			errorChan <- err
			break
		}

		wg.Add(1)

		go func(t model.Topic) {
			running := true
			for running {
				select {
				case message := <-consumer.Incoming():
					t.Info("bleeding message from topic", nil)
					message.Commit()
				case err := <-consumer.Errors():
					t.ErrorC("consumer errors chan received error", err, nil)
					closeConsumer(t, consumer)
					running = false
					errorChan <- err
				case <-shutdown:
					t.Info("received shutting down notification, attempting to close consumer", nil)
					closeConsumer(t, consumer)
					running = false
				}
			}

			wg.Done()
			t.Info("exiting bleeder", nil)
		}(topic)
	}
}

func loadConfig(filename string) model.Config {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.ErrorC("config file not found", err, log.Data{"filename": filename})
		os.Exit(1)
	}

	log.Info("loading configuration file", log.Data{"filename": filename})
	var config model.Config
	if err := yaml.Unmarshal(source, &config); err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	if err := config.Validate(); err != nil {
		os.Exit(1)
	}

	log.Info("successfully loaded config", log.Data{"config": config.String()})
	return config
}

func closeConsumer(t model.Topic, consumer *kafka.ConsumerGroup) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	err := consumer.Close(ctx)
	if err != nil {
		t.ErrorC("consumer.Close returned an error", err, nil)
	} else {
		t.Info("consumer closed successfully", nil)
	}
}
