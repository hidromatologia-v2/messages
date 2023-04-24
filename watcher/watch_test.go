package watcher

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/common/sqlite"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/memphisdev/memphis.go"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wneessen/go-mail"
)

const (
	testingStation  = "testing"
	testingConsumer = "testing-consumer"
	testingProducer = "testing-producer"
)

func testWatcher(t *testing.T, w *Watcher) {
	t.Run("Valid", func(tt *testing.T) {
		go w.Run()
		defer w.Close()
		prod, err := w.Conn.CreateProducer(testingStation, testingProducer)
		assert.Nil(tt, err)
		defer prod.Destroy()
		m := tables.RandomMessage(tables.Email)
		var buffer bytes.Buffer
		eErr := json.NewEncoder(&buffer).Encode(m)
		assert.Nil(tt, eErr)
		assert.Nil(tt, prod.Produce(buffer.Bytes(), memphis.AckWaitSec(5)))
		var message tables.Message
		for i := 0; i < 1000; i++ {
			if w.Controller.DB.
				Where("recipient = ?", m.Recipient).
				Where("subject = ?", m.Subject).
				Where("body = ?", m.Body).
				First(&message).Error == nil {
				break
			}
		}
		assert.NotEqual(tt, uuid.Nil, message.UUID)
		assert.Equal(tt, m.Subject, message.Subject)
	})
}

func TestWatcher(t *testing.T) {
	t.Run("SQLite", func(tt *testing.T) {
		opts := models.Options{
			Database:  sqlite.NewMem(),
			Cache:     cache.Bigcache(),
			JWTSecret: []byte(random.String()),
			Mail: &models.MailOptions{
				From: "sulcud@mail.com",
				Host: "127.0.0.1",
				Options: []mail.Option{
					mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
					mail.WithUsername(""), mail.WithPassword(""),
					mail.WithTLSPolicy(mail.NoTLS),
				},
			},
		}
		c := models.NewController(&opts)
		conn, cErr := memphis.Connect(
			"127.0.0.1",
			"root",
			memphis.Password("memphis"),
			// memphis.ConnectionToken("memphis"),
		)
		assert.Nil(t, cErr)
		w := &Watcher{
			Controller:   c,
			StationName:  testingStation,
			ConsumerName: testingConsumer,
			ConsumerOpts: []memphis.ConsumerOpt{},
			Conn:         conn,
		}
		testWatcher(tt, w)
	})
	t.Run("PostgreSQL", func(tt *testing.T) {
		opts := models.Options{
			Database:  postgres.NewDefault(),
			Cache:     cache.RedisDefault(),
			JWTSecret: []byte(random.String()),
			Mail: &models.MailOptions{
				From: "sulcud@mail.com",
				Host: "127.0.0.1",
				Options: []mail.Option{
					mail.WithPort(1025), mail.WithSMTPAuth(mail.SMTPAuthPlain),
					mail.WithUsername(""), mail.WithPassword(""),
					mail.WithTLSPolicy(mail.NoTLS),
				},
			},
		}
		c := models.NewController(&opts)
		conn, cErr := memphis.Connect(
			"127.0.0.1",
			"root",
			memphis.Password("memphis"),
			// memphis.ConnectionToken("memphis"),
		)
		assert.Nil(t, cErr)
		w := &Watcher{
			Controller:   c,
			StationName:  testingStation,
			ConsumerName: testingConsumer,
			ConsumerOpts: []memphis.ConsumerOpt{},
			Conn:         conn,
		}
		testWatcher(tt, w)
	})
}
