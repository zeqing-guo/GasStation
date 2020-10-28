package main

import (
	"context"
	"os"
	"strconv"

	r "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zooFinance/bebop/store"
	"github.com/zooFinance/tradebot/gas"
)

func main() {

	app := cli.NewApp()
	app.Name = "gas"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "redis_url", Value: "localhost:6379", Usage: "redis url", EnvVar: "REDIS_URL"},
		cli.StringFlag{Name: "gasnow_url", Value: "https://www.gasnow.org/api/v3/gas/data", Usage: "GasNow URL", EnvVar: "GASNOW_URL"},
		cli.StringFlag{Name: "log_level", Value: "debug", Usage: "log level", EnvVar: "LOG_LEVEL"},
	}
	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) (err error) {
	lvl, err := log.ParseLevel(ctx.String("log_level"))
	if err != nil {
		log.WithError(err).Fatalln("unvalid log level")
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	g := gas.New([]string{ctx.String("gasnow_url")}, gas.GetGasNow)
	g.Run()
	store.InitRedis(&r.Options{
		Addr:     ctx.String("redis_url"),
		Password: "",
		DB:       0,
	})
	for {
		fastest, fast, average, err := g.GetGasPrice()
		if err != nil {
			log.WithError(err).Error("can't get gas price")
			continue
		}
		fastestStr := strconv.FormatFloat(fastest, 'E', -1, 64)
		fastStr := strconv.FormatFloat(fast, 'E', -1, 64)
		averageStr := strconv.FormatFloat(average, 'E', -1, 64)
		store.RedisClient.Set(context.Background(), "gas.fastest", fastestStr, -1)
		store.RedisClient.Set(context.Background(), "gas.fast", fastStr, -1)
		store.RedisClient.Set(context.Background(), "gas.average", averageStr, -1)
	}
}