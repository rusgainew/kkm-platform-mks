package middleware

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// JWTUnaryClientInterceptor добавляет JWT token в gRPC metadata для исходящих запросов
func JWTUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Извлекаем JWT из контекста (добавлено middleware JWT auth)
		token := extractTokenFromContext(ctx)

		// Если token есть, добавляем в metadata
		if token != "" {
			// Создаем новый контекст с метаданными
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		// Вызываем оригинальный invoker с обновленным контекстом
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// JWTStreamClientInterceptor добавляет JWT token в gRPC metadata для stream запросов
func JWTStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Извлекаем JWT из контекста
		token := extractTokenFromContext(ctx)

		// Если token есть, добавляем в metadata
		if token != "" {
			// Создаем новый контекст с метаданными
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		// Вызываем оригинальный streamer с обновленным контекстом
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// extractTokenFromContext извлекает JWT token из контекста
// Token может быть добавлен middleware JWT auth в значение "user_token"
func extractTokenFromContext(ctx context.Context) string {
	// Проверяем значение "user_token" в контексте (добавленное JWT middleware)
	if token, ok := ctx.Value("user_token").(string); ok && token != "" {
		return token
	}

	// Проверяем значение "authorization" в контексте
	if authHeader, ok := ctx.Value("authorization").(string); ok && authHeader != "" {
		// Извлекаем token из "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
		return authHeader
	}

	return ""
}
