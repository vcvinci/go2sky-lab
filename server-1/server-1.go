package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	"go-skywalking/common"
	"github.com/SkyAPM/go2sky/propagation"
	"github.com/SkyAPM/go2sky/reporter"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/common"
	"github.com/gin-gonic/gin"

	gg "github.com/SkyAPM/go2sky/plugins/gin"
)

const (
	server_name        = "server-1"
	server_port        = 8081
	remote_server_name = "server-2"
	remote_server_addr = "localhost:8082"
	remoto_path        = "/user/info"
)

func main() {

	r := gin.New()
	rp, err := reporter.NewGRPCReporter(common.SkyAddr)
	common.PanicError(err)
	tracer, err := go2sky.NewTracer(server_name, go2sky.WithReporter(rp))
	common.PanicError(err)
	r.Use(gg.Middleware(r, tracer))

	r.GET("/trace", func(context *gin.Context) {
		span, ctx, err := tracer.CreateLocalSpan(context.Request.Context())
		common.PanicError(err)
		span.SetOperationName("Trace")

		context.Request = context.Request.WithContext(ctx)
		span.Log(time.Now(), "[Trace]", fmt.Sprintf(server_name+" satrt, params : %s", time.Now().Format("15:04:05")))

		result := make([]map[string]interface{}, 0)

		{
			url := fmt.Sprintf("http://%s%s", remote_server_addr, remoto_path)
			params := common.Params{
				Name: server_name + time.Now().Format("15:04:05"),
			}
			buffer := &bytes.Buffer{}
			_ = json.NewEncoder(buffer).Encode(params)
			req, err := http.NewRequest(http.MethodPost, url, buffer)
			common.PanicError(err)

			// op_name 是每一个操作的名称
			reqSpan, err := tracer.CreateExitSpan(context.Request.Context(), "invoke - "+remote_server_name, fmt.Sprintf("localhost:8082/user/info"), func(header string) error {
				req.Header.Set(propagation.Header, header)
				return nil
			})
			common.PanicError(err)
			reqSpan.SetComponent(2)                         //HttpClient,看 https://github.com/apache/skywalking/blob/master/docs/en/guides/Component-library-settings.md ， 目录在component-libraries.yml文件配置
			reqSpan.SetSpanLayer(v3.SpanLayer_RPCFramework) // rpc 调用

			resp, err := http.DefaultClient.Do(req)
			common.PanicError(err)
			defer resp.Body.Close()

			reqSpan.Log(time.Now(), "[HttpRequest]", fmt.Sprintf("开始请求,请求服务:%s,请求地址:%s,请求参数:%+v", remote_server_name, url, params))
			body, err := ioutil.ReadAll(resp.Body)
			common.PanicError(err)
			fmt.Printf("接受到消息： %s\n", body)
			reqSpan.Tag(go2sky.TagHTTPMethod, http.MethodPost)
			reqSpan.Tag(go2sky.TagURL, url)
			reqSpan.Log(time.Now(), "[HttpRequest]", fmt.Sprintf("结束请求,响应结果: %s", body))
			reqSpan.End()
			res := map[string]interface{}{}
			err = json.Unmarshal(body, &res)
			common.PanicError(err)
			result = append(result, res)
		}

		{
			url := fmt.Sprintf("http://%s%s", remote_server_addr, remoto_path)

			params := common.Params{
				Name: server_name + time.Now().Format("15:04:05"),
			}
			buffer := &bytes.Buffer{}
			_ = json.NewEncoder(buffer).Encode(params)
			req, err := http.NewRequest(http.MethodPost, url, buffer)
			common.PanicError(err)

			// op_name 是每一个操作的名称
			reqSpan, err := tracer.CreateExitSpan(context.Request.Context(), "invoke - "+remote_server_name, fmt.Sprintf("localhost:8082/user/info"), func(header string) error {
				req.Header.Set(propagation.Header, header)
				return nil
			})
			common.PanicError(err)
			reqSpan.SetComponent(2)                         //HttpClient,看 https://github.com/apache/skywalking/blob/master/docs/en/guides/Component-library-settings.md ， 目录在component-libraries.yml文件配置
			reqSpan.SetSpanLayer(v3.SpanLayer_RPCFramework) // rpc 调用

			resp, err := http.DefaultClient.Do(req)
			common.PanicError(err)
			defer resp.Body.Close()

			reqSpan.Log(time.Now(), "[HttpRequest]", fmt.Sprintf("开始请求,请求服务:%s,请求地址:%s,请求参数:%+v", remote_server_name, url, params))
			body, err := ioutil.ReadAll(resp.Body)
			common.PanicError(err)
			fmt.Printf("接受到消息： %s\n", body)

			reqSpan.Tag(go2sky.TagHTTPMethod, http.MethodPost)
			reqSpan.Tag(go2sky.TagURL, url)
			reqSpan.Log(time.Now(), "[HttpRequest]", fmt.Sprintf("结束请求,响应结果: %s", body))
			reqSpan.End()
			res := map[string]interface{}{}
			err = json.Unmarshal(body, &res)
			common.PanicError(err)
			result = append(result, res)
		}

		local := gin.H{
			"msg": result,
		}
		context.JSON(200, local)
		span.Log(time.Now(), "[Trace]", fmt.Sprintf(server_name+" end, resp : %s", local))
		span.End()
		{
			span, ctx, err := tracer.CreateEntrySpan(context.Request.Context(), "Send", func() (s string, e error) {
				return "", nil
			})
			context.Request = context.Request.WithContext(ctx)
			common.PanicError(err)
			span.SetOperationName("Send")
			//span.Error(time.Now(), "[Error]", "time is too long")
			span.Log(time.Now(), "[Info]", "send resp")
			span.End()
		}
	})

	common.PanicError(http.ListenAndServe(fmt.Sprintf(":%d", server_port), r))

}
