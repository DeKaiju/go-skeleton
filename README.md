# go-skeleton

## Run

```bash
// env
cp .env.example .env

// build
go build main.go

// cpu core number
go build -ldflags "-X main.SetCpuCount=1" main.go

// run server
./main server

// run script
./main cmd demo hello
./main cmd demo world
```

## Test

```bash
go get github.com/smartystreets/goconvey

cd test/redis
go test -v

cd test/curl
goconvey
```

## Lib

| Role     | Package     |   Link   |
|----------|-------------| ---- |
| Script   | cobra       |   https://github.com/spf13/cobra     |
| Router   | gin         |   https://github.com/gin-gonic/gin     |
| Env      | godotenv    |   https://github.com/joho/godotenv     |
| ORM      | gorm        |   https://github.com/go-gorm/gorm      |
| Redis    | redigo      |   https://github.com/gomodule/redigo   |
| Curl     | goz         |   https://github.com/idoubi/goz        |
| Json     | gjson       |   https://github.com/tidwall/gjson     |
| Log      | logrus      |   https://github.com/sirupsen/logrus   |
| JWT      | jwt-go      |   https://github.com/golang-jwt/jwt  |
| Doc      | gin-swagger |   https://github.com/swaggo/gin-swagger |
| Test     | goconvey    |   https://github.com/smartystreets/goconvey |
| Kafka    | kafka-go    |   https://github.com/segmentio/kafka-go  |
| RabbitMq | amqp        |   https://github.com/streadway/amqp  |
| Etcd     | etcd        |   https://github.com/coreos/etcd/clientv3 |
