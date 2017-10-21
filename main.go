package main

import (
	"context"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/daiLlew/topic-bleed/model"
	"io/ioutil"
	"sync"
)

func main() {
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.ErrorC("config.yml not found", err, nil)
		os.Exit(1)
	}

	var config model.Config
	if err := yaml.Unmarshal(source, &config); err != nil {
		panic(err)
	}

	log.Info("successfully loaded config", log.Data{"config": config.Topics})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(config.Topics))

	for _, topic := range config.Topics {
		go func() {
			consumer, err := kafka.NewConsumerGroup(topic.Brokers, topic.Name, topic.ConsumerGroup, kafka.OffsetNewest)
			if err != nil {
				log.ErrorC("failed to create kafka consumer", err, log.Data{"topic": topic})
				wg.Done()
				os.Exit(1)
			}

			running := true

			for running {
				select {
				case message := <-consumer.Incoming():
					topic.Info("bleeding message from topic", nil)
					message.Commit()
				case err := <-consumer.Errors():
					topic.ErrorC("consumer errors chan received error", err, nil)
					errorChan <- err
					running = false
				case <-shutdown:
					topic.Info("received shutting down notification, attempting to close consumer", nil)

					ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
					err := consumer.Close(ctx);
					if err != nil {
						topic.ErrorC("consumer.Close returned an error", err, nil)
					} else {
						topic.Info("consumer closed successfully", nil)
					}
					running = false
				}
			}
			wg.Done()
		}()
	}

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
