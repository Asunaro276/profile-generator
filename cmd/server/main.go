package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryuhei/randomuser-go/internal/config"
	"github.com/ryuhei/randomuser-go/internal/generator"
	"github.com/ryuhei/randomuser-go/internal/infrastructure/controller"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗: %v", err)
	}

	// 設定値を環境変数として設定
	if err := config.SetEnv(cfg); err != nil {
		log.Printf("環境変数の設定に失敗: %v", err)
		// 続行 - 環境変数の設定が失敗してもアプリケーションは起動可能
	}

	// ジェネレーターの読み込み
	gen := &generator.Generator{}
	err = gen.LoadGenerators()
	if err != nil {
		log.Fatalf("ジェネレーターの読み込みに失敗: %v", err)
	}

	// ルーターの設定
	router := gin.Default()
	// CORSミドルウェア設定
	router.Use(corsMiddleware())

	// APIルート
	api := router.Group("/api")
	{
		api.GET("", func(c *gin.Context) {
			controller.GenerateUser(c, gen, cfg)
		})
	}
	// HTTPサーバーの設定
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// グレースフルシャットダウンの設定
	go func() {
		log.Printf("サーバーを起動しました: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバーエラー: %v", err)
		}
	}()

	// シグナル処理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// サーバーのシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("サーバーのシャットダウンに失敗: %v", err)
	}
	log.Println("サーバーをシャットダウンしました")
}

// CORSミドルウェア
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
