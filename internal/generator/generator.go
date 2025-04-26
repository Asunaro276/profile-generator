package generator

import (
	"context"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3クライアントのシングルトン実装
var (
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	clientOnce    sync.Once
)

// initS3Client は S3クライアントを初期化します
func initS3Client() {
	clientOnce.Do(func() {
		// AWS SDKの設定をロード
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Printf("AWS設定のロードに失敗: %v\n", err)
			return
		}

		// S3クライアントの作成
		s3Client = s3.NewFromConfig(cfg)
		// 署名付きURLジェネレーターの作成
		presignClient = s3.NewPresignClient(s3Client)
	})
}

// Generator はユーザー生成器
type Generator struct {
	Data map[string]interface{}
}

// UserResult はランダムユーザーの生成結果
type UserResult struct {
	Results []map[string]interface{} `json:"results" xml:"results"`
	Info    struct {
		Seed    string `json:"seed" xml:"seed"`
		Results int    `json:"results" xml:"results"`
		Page    int    `json:"page" xml:"page"`
	} `json:"info" xml:"info"`
}

// LoadGenerators はジェネレーターをロードする
func (g *Generator) LoadGenerators() error {
	// APIディレクトリの確認
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ワーキングディレクトリの取得に失敗: %v", err)
	}
	dataDir := filepath.Join(workDir, "internal", "data")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return fmt.Errorf("APIディレクトリが見つかりません: %v", err)
	}

	g.Data = make(map[string]interface{})
	// 名前リスト
	firstNamesMale, err := readLines(filepath.Join(dataDir, "male_first.txt"))
	if err == nil {
		g.Data["first_names_male"] = firstNamesMale
	}

	firstNamesFemale, err := readLines(filepath.Join(dataDir, "female_first.txt"))
	if err == nil {
		g.Data["first_names_female"] = firstNamesFemale
	}

	lastNames, err := readLines(filepath.Join(dataDir, "last.txt"))
	if err == nil {
		g.Data["last_names"] = lastNames
	}

	return nil
}

// readLines はファイルから行を読み込みリストで返す
func readLines(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

// Generate は指定された数のユーザーを生成
func (g *Generator) Generate(count int, seed int64, page int, gender string) (string, error) {
	// 乱数ジェネレーターの初期化 - これにより決定論的な結果が得られる
	rnd := mathrand.New(mathrand.NewSource(seed))

	// ユーザー生成
	users := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		users[i] = g.generateUser(gender, rnd)
	}

	// レスポンス形式
	result := UserResult{
		Results: users,
		Info: struct {
			Seed    string `json:"seed" xml:"seed"`
			Results int    `json:"results" xml:"results"`
			Page    int    `json:"page" xml:"page"`
		}{
			Seed:    fmt.Sprintf("%d", seed),
			Results: count,
			Page:    page,
		},
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// generateUser は1人のユーザーを生成
func (g *Generator) generateUser(gender string, rnd *mathrand.Rand) map[string]interface{} {
	// ユーザーデータの生成
	if gender == "" {
		if rnd.Intn(2) == 1 {
			gender = "male"
		} else {
			gender = "female"
		}
	}

	// 名前の取得
	var firstName, lastName string

	// 名前の生成
	if gender == "male" && g.Data["first_names_male"] != nil {
		firstNames := g.Data["first_names_male"].([]string)
		firstName = firstNames[rnd.Intn(len(firstNames))]
	} else if gender == "female" && g.Data["first_names_female"] != nil {
		firstNames := g.Data["first_names_female"].([]string)
		firstName = firstNames[rnd.Intn(len(firstNames))]
	} else {
		// デフォルト名 (国籍データがない場合)
		maleNames := []string{"John", "Robert", "Michael", "David", "William"}
		femaleNames := []string{"Mary", "Patricia", "Jennifer", "Linda", "Elizabeth"}

		if gender == "male" {
			firstName = maleNames[rnd.Intn(len(maleNames))]
		} else {
			firstName = femaleNames[rnd.Intn(len(femaleNames))]
		}
	}

	// 姓の生成
	if g.Data != nil && g.Data["last_names"] != nil {
		lastNames := g.Data["last_names"].([]string)
		lastName = lastNames[rnd.Intn(len(lastNames))]
	} else {
		// デフォルト姓 (国籍データがない場合)
		defaultLastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones"}
		lastName = defaultLastNames[rnd.Intn(len(defaultLastNames))]
	}

	// タイトル
	title := "Mr"
	if gender == "female" {
		title = "Ms"
	}

	// 写真番号
	var photoNumber int
	if gender == "male" {
		photoNumber = rnd.Intn(46) + 1
	} else {
		photoNumber = rnd.Intn(24) + 1
	}

	// メールアドレス
	email := strings.ToLower(firstName) + "." + strings.ToLower(lastName) + "@example.com"

	// ユーザーID生成 - 決定論的に生成
	userID := rnd.Int63n(100000000)

	// プロフィール画像のパス
	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		bucket = "profile-generator" // デフォルト値
	}
	thumbnailKey := fmt.Sprintf("%s/portrait (%d).png", gender, photoNumber)

	// 署名付きURL
	thumbnailURL, _ := generateSignedURL(bucket, thumbnailKey, 10*time.Minute)

	// エラー処理（本番環境では適切に処理すること）
	largeURL := fmt.Sprintf("https://example.com/placeholder/%s/large.png", gender)
	mediumURL := fmt.Sprintf("https://example.com/placeholder/%s/medium.png", gender)
	if thumbnailURL == "" {
		thumbnailURL = fmt.Sprintf("https://example.com/placeholder/%s/thumbnail.png", gender)
	}

	// ユーザーオブジェクトの作成
	return map[string]interface{}{
		"gender": gender,
		"name": map[string]interface{}{
			"title": title,
			"first": firstName,
			"last":  lastName,
		},
		"location": map[string]interface{}{
			"street": map[string]interface{}{
				"number": rnd.Intn(9999) + 1,
				"name":   generateRandomStreetWithRand(rnd),
			},
			"city":     generateRandomCityWithRand(rnd),
			"state":    generateRandomStateWithRand(rnd),
			"country":  "US",
			"postcode": fmt.Sprintf("%05d", rnd.Intn(99999)),
			"coordinates": map[string]interface{}{
				"latitude":  fmt.Sprintf("%.4f", -90.0+rnd.Float64()*180.0),
				"longitude": fmt.Sprintf("%.4f", -180.0+rnd.Float64()*360.0),
			},
		},
		"email": email,
		"login": map[string]interface{}{
			"uuid":     generateUUIDWithRand(rnd),
			"username": strings.ToLower(firstName + lastName + strconv.Itoa(rnd.Intn(99))),
			"password": generateRandomPasswordWithRand(rnd),
			"salt":     generateRandomStringWithRand(rnd, 16),
			"md5":      generateRandomStringWithRand(rnd, 32),
			"sha1":     generateRandomStringWithRand(rnd, 40),
			"sha256":   generateRandomStringWithRand(rnd, 64),
		},
		"dob": map[string]interface{}{
			"date": time.Now().AddDate(-rnd.Intn(80)-18, -rnd.Intn(12), -rnd.Intn(28)).Format(time.RFC3339),
			"age":  rnd.Intn(80) + 18,
		},
		"registered": map[string]interface{}{
			"date": time.Now().AddDate(-rnd.Intn(20), -rnd.Intn(12), -rnd.Intn(28)).Format(time.RFC3339),
			"age":  rnd.Intn(20),
		},
		"phone": fmt.Sprintf("(%03d)-%03d-%04d", rnd.Intn(1000), rnd.Intn(1000), rnd.Intn(10000)),
		"cell":  fmt.Sprintf("(%03d)-%03d-%04d", rnd.Intn(1000), rnd.Intn(1000), rnd.Intn(10000)),
		"id": map[string]interface{}{
			"name":  "ID",
			"value": fmt.Sprintf("%08d", userID),
		},
		"picture": map[string]interface{}{
			"large":     largeURL,
			"medium":    mediumURL,
			"thumbnail": thumbnailURL,
		},
		"nat": "US",
	}
}

// 決定論的なヘルパー関数
func generateRandomStreetWithRand(rnd *mathrand.Rand) string {
	streets := []string{"Main Street", "Park Avenue", "Oak Street", "Maple Avenue", "Cedar Road"}
	return streets[rnd.Intn(len(streets))]
}

func generateRandomCityWithRand(rnd *mathrand.Rand) string {
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"}
	return cities[rnd.Intn(len(cities))]
}

func generateRandomStateWithRand(rnd *mathrand.Rand) string {
	states := []string{"California", "New York", "Texas", "Florida", "Illinois"}
	return states[rnd.Intn(len(states))]
}

// シード値に依存したUUIDを生成
func generateUUIDWithRand(rnd *mathrand.Rand) string {
	uuid := make([]byte, 16)
	for i := range uuid {
		uuid[i] = byte(rnd.Intn(256))
	}
	// RFC 4122 バージョン4 UUID形式に準拠
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func generateRandomPasswordWithRand(rnd *mathrand.Rand) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 12)
	for i := range result {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

func generateRandomStringWithRand(rnd *mathrand.Rand, length int) string {
	const chars = "abcdef0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

// 署名付きURLを生成する関数
func generateSignedURL(bucket, key string, duration time.Duration) (string, error) {
	// クライアントの初期化（一度だけ実行される）
	initS3Client()

	// クライアントが初期化されていない場合はエラー
	if presignClient == nil {
		return "", fmt.Errorf("S3クライアントの初期化に失敗")
	}

	// 署名付きURLのリクエスト作成
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", fmt.Errorf("署名付きURLの生成に失敗: %v", err)
	}

	return request.URL, nil
}
