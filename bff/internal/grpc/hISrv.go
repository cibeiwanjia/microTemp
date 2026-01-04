package grpc

import (
	"fmt"
	"log"

	pb "github.com/cibeiwanjia/microTemp/proto/hipb"

	"github.com/cibeiwanjia/microTemp/pkg/storage/conf"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	HiClient pb.GreeterClient
	Tracer   trace.Tracer
)

func HiInit(cfg *conf.ConsulConfig) (err error) {
	fmt.Println("--------------->cfg.Services:", cfg.Services[0].Name)
	Tracer = otel.Tracer(cfg.Services[0].Name)

	//从consul 获取grpc服务地址
	grpcAddr, err := discoverService(cfg.Services[0].Name, cfg.Host, cfg.Port)
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // 设置 StatsHandler
		//grpc.WithChainUnaryInterceptor(clientUnaryInterceptor),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	HiClient = pb.NewGreeterClient(conn)
	return nil
}
