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
	Controller      *models.Controller
	MessageConsumer *memphis.Consumer
}

func (w *Watcher) Close() error {
	return w.MessageConsumer.Destroy()
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

func LogOnError(err error) {
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		log.Print(err)
	}
}

func (w *Watcher) Run() error {
	return w.MessageConsumer.Consume(
		func(messages []*memphis.Msg, err error, ctx context.Context) {
			LogOnError(err)
			for _, message := range messages {
				LogOnError(w.handleMessage(message))
			}
		},
	)
}
