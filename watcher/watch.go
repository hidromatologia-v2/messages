package watcher

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/logs"
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

func (w *Watcher) HandleMessage(message *memphis.Msg) {
	defer func() {
		err, _ := recover().(error)
		if err == nil {
			return
		}
		logs.LogOnError(err)
	}()
	var m tables.Message
	dErr := json.NewDecoder(bytes.NewReader(message.Data())).Decode(&m)
	logs.PanicOnError(dErr)
	sendErr := w.Controller.SendMessage(&m)
	logs.PanicOnError(sendErr)
	ackErr := message.Ack()
	logs.PanicOnError(ackErr)
}

func (w *Watcher) Run() error {
	return w.MessageConsumer.Consume(
		func(messages []*memphis.Msg, err error, ctx context.Context) {
			logs.LogOnError(err)
			for _, message := range messages {
				w.HandleMessage(message)
			}
		},
	)
}
