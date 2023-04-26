package connection

import (
	"testing"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/memphisdev/memphis.go"
	"github.com/stretchr/testify/assert"
	"github.com/wneessen/go-mail"
)

func PostgresController() *models.Controller {
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
	return models.NewController(&opts)
}

const (
	testingStation = "testing"
)

func DefaultConsumer(t *testing.T) *memphis.Consumer {
	conn, cErr := memphis.Connect(
		"127.0.0.1",
		"root",
		memphis.Password("memphis"),
		// memphis.ConnectionToken("memphis"),
	)
	assert.Nil(t, cErr)
	c, err := conn.CreateConsumer(testingStation, random.String())
	assert.Nil(t, err)
	return c
}

func DefaultProducer(t *testing.T) *memphis.Producer {
	conn, cErr := memphis.Connect(
		"127.0.0.1",
		"root",
		memphis.Password("memphis"),
		// memphis.ConnectionToken("memphis"),
	)
	assert.Nil(t, cErr)
	prod, err := conn.CreateProducer(testingStation, random.String())
	assert.Nil(t, err)
	return prod
}
