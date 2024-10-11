package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitConfig(path string) {

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

var (
	DB    *gorm.DB
	REDIS *redis.Client
)

func InitMysql() {
	mysqllogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      true,
	},
	)
	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: mysqllogger})
	DB = db
	if err != nil {
		panic("failed to connect database")
	}
}
func InitRedis() {
	REDIS = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdelConn"),
	})
}

const (
	PublishKey = "websocket"
)

func Publish(ctx context.Context, channel string, msg string) error {
	err := REDIS.Publish(ctx, channel, msg).Err()
	fmt.Println("Publish", msg)
	if err != nil {
		fmt.Println("err :", err)
	}
	return nil
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	str := REDIS.Subscribe(ctx, channel)

	msg, err := str.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("err :", err)
	}
	fmt.Println("subscrib", msg.Payload)
	return msg.Payload, nil
}
