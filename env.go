package main

import "github.com/hidromatologia-v2/models/common/config"

type (
	Config struct {
		config.Consumer `env:",prefix=MEMPHIS_"`  // Memphis
		config.Postgres `env:",prefix=POSTGRES_"` // POSTGRESQL
		config.Redis    `env:",prefix=REDIS_"`    // Redis
		config.SMTP     `env:",prefix=SMTP_"`     // SMTP
	}
)
