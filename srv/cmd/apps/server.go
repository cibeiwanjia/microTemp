// Package app

package cmd

import "C"
import (
	"context"
	"fmt"
	"log"

	hi "github.com/cibeiwanjia/microTemp/proto/hipb"

	"net"
	"os"
	"strconv"

	"github.com/cibeiwanjia/microTemp/pkg/logger"
	"github.com/cibeiwanjia/microTemp/pkg/storage/cache/local"
	"github.com/cibeiwanjia/microTemp/pkg/storage/cache/redis"
	"github.com/cibeiwanjia/microTemp/pkg/storage/conf"
	"github.com/cibeiwanjia/microTemp/pkg/storage/db/mysql"
	"github.com/cibeiwanjia/microTemp/srv/internal/logic"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	componentServer = "./server"
)

var (
	configFilePath string
)

var cmdRun = &cobra.Command{
	Use:     "srv",
	Example: fmt.Sprintf("%s srv -c apps", componentServer),
	Short:   "srv",
	Long:    `a srv test`,
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

	fmt.Println("-------->conf1:", conf.Cfg.Mysql)
	gr := run.Group{}
	ctx, logCancel := context.WithCancel(context.Background())
	if err := logger.Init(ctx, conf.Cfg.Log); err != nil {
		fmt.Println("err init failed", err)
		os.Exit(1)
	}

	//if conf.Cfg.Chain.Enable == true {
	//	if err := jaeger.Init(conf.Cfg.Chain); err != nil {
	//		logger.L.Error("local init failed: " + err.Error())
	//		os.Exit(1)
	//	}
	//}

	if _, err := mysql.Init(conf.Cfg.Chain.Enable, conf.Cfg.Mode, conf.Cfg.Mysql); err != nil {
		logger.L.Error("mysql init failed: " + err.Error())
		os.Exit(1)
	}
	//chainEnable bool, addr string, db int, password string, poolSize int, maxIdle int
	if err := redis.Init(conf.Cfg.Chain.Enable, conf.Cfg.Redis); err != nil {
		logger.L.Error("redis init failed: " + err.Error())
		os.Exit(1)
	}

	if err := local.Init(); err != nil {
		logger.L.Error("local init failed: " + err.Error())
		os.Exit(1)
	}

	{
		address := conf.Cfg.SRV.Host
		port := conf.Cfg.SRV.Port
		//监听端口
		lis, err := net.Listen("tcp", net.JoinHostPort(address, strconv.Itoa(port)))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		chainEnable := conf.Cfg.Chain.Enable
		s := grpc.NewServer()
		//调用链的服务地址不为空，说明需要调用链的服务
		if chainEnable == true {
			s = grpc.NewServer(
				grpc.StatsHandler(otelgrpc.NewServerHandler()), // 设置 StatsHandler
			)
		}

		//微服务注册
		hi.RegisterGreeterServer(s, &logic.HiServer{})

		//grpc的健康检查服务
		// 确保这部分代码被执行
		healthSrv := health.NewServer()
		healthpb.RegisterHealthServer(s, healthSrv)
		healthSrv.SetServingStatus("health.check", healthpb.HealthCheckResponse_SERVING)

		// 连接到Consul并注册服务
		consulClient, err := NewConsul(conf.Cfg.Consul.Host + ":" + strconv.Itoa(conf.Cfg.Consul.Port)) // Consul地址
		if err != nil {
			log.Fatalf("failed to connect to consul: %v", err)
		}
		fmt.Println("------------>consulClient:", conf.Cfg.Consul)
		// consul 注册中心
		services := conf.Cfg.Consul.Services
		for _, srv := range services {
			err = consulClient.RegisterService(srv)
			if err != nil {
				log.Fatalf("failed to register service with consul: %v", err)
			}
		}
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}
	if err := gr.Run(); err != nil {
		logger.L.Error(err.Error())
	}

	logger.L.Info("exiting")
	logCancel()

}

type Consul struct {
	client *api.Client
}

func NewConsul(addr string) (*Consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Consul{c}, nil
}

func (c *Consul) RegisterService(serviceInfo conf.ServiceInfo) error {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", serviceInfo.Host, serviceInfo.Port),
		Timeout:                        "5s",  // 缩短超时时间
		Interval:                       "30s", // 延长检查间隔
		DeregisterCriticalServiceAfter: "2m",  // 延长注销时间
	}
	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceInfo.Name, serviceInfo.Host, serviceInfo.Port), // 服务唯一ID
		Name:    serviceInfo.Name,                                                              // 服务名称
		Tags:    serviceInfo.Tags,                                                              // 为服务打标签
		Address: serviceInfo.Host,
		Port:    serviceInfo.Port,
		Check:   check,
	}
	return c.client.Agent().ServiceRegister(srv)
}
