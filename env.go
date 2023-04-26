package main

type (
	Consumer struct {
		StationName     string  `env:"STATION,required"`
		ConsumerName    string  `env:"CONSUMER,required"`
		Host            string  `env:"HOST,required"`
		Username        string  `env:"USERNAME,required"`
		Password        *string `env:"PASSWORD,noinit"`
		ConnectionToken *string `env:"CONN_TOKEN,noinit"`
	}
	Redis struct {
		Addr string `env:"ADDR,required"`
		DB   int    `env:"DB,required"`
	}
	SMTP struct {
		From     string  `env:"FROM,required"`
		Host     string  `env:"HOST,required"`
		Port     int     `env:"PORT,required"`
		Username *string `env:"USERNAME,noinit"`
		Password *string `env:"PASSWORD,noinit"`
		NoTLS    *bool   `env:"NO_TLS,noinit"`
	}
	Postgres struct {
		DSN string `env:"DSN,required"`
	}
	Config struct {
		Consumer `env:",prefix=MEMPHIS_"`  // Memphis
		Postgres `env:",prefix=POSTGRES_"` // POSTGRESQL
		Redis    `env:",prefix=REDIS_"`    // Redis
		SMTP     `env:",prefix=SMTP_"`     // SMTP
	}
)
