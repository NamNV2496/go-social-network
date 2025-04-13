package interceptors

import (
	"context"
	"time"

	"github.com/namnv2496/post-service/internal/pkg/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			metric.MetricIncError("400", "get", "post")
			return nil, err
		}
		metric.MetricIncHits("400", "get", "post")
		metric.MetricObserveTime("400", "get", "post", float64(time.Now().Sub(start)))
		return resp, nil

	}
}

// func MetricsInterceptor() grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {
// 		start := time.Now()
// 		resp, err := handler(ctx, req)
// 		duration := time.Since(start).Seconds()

// 		method := "POST"
// 		path := info.FullMethod
// 		status := "OK"

// 		if err != nil {
// 			status = "ERROR"
// 			metric.MetricIncError(status, method, path)
// 			return resp, err
// 		}

// 		metric.MetricIncHits(status, method, path)
// 		metric.MetricObserveTime(status, method, path, duration)

// 		return resp, nil
// 	}
// }
