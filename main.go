package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ONSdigital/go-ns/kafka"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	topic := flag.String("topic", "", "the kafka topic to consumer messages from")
	group := flag.String("group", "", "the consumer group")
	flag.Parse()

	if len(*topic) == 0 {
		log.Error(errors.New("please specify a topic to consume from"), nil)
		os.Exit(1)
	}
	if len(*group) == 0 {
		log.Error(errors.New("please specify a consumer group"), nil)
		os.Exit(1)
	}

	kafkaConsumer, _ := kafka.NewConsumerGroup([]string{"localhost:9092"}, *topic, *group, kafka.OffsetNewest)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case message := <-kafkaConsumer.Incoming():
			log.Info("removing message from topic", nil)
			message.Commit()
		case <-sigChan:
			log.Info("shutting down...", nil)
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			if err := kafkaConsumer.Close(ctx); err != nil {
				log.Error(err, nil)
			}
			return
		}
	}
}
