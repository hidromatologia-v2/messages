package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/memphisdev/memphis.go"
)

type Watcher struct {
	Controller   *models.Controller
	StationName  string
	ConsumerName string
	ConsumerOpts []memphis.ConsumerOpt
	Conn         *memphis.Conn
	consumer     *memphis.Consumer
}

func (w *Watcher) Close() error {
	w.consumer.StopConsume()
	w.Conn.Close()
	return w.Controller.Close()
}

func (w *Watcher) Run() error {
	var cErr error
	w.consumer, cErr = w.Conn.CreateConsumer(
		w.StationName,
		w.ConsumerName,
		w.ConsumerOpts...,
	)
	if cErr != nil {
		return fmt.Errorf("consumer creation error: %w", cErr)
	}
	return w.consumer.Consume(
		func(messages []*memphis.Msg, err error, ctx context.Context) {
			if err != nil {
				log.Println("consumer err:", err)
				return
			}
			for _, message := range messages {
				var m tables.Message
				dErr := json.NewDecoder(bytes.NewReader(message.Data())).Decode(&m)
				if dErr != nil {
					log.Println("error while decoding the message:", dErr)
					continue
				}
				sendErr := w.Controller.SendMessage(&m)
				if sendErr != nil {
					log.Println("error while sending the message:", sendErr)
					continue
				}
				ackErr := message.Ack()
				if ackErr != nil {
					log.Println("error while ack the message:", ackErr)
				}
			}
		},
	)
}
