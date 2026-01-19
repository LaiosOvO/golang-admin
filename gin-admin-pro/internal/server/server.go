package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-admin-pro/internal/pkg/config"
	"gin-admin-pro/internal/router"
)

// Server 服务器结构
type Server struct {
	httpServer *http.Server
}

// NewServer 创建新的服务器实例
func NewServer() *Server {
	return &Server{}
}

// Start 启动服务器
func (s *Server) Start() error {
	cfg := config.GetConfig()

	// 初始化路由
	r := router.InitRouter()

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	fmt.Printf("服务器启动在端口 %d\n", cfg.Server.Port)
	fmt.Printf("健康检查: http://localhost:%d/health\n", cfg.Server.Port)
	fmt.Printf("API文档: http://localhost:%d/api/v1\n", cfg.Server.Port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %v", err)
	}

	return nil
}

// Stop 优雅关闭服务器
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("正在关闭服务器...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %v", err)
	}

	fmt.Println("服务器已关闭")
	return nil
}

// WaitForShutdown 等待关闭信号
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := s.Stop(); err != nil {
		fmt.Printf("关闭服务器时出错: %v\n", err)
	}
}
