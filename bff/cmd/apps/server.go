// Package app

package cmd

import (
	"context"
	"fmt"
	"github.com/cibeiwanjia/microTemp/bff/internal/grpc"
	"github.com/cibeiwanjia/microTemp/bff/internal/router"
	"github.com/cibeiwanjia/microTemp/pkg/logger"
	"github.com/cibeiwanjia/microTemp/pkg/storage/chain/jaeger"
	"github.com/cibeiwanjia/microTemp/pkg/storage/conf"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
)

const (
	componentServer = "./server"
)

var (
	configFilePath string
)

var cmdRun = &cobra.Command{
	Use:     "bff",
	Example: fmt.Sprintf("%s bff -c apps", componentServer),
	Short:   "bff",
	Long:    `a bff test`,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		Run()
		os.Exit(0)
	},
}

func NewServerCommand() *cobra.Command {
	command := &cobra.Command{
		Use: componentServer,
	}
	cmdRun.PersistentFlags().StringVarP(&configFilePath, "config", "c", "./apps", "config path")
	command.AddCommand(cmdRun)
	return command
}

func Run() {
	conf.Init(configFilePath)
	gr := run.Group{}
	ctx, logCancel := context.WithCancel(context.Background())

	if err := logger.Init(ctx, conf.Cfg.Log); err != nil {
		fmt.Println("err init failed", err)
		os.Exit(1)
	}
	if conf.Cfg.Chain.Enable == true {
		//logger.OtlpInit()
		if err := jaeger.Init(conf.Cfg.Chain); err != nil {
			logger.L.Error("local init failed: " + err.Error())
			os.Exit(1)
		}
	}
	fmt.Println("consul config:", conf.Cfg.Consul)
	//微服发现
	if err := grpc.HiInit(conf.Cfg.Consul); err != nil {
		logger.L.Error("grpc init failed: " + err.Error())
		os.Exit(1)
	}

	{
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		gr.Add(
			func() error {
				<-term
				logger.L.Warn("Received SIGTERM, exiting gracefully...")
				return nil
			},
			func(err error) {},
		)
	}
	{
		cancel := make(chan struct{})
		gr.Add(func() error {
			gin.SetMode(conf.Cfg.Mode)
			srv := router.Server(conf.Cfg)
			router.GracefulExit(srv, cancel)
			return nil
		}, func(err error) {
			close(cancel)
		})
	}

	if err := gr.Run(); err != nil {
		logger.L.Error(err.Error())
	}

	logger.L.Info("exiting")
	logCancel()
}
