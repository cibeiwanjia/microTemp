package dao

import (
	"context"
	"fmt"
	db "github.com/cibeiwanjia/microTemp/pkg/storage/db/mysql"

	"github.com/cibeiwanjia/microTemp/srv/internal/model/po"
	"go.opentelemetry.io/otel"
)

func GetGood() string {
	db.DB.WithContext(TraceCtx()).Exec("SELECT * FROM goods")
	if err := db.DB.WithContext(TraceCtx()).Take(&po.Goods{}).Error; err != nil {
		//return err
	}
	return "good"
}

func TraceCtx() context.Context {
	var tracer = otel.Tracer("goods table")
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "goods table")
	defer span.End()
	fmt.Println("-------------------->Received md: ", ctx)
	return ctx
}
