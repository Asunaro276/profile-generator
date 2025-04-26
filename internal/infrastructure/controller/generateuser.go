package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryuhei/randomuser-go/internal/config"
	"github.com/ryuhei/randomuser-go/internal/generator"
)

var (
	clients   = make(map[string]int)
	clientsMu sync.Mutex
)

// クライアント制限のリセット
func init() {
	go func() {
		// ミリ秒をナノ秒に変換（デフォルト: 5分）
		resetInterval := 5 * 60 * 1000 * 1000000 // 5分

		for {
			time.Sleep(time.Duration(resetInterval))
			clientsMu.Lock()
			// 制限超過のクライアントを表示
			for client, count := range clients {
				if count >= 1000 { // 設定から読み込むべき
					fmt.Printf("制限超過: %s: %d\n", client, count)
				}
			}
			// リセット
			clients = make(map[string]int)
			clientsMu.Unlock()
		}
	}()
}

// generateUser はランダムユーザーを生成する
func GenerateUser(c *gin.Context, gen *generator.Generator, cfg *config.Config) {
	// IPアドレスの取得
	ip := c.ClientIP()

	// レート制限のチェック
	clientsMu.Lock()
	if clients[ip] >= cfg.Limit {
		clientsMu.Unlock()
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("制限超過: %d ユーザーのリクエストが1分間に行われました", clients[ip]),
		})
		return
	}
	clientsMu.Unlock()

	// シードの設定
	seed := time.Now().UnixNano()
	if seedParam := c.DefaultQuery("seed", ""); seedParam != "" {
		// 文字列からシードを生成
		seedInt, err := strconv.ParseInt(seedParam, 10, 64)
		if err == nil {
			seed = seedInt
		}
	}

	// ページ情報
	page := 1
	if pageparam := c.DefaultQuery("page", "1"); pageparam != "" {
		pagenum, err := strconv.Atoi(pageparam)
		if err == nil && pagenum > 0 {
			page = pagenum
		}
	}
	seed += int64(page)

	// リクエストパラメータの解析
	resultsStr := c.DefaultQuery("results", "1")
	results, err := strconv.Atoi(resultsStr)
	if err != nil || results < 1 || results > cfg.MaxResults {
		results = 1
	}

	// 性別の設定
	gender := c.DefaultQuery("gender", "")

	// ジェネレーターの実行
	output, err := gen.Generate(results, seed, page, gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レート制限のカウント更新
	clientsMu.Lock()
	if _, exists := clients[ip]; !exists {
		clients[ip] = results
	} else {
		clients[ip] += results
	}
	clientsMu.Unlock()

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(http.StatusOK, output)
}
