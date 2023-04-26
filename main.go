package main

import (
	"context"
	"log"

	"github.com/hidromatologia-v2/messages/watcher"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/memphisdev/memphis.go"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
	"github.com/wneessen/go-mail"
)

func main() {
	var config Config
	eErr := envconfig.Process(context.Background(), &config)
	if eErr != nil {
		log.Fatal(eErr)
	}
	mailOpts := []mail.Option{
		mail.WithPort(config.SMTP.Port),
	}
	if config.SMTP.Username != nil && config.SMTP.Password != nil {
		mailOpts = append(mailOpts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
		mailOpts = append(mailOpts, mail.WithUsername(*config.SMTP.Username))
		mailOpts = append(mailOpts, mail.WithUsername(*config.SMTP.Password))
	}
	if config.SMTP.NoTLS != nil {
		mailOpts = append(mailOpts, mail.WithTLSPolicy(mail.NoTLS))
	}
	controllerOpts := models.Options{
		Database: postgres.New(config.Postgres.DSN),
		Cache: cache.Redis(&redis.Options{
			Addr: config.Redis.Addr,
			DB:   config.Redis.DB,
		}),
		Mail: &models.MailOptions{
			From:    config.SMTP.From,
			Host:    config.SMTP.Host,
			Options: mailOpts,
		},
	}
	var connOpts []memphis.Option
	if config.Consumer.Password != nil {
		connOpts = append(connOpts, memphis.Password(*config.Consumer.Password))
	}
	if config.Consumer.ConnectionToken != nil {
		connOpts = append(connOpts, memphis.ConnectionToken(*config.Consumer.ConnectionToken))
	}
	conn, err := memphis.Connect(
		config.Consumer.Host,
		config.Consumer.Username,
		connOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
	consumer, cErr := conn.CreateConsumer(
		config.Consumer.Station,
		config.Consumer.Consumer,
	)
	if cErr != nil {
		log.Fatal(cErr)
	}
	w := &watcher.Watcher{
		Controller:      models.NewController(&controllerOpts),
		MessageConsumer: consumer,
	}
	rErr := w.Run()
	if rErr != nil {
		log.Fatal(rErr)
	}
	<-make(chan struct{}, 1)
}
