package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/memphisdev/memphis.go"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "messaging",
		Usage: "messaging microservice for the ResupplyOrg project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Value: "localhost",
				Usage: "Memphis host",
			},
			&cli.StringFlag{
				Name:  "username",
				Value: "root",
				Usage: "Memphis username",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "memphis",
				Usage: "Memphis password",
			},
			&cli.BoolFlag{
				Name:  "no-pwd",
				Value: false,
				Usage: "Ignore password",
			},
			&cli.StringFlag{
				Name:  "token",
				Value: "",
				Usage: "Memphis connection token",
			},
			&cli.StringFlag{
				Name:  "station",
				Value: "messages",
				Usage: "Memphis station",
			},
			&cli.StringFlag{
				Name:  "producer",
				Value: "messages",
				Usage: "Memphis producer name",
			},
			&cli.StringFlag{
				Name:  "message",
				Value: tables.Email,
				Usage: "Message type to send",
			},
			&cli.StringFlag{
				Name:  "recipient",
				Value: "",
				Usage: "Message recipient",
			},
			&cli.StringFlag{
				Name:  "subject",
				Value: "Testing message",
				Usage: "Msg subject",
			},
			&cli.StringFlag{
				Name:  "body",
				Value: "Please ignore this message, is for testing",
				Usage: "Msg body",
			},
		},
		Action: func(ctx *cli.Context) error {
			var opts []memphis.Option
			if !ctx.Bool("no-pwd") {
				opts = append(opts, memphis.Password(ctx.String("password")))
			}
			if token := ctx.String("token"); token != "" {
				opts = append(opts, memphis.ConnectionToken(token))
			}
			conn, err := memphis.Connect(
				ctx.String("host"),
				ctx.String("username"),
				opts...,
			)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			prod, pErr := conn.CreateProducer(ctx.String("station"), ctx.String("producer"))
			if pErr != nil {
				log.Fatal(pErr)
			}
			defer prod.Destroy()
			message := tables.RandomMessage(ctx.String("message"))
			message.Recipient = ctx.String("recipient")
			message.Subject = ctx.String("subject")
			message.Message = ctx.String("body")
			msgBytes, mErr := json.Marshal(message)
			if mErr != nil {
				log.Fatal(mErr)
			}
			prErr := prod.Produce(msgBytes)
			if pErr != nil {
				log.Fatal(prErr)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
