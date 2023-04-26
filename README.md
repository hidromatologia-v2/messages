# messages

[![Build](https://github.com/hidromatologia-v2/messages/actions/workflows/build.yaml/badge.svg)](https://github.com/hidromatologia-v2/messages/actions/workflows/build.yaml)
[![codecov](https://codecov.io/gh/hidromatologia-v2/messages/branch/main/graph/badge.svg?token=64EUME4QU2)](https://codecov.io/gh/hidromatologia-v2/messages)

## Docker

```shell
docker pull ghcr.io/hidromatologia-v2/messages:latest
```

### Compose example

```shell
docker compose up -d
```

## Binary

You can use the binary present in [Releases](https://github.com/hidromatologia-v2/messages/releases/latest). Or compile your own with.

```shell
go install github.com/hidromatologia-v2/messages@latest
```

## Config

List of environment variables used by the binary. Make sure to setup them as well in your deployments

| Variable             | Description                                                  | Example                                                      |
| -------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `MEMPHIS_STATION`    | Name for the station to **CREATE**/**CONNECT**               | `messages`                                                   |
| `MEMPHIS_CONSUMER`   | Name of the consumer.                                        | `consumer-1`                                                 |
| `MEMPHIS_HOST`       | Host or IP of the Memphis service                            | `10.10.10.10`                                                |
| `MEMPHIS_USERNAME`   | Memphis Username                                             | `root`                                                       |
| `MEMPHIS_PASSWORD`   | Memphis password, if ignored `MEMPHIS_CONN_TOKEN` will be used | `password`                                                   |
| `MEMPHIS_CONN_TOKEN` | Memphis connection token, if ignored `MEMPHIS_PASSWORD` will be used | `ABCD`                                                       |
| `POSTGRES_DSN`       | Postgres DSN to be used                                      | `host=127.0.0.1 user=sulcud password=sulcud dbname=sulcud port=5432 sslmode=disable` |
| `REDIS_ADDR`         | Host and port of the Redis service.                          | `10.10.10.10:9999`                                           |
| `REDIS_DB`           | Redis Database number.                                       | `1`                                                          |
| `SMTP_FROM`          | Email address to setup as **FROM** in the sent messages      | `sulcud@mail.com`                                            |
| `SMTP_HOST`          | Host of the SMTP server                                      | `10.10.10.10`                                                |
| `SMTP_PORT`          | Port of the SMTP server                                      | `25`                                                         |
| `SMTP_USERNAME`      | Username to setup authentication                             | `sulcud`                                                     |
| `SMTP_PASSWORD`      | Password to setup authentication                             | `password`                                                   |
| `SMTP_NO_TLS`        | Bool Environment variable to setup                           | `1`                                                          |

## Coverage

| [![coverage](https://codecov.io/gh/hidromatologia-v2/messages/branch/main/graphs/sunburst.svg?token=64EUME4QU2)](https://app.codecov.io/gh/hidromatologia-v2/messages) | [![coverage](https://codecov.io/gh/hidromatologia-v2/messages/branch/main/graphs/tree.svg?token=64EUME4QU2)](https://app.codecov.io/gh/hidromatologia-v2/messages) |
| ------------------------------------------------------------ | ------------------------------------------------------------ |

