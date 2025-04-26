package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config はアプリケーション設定を保持する構造体
type Config struct {
	Port          int    `json:"port"`
	Limit         int    `json:"limit"`
	MaxResults    int    `json:"maxResults"`
	ResetInterval int    `json:"resetInterval"`
	LatestVersion string `json:"latestVersion"`
}

// Load は設定ファイルから設定を読み込む
func Load() (*Config, error) {
	// 設定ファイルのパスを取得
	configPath := filepath.Join(".", "config.json")

	// 設定ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Port:          8080,
			Limit:         100,
			MaxResults:    5000,
			ResetInterval: 300000, // 5分 (ms)
			LatestVersion: "1.4",
		}, nil
	}

	// 設定ファイルを開く
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// JSONデコード
	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
