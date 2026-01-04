package hi

import (
	"fmt"
	"log"

	"net/http"
	"sync"
	"time"

	pb "github.com/cibeiwanjia/microTemp/proto/hipb"

	"github.com/cibeiwanjia/microTemp/bff/internal/grpc"
	"github.com/cibeiwanjia/microTemp/pkg/logger"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type HelloResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Hello  `json:"data"`
}

type Hello struct {
	Time time.Time
}

// Hello godoc
// @Tags Hello
// @Summary get Hello information
// @Accept text/plain
// @Produce json
// @Success 200 {object} HelloResponse
// @Failure default {object} HTTPError
// @Router /gethello [get]
func (c *HiController) GetHello(ctx *gin.Context) {
	//res, err := hiSrv.HiClient.SayHello(ctx, &pb.HelloRequest{Name: "world"})
	//md := metadata.Pairs(
	//	"timestamp", time.Now().Format(time.StampNano),
	//	//"client-id", "from bff id 1233213434324322",
	//	//"user-id", "211",
	//)
	chainCtx := grpc.TraceCtx(ctx)
	//租房
	res, err := grpc.HiClient.SayHello(chainCtx, &pb.HelloRequest{Name: "world"}) //200ms
	//二手房
	//res, err := grpc.HiClient.SayHi(chainCtx, &pb.HelloRequest{Name: "world"})   //200ms
	//xxx

	logger.L.Info("get hello", zap.Any("res", res))
	//添加tag
	otelLogger := otelzap.New(logger.L)
	otelLogger.Ctx(chainCtx).Error(">>>>>>>get hello", zap.Any("res", res))
	//添加日志
	grpc.Span.AddEvent("2222222	 -------------->")
	defer grpc.Span.End()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("----------------->Greeting: %s", res.Message)
	ctx.JSON(http.StatusOK, HelloResponse{Code: http.StatusOK, Message: "hello world", Data: Hello{Time: time.Now()}})
	//ctx.JSON(200, gin.H{"code": 0, "message": "hello world", "data": model.Hello{Time: time.Now()}})
}

func (c *HiController) GetHello2(ctx *gin.Context) {
	// 创建缓冲大小为2的通道，用于控制并发数
	concurrencyLimit := make(chan struct{}, 2)

	// 创建WaitGroup，用于等待所有goroutine完成
	var wg sync.WaitGroup

	// 获取追踪上下文
	chainCtx := grpc.TraceCtx(ctx)

	// 定义响应变量
	var helloRes *pb.HelloReply
	var hiRes *pb.HelloReply
	var helloErr, hiErr error

	// 调用SayHello服务
	concurrencyLimit <- struct{}{}
	wg.Add(1)
	go func() {
		defer func() {
			<-concurrencyLimit
			wg.Done()
		}()
		res, err := grpc.HiClient.SayHello(chainCtx, &pb.HelloRequest{Name: "world"}) //200ms
		if err == nil {
			helloRes = res
		} else {
			helloErr = err
		}
	}()

	// 调用sayHi服务（假设存在）
	concurrencyLimit <- struct{}{}
	wg.Add(1)
	go func() {
		defer func() {
			<-concurrencyLimit
			wg.Done()
		}()
		// 假设存在sayHi方法，接口类似SayHello
		res, err := grpc.HiClient.SayHello(chainCtx, &pb.HelloRequest{Name: "hi"}) //300ms
		if err == nil {
			hiRes = res
		} else {
			hiErr = err
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()
	close(concurrencyLimit)

	// 处理结果
	if helloErr != nil || hiErr != nil {
		logger.L.Error("gRPC调用失败", zap.Error(helloErr), zap.Error(hiErr))
		ctx.JSON(http.StatusInternalServerError, HelloResponse{
			Code:    http.StatusInternalServerError,
			Message: "服务调用失败",
		})
		return
	}

	// 记录日志
	logger.L.Info("get hello", zap.Any("helloRes", helloRes), zap.Any("hiRes", hiRes))

	// 合并两个服务的响应
	combinedMessage := fmt.Sprintf("%s + %s", helloRes.Message, hiRes.Message)

	// 返回结果
	ctx.JSON(http.StatusOK, HelloResponse{
		Code:    http.StatusOK,
		Message: combinedMessage,
		Data:    Hello{Time: time.Now()},
	})
}

func (c *HiController) GetHello3(ctx *gin.Context) {

	grpc.HiClient.SayHello(ctx, &pb.HelloRequest{Name: "world"}) //200ms
	//grpc.HiClient.SayHello2(ctx, &pb.HelloRequest{Name: "world"}) //200ms
	//grpc.HiClient.SayHello3(ctx, &pb.HelloRequest{Name: "world"}) //200ms
	//grpc.HiClient.SayHello4(ctx, &pb.HelloRequest{Name: "world"}) //200ms
	//grpc.HiClient.SayHello5(ctx, &pb.HelloRequest{Name: "world"}) //200ms

}
