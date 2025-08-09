package observability

import (
	"context"
	"log"
)

// InitTracer 初始化 OpenTelemetry TracerProvider
// 暂时简化，只是一个占位符
func InitTracer(ctx context.Context, serviceName, environment, version string) (func(context.Context) error, error) {
	log.Printf("Tracer initialized for service: %s, env: %s, version: %s", serviceName, environment, version)

	// 返回一个空的 cleanup 函数
	cleanup := func(ctx context.Context) error {
		log.Println("Tracer cleanup called")
		return nil
	}

	return cleanup, nil
}
