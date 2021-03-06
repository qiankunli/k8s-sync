package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s-sync/pkg/config"
	"k8s-sync/pkg/kube"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"io"
)

func main() {
	app := &cli.App{
		Name:  "k8s-sync",
		Usage: "sync apiServer to mysql",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			run(c)
			return nil
		},
	}
	app.Run(os.Args)
}

func run(cli *cli.Context) {
	configPath := cli.String("config")
	if len(configPath) == 0 {
		fmt.Println("please input config file path use -c flag")
		return
	}
	config, err := config.ReadYaml(configPath)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	setLog(config)

	kubeClient, err := kube.GetKubernetesClient(config)
	if err != nil {
		log.Err(err)
		return
	}
	w := kube.NewPodWatcher(kubeClient,config)
	w.Start()
	defer w.Stop()
	c := kube.NewCleaner(kubeClient, config)
	c.Start()
	defer c.Stop()
	stopper := make(chan int)
	<-stopper
}

func setLog(config *config.Config) {
	// 设置全局logger
	var writer io.Writer
	file, err := os.OpenFile(config.File, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("open log file %s error\n,print log to console", config.File)
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		writer = zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
		}
	}
	log.Logger = log.With().Caller().Logger().Output(writer)
	if config.Debug {
		log.Logger.Level(zerolog.DebugLevel)
	} else {
		log.Logger.Level(zerolog.InfoLevel)
	}
}
