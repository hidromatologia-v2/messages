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

type Messages struct {
	MemphisStationName     string  `env:"MEMPHIS_STATION,required"`  // MEMPHIS
	MemphisConsumerName    string  `env:"MEMPHIS_CONSUMER,required"` //
	MemphisHost            string  `env:"MEMPHIS_HOST,required"`     //
	MemphisUsername        string  `env:"MEMPHIS_USERNAME,required"` //
	MemphisPassword        *string `env:"MEMPHIS_PASSWORD,noinit"`   //
	MemphisConnectionToken *string `env:"MEMPHIS_CONN_TOKEN,noinit"` //
	PostgresDsn            string  `env:"POSTGRES_DSN,required"`     // POSTGRESQL
	RedisAddr              string  `env:"REDIS_ADDR,required"`       // REDIS
	RedisDB                int     `env:"REDIS_DB,required"`         //
	VonageSecret           string  `env:"VONAGE_SECRET,required"`    // VONAGE
	VonageAPIKey           string  `env:"VONAGE_APIKEY,required"`    //
	SMTPFrom               string  `env:"SMTP_FROM,required"`        // SMTP
	SMTPHost               string  `env:"SMTP_HOST,required"`        //
	SMTPPort               int     `env:"SMTP_PORT,required"`        //
	SMTPUsername           *string `env:"SMTP_USERNAME,noinit"`      //
	SMTPPassword           *string `env:"SMTP_PASSWORD,noinit"`      //
	SMTPNoTLS              *bool   `env:"SMTP_NO_TLS,noinit"`        //
}

func main() {
	var config Messages
	eErr := envconfig.Process(context.Background(), &config)
	if eErr != nil {
		log.Fatal(eErr)
	}
	mailOpts := []mail.Option{
		mail.WithPort(config.SMTPPort),
	}
	if config.SMTPUsername != nil && config.SMTPPassword != nil {
		mailOpts = append(mailOpts, mail.WithSMTPAuth(mail.SMTPAuthPlain))
		mailOpts = append(mailOpts, mail.WithUsername(*config.SMTPUsername))
		mailOpts = append(mailOpts, mail.WithUsername(*config.SMTPPassword))
	}
	if config.SMTPNoTLS != nil {
		mailOpts = append(mailOpts, mail.WithTLSPolicy(mail.NoTLS))
	}
	controllerOpts := models.Options{
		Database: postgres.New(config.PostgresDsn),
		Cache: cache.Redis(&redis.Options{
			Addr: config.RedisAddr,
			DB:   config.RedisDB,
		}),
		Vonage: &models.VonageOptions{
			Secret: config.VonageSecret,
			APIKey: config.VonageAPIKey,
		},
		Mail: &models.MailOptions{
			From:    config.SMTPFrom,
			Host:    config.SMTPHost,
			Options: mailOpts,
		},
	}
	var connOpts []memphis.Option
	if config.MemphisPassword != nil {
		connOpts = append(connOpts, memphis.Password(*config.MemphisPassword))
	}
	if config.MemphisConnectionToken != nil {
		connOpts = append(connOpts, memphis.ConnectionToken(*config.MemphisConnectionToken))
	}
	conn, err := memphis.Connect(
		config.MemphisHost,
		config.MemphisUsername,
		connOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
	w := &watcher.Watcher{
		Controller:   models.NewController(&controllerOpts),
		StationName:  config.MemphisStationName,
		ConsumerName: config.MemphisConsumerName,
		Conn:         conn,
	}
	rErr := w.Run()
	if rErr != nil {
		log.Fatal(rErr)
	}
	<-make(chan struct{}, 1)
}
