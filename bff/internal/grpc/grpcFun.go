package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

var (
	Span trace.Span
)

func TraceCtx(c *gin.Context) context.Context {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		//"client-id", "from bff id 1233213434324322",
		//"user-id", "211",
	)
	Tracer := otel.Tracer(c.Request.URL.Path)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ctx, Span = Tracer.Start(ctx, c.Request.URL.Path, trace.WithAttributes(attribute.String("params", c.Request.URL.RawQuery)))
	return ctx
}

func discoverService(serviceName string, host string, port int) (string, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", host, port)

	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}

	services, err := client.Agent().Services()
	if err != nil {
		return "", err
	}

	for _, service := range services {
		if service.Service == serviceName {
			return net.JoinHostPort(service.Address, strconv.Itoa(service.Port)), nil
		}
	}

	return "", nil
}
