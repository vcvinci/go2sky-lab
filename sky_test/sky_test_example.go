package go2sky_test

import (
	"context"
	"log"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

func ExampleNewTracer() {
	// Use gRPC reporter for production
	r, err := reporter.NewGRPCReporter("127.0.0.1:11800")
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()
	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	// This for test
	tracer.WaitUntilRegister()
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatalf("create new local span error %v \n", err)
	}
	span.SetOperationName("invoke data")
	span.Tag("kind", "outer")
	time.Sleep(time.Second)
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		log.Fatalf("create new sub local span error %v \n", err)
	}
	subSpan.SetOperationName("invoke inner")
	subSpan.Log(time.Now(), "inner", "this is right")
	time.Sleep(time.Second)
	subSpan.End()
	time.Sleep(500 * time.Millisecond)
	span.End()
	time.Sleep(time.Second)
	// Output:
}
