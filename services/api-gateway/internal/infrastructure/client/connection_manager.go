package client

import (
	"context"
	"fmt"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnectionManager управляет gRPC соединениями
type ConnectionManager struct {
	connections map[string]*grpc.ClientConn
	dialTimeout time.Duration
	logger      *zap.Logger
	tlsConfig   *GRPCTLSConfig
}

// NewConnectionManager создает новый ConnectionManager
func NewConnectionManager(dialTimeout time.Duration, logger *zap.Logger) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*grpc.ClientConn),
		dialTimeout: dialTimeout,
		logger:      logger,
		tlsConfig:   &GRPCTLSConfig{Enabled: false},
	}
}

// NewConnectionManagerWithTLS создает ConnectionManager с TLS поддержкой
func NewConnectionManagerWithTLS(dialTimeout time.Duration, logger *zap.Logger, tlsConfig *GRPCTLSConfig) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*grpc.ClientConn),
		dialTimeout: dialTimeout,
		logger:      logger,
		tlsConfig:   tlsConfig,
	}
}

// GetConnection возвращает gRPC соединение для указанного адреса
func (cm *ConnectionManager) GetConnection(ctx context.Context, address string) (*grpc.ClientConn, error) {
	// Проверяем существующее соединение
	if conn, exists := cm.connections[address]; exists {
		// Проверяем состояние соединения
		state := conn.GetState()
		if state.String() != "SHUTDOWN" && state.String() != "TRANSIENT_FAILURE" {
			return conn, nil
		}

		// Закрываем неработающее соединение
		conn.Close()
		delete(cm.connections, address)
	}

	// Создаем новое соединение
	conn, err := cm.dial(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	cm.connections[address] = conn
	cm.logger.Info("gRPC connection established", zap.String("address", address))

	return conn, nil
}

// dial создает новое gRPC соединение
func (cm *ConnectionManager) dial(ctx context.Context, address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, cm.dialTimeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024), // 50MB
			grpc.MaxCallSendMsgSize(50*1024*1024), // 50MB
		),
		// Добавляем JWT interceptors для пропагирования токена в metadata
		grpc.WithUnaryInterceptor(middleware.JWTUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(middleware.JWTStreamClientInterceptor()),
	}

	// Добавляем TLS credentials если включен
	if cm.tlsConfig != nil && cm.tlsConfig.Enabled {
		cm.logger.Warn("TLS is configured but implementation pending - using insecure connection", zap.String("address", address))
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		cm.logger.Warn("Using insecure gRPC connection - not recommended for production", zap.String("address", address))
	}

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CloseAll закрывает все gRPC соединения
func (cm *ConnectionManager) CloseAll() {
	for address, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			cm.logger.Error("Failed to close gRPC connection",
				zap.String("address", address),
				zap.Error(err),
			)
		} else {
			cm.logger.Info("gRPC connection closed", zap.String("address", address))
		}
	}
	cm.connections = make(map[string]*grpc.ClientConn)
}

// Close закрывает конкретное соединение
func (cm *ConnectionManager) Close(address string) error {
	if conn, exists := cm.connections[address]; exists {
		delete(cm.connections, address)
		return conn.Close()
	}
	return nil
}
