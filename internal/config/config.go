package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

// Config はアプリケーション設定を保持する構造体
type Config struct {
	Port          int    `json:"port"`
	Limit         int    `json:"limit"`
	MaxResults    int    `json:"maxResults"`
	ResetInterval int    `json:"resetInterval"`
	BucketName    string `json:"bucketName"`
}

// Load は設定ファイルから設定を読み込む
func Load() (*Config, error) {
	// 設定ファイルのパスを取得
	configPath := filepath.Join(".", "config.json")

	// 設定ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Port:          8080,
			MaxResults:    5000,
			BucketName:    "profile-generator",
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

// SetEnv は設定値を環境変数として設定します
func SetEnv(config *Config) error {
	if config == nil {
		return nil
	}

	// 数値を文字列に変換して環境変数に設定
	if err := os.Setenv("PORT", strconv.Itoa(config.Port)); err != nil {
		return err
	}
	if err := os.Setenv("MAX_RESULTS", strconv.Itoa(config.MaxResults)); err != nil {
		return err
	}
	if err := os.Setenv("BUCKET_NAME", config.BucketName); err != nil {
		return err
	}

	return nil
}
