package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

func (w *Watcher) handleMessage(message *memphis.Msg) error {
	var m tables.Message
	dErr := json.NewDecoder(bytes.NewReader(message.Data())).Decode(&m)
	if dErr != nil {
		return fmt.Errorf("error while decoding the message: %w", dErr)
	}
	sendErr := w.Controller.SendMessage(&m)
	if sendErr != nil {
		return fmt.Errorf("error while sending the message: %w", sendErr)
	}
	ackErr := message.Ack()
	if ackErr != nil {
		return fmt.Errorf("error while ack the message: %w", ackErr)
	}
	return nil
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
			if err != nil && !strings.Contains(err.Error(), "timeout") {
				log.Println("consumer err:", err)
				return
			}
			for _, message := range messages {
				hErr := w.handleMessage(message)
				if hErr != nil {
					log.Print(hErr)
				}
			}
		},
	)
}
