package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/cibeiwanjia/microTemp/pkg/storage/cache/redis"
	db "github.com/cibeiwanjia/microTemp/pkg/storage/db/mysql"
	"github.com/cibeiwanjia/microTemp/proto/hipb"

	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type HiServer struct {
	hipb.UnimplementedGreeterServer
}

func (s *HiServer) SayHello(ctx context.Context, in *hipb.HelloRequest) (*hipb.HelloReply, error) {
	var tracer = otel.Tracer("goods")

	ctx, span := tracer.Start(ctx, "goods_gorm")
	defer span.End()

	//redis
	if err := redis.RC.Set(ctx, "name", "Q1mi", time.Minute).Err(); err != nil {
		//return err
	}
	if err := redis.RC.Set(ctx, "tag", "OTel", time.Minute).Err(); err != nil {
		//return err
	}
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := redis.RC.Get(ctx, "tag").Val()
			if val != "OTel" {
				log.Printf("%q != %q", val, "OTel")
			}
		}()
	}
	wg.Wait()

	if err := redis.RC.Del(ctx, "name").Err(); err != nil {
		//return err
	}
	if err := redis.RC.Del(ctx, "tag").Err(); err != nil {
		//return err
	}
	//gorm
	type Book struct {
		gorm.Model
		Title string
	}
	// 迁移 schema
	db.DB.AutoMigrate(&Book{})
	// Create
	db.DB.WithContext(ctx).Create(&Book{Title: "《Go语言之路》"})
	// Read
	var book Book
	if err := db.DB.WithContext(ctx).Take(&book).Error; err != nil {
		//return err
	}
	// delete
	if err := db.DB.WithContext(ctx).Delete(&book).Error; err != nil {
		//return err
	}
	//fmt.Println("-------------------->Received md: ", md.Get("client-id"))
	//log.Printf("--------------------->Received: %v", in.GetName())
	//DB.table("user").Create(&User{Name: in.GetName()})
	return &hipb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *HiServer) Hi() {
	time.Sleep(time.Second * 10)
	fmt.Println("Hi")
}
