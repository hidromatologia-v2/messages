package watcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hidromatologia-v2/messages/common/connection"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/memphisdev/memphis.go"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func testWatcher(t *testing.T, w *Watcher) {
	t.Run("Valid", func(tt *testing.T) {
		prod := connection.DefaultProducer(tt)
		defer prod.Destroy()
		m := tables.RandomMessage(tables.Email)
		var buffer bytes.Buffer
		eErr := json.NewEncoder(&buffer).Encode(m)
		assert.Nil(tt, eErr)
		assert.Nil(tt, prod.Produce(buffer.Bytes(), memphis.AckWaitSec(5)))
		var message tables.Message
		t := time.NewTicker(time.Millisecond)
		defer t.Stop()
		for i := 0; i < 1000; i++ {
			if w.Controller.DB.
				Where("recipient = ?", m.Recipient).
				Where("subject = ?", m.Subject).
				Where("body = ?", m.Body).
				First(&message).Error == nil {
				break
			}
			<-t.C
		}
		assert.NotEqual(tt, uuid.Nil, message.UUID)
		assert.Equal(tt, m.Subject, message.Subject)
	})
}

func TestLogOnError(t *testing.T) {
	t.Run("Nil", func(tt *testing.T) {
		LogOnError(nil)
	})
	t.Run("Error", func(tt *testing.T) {
		LogOnError(fmt.Errorf("an error"))
	})
}

func TestWatcher(t *testing.T) {
	w := &Watcher{
		Controller:      connection.PostgresController(),
		MessageConsumer: connection.DefaultConsumer(t),
	}
	go w.Run()
	defer w.Close()
	testWatcher(t, w)
}
