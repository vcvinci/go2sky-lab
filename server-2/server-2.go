package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	"go-skywalking/common"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"

	gg "github.com/SkyAPM/go2sky/plugins/gin"
)

const (
	server_name_2 = "server-2"
	server_port_2 = 8082
)

func main() {

	r := gin.New()
	rp, err := reporter.NewGRPCReporter(common.SkyAddr)
	common.PanicError(err)
	tracer, err := go2sky.NewTracer(server_name_2, go2sky.WithReporter(rp))
	common.PanicError(err)
	r.Use(gg.Middleware(r, tracer))

	r.POST("/user/info", func(context *gin.Context) {
		span, ctx, err := tracer.CreateLocalSpan(context.Request.Context())
		common.PanicError(err)
		span.SetOperationName("UserInfo")
		context.Request = context.Request.WithContext(ctx)
		params := new(common.Params)
		err = context.BindJSON(params)
		common.PanicError(err)

		span.Log(time.Now(), "[UserInfo]", fmt.Sprintf(server_name_2+" satrt, req : %+v", params))
		local := gin.H{
			"msg": fmt.Sprintf(server_name_2+" time : %s", time.Now().Format("15:04:05")),
		}
		context.JSON(200, local)
		span.Log(time.Now(), "[UserInfo]", fmt.Sprintf(server_name_2+" end, resp : %s", local))
		span.End()
	})

	common.PanicError(http.ListenAndServe(fmt.Sprintf(":%d", server_port_2), r))

}
