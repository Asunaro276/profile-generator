package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryuhei/randomuser-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateUser(t *testing.T) {
	// テスト用の設定
	cfg := &config.Config{
		MaxResults: 50,
		Limit:      100,
	}

	// テストケース
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockReturnJSON string
		mockError      error
		expectedStatus int
		expectedBody   string
		setUpMock      func(*MockUserGenerator)
	}{
		{
			name:           "正常なリクエスト",
			queryParams:    map[string]string{"results": "2", "gender": "male", "seed": "12345", "page": "2"},
			mockReturnJSON: `{"results":[{"name":{"first":"Test","last":"User"}, "gender": "male"}, {"name":{"first":"Test2","last":"User2"}, "gender": "male"}}]}`,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"results":[{"name":{"first":"Test","last":"User"}, "gender": "male"}, {"name":{"first":"Test2","last":"User2"}, "gender": "male"}}]}`,
			setUpMock: func(m *MockUserGenerator) {
				m.EXPECT().Generate(2, mock.AnythingOfType("int64"), 2, "male").Return(
					`{"results":[{"name":{"first":"Test","last":"User"}, "gender": "male"}, {"name":{"first":"Test2","last":"User2"}, "gender": "male"}}]}`, nil,
				)
			},
		},
		{
			name:           "ジェネレーターエラー",
			queryParams:    map[string]string{"results": "1"},
			mockReturnJSON: "",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"assert.AnError general error for testing"}`,
			setUpMock: func(m *MockUserGenerator) {
				m.EXPECT().Generate(1, mock.AnythingOfType("int64"), 1, "").Return("", assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ginのテスト用コンテキストを作成
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			// mockUserGeneratorの作成
			mockGen := NewMockUserGenerator(t)
			tt.setUpMock(mockGen)

			// テストリクエストの作成
			req, _ := http.NewRequest("GET", "/api", nil)

			// クエリパラメータの設定
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			// リクエストをコンテキストに設定
			r.GET("/api", func(c *gin.Context) {
				GenerateUser(c, mockGen, cfg)
			})

			// ハンドラを実行
			r.ServeHTTP(w, req)

			// レスポンスのアサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			// JSONレスポンスの検証（エラーレスポンスの場合）
			if tt.expectedStatus != http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], assert.AnError.Error())
			} else {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

// レート制限のテスト
func TestGenerateUserRateLimit(t *testing.T) {
	// テスト用の設定（低めの制限値）
	cfg := &config.Config{
		MaxResults: 50,
		Limit:      5, // 低い制限値を設定
	}

	// Ginのテスト用コンテキストを作成
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	// mockUserGeneratorの作成
	mockGen := NewMockUserGenerator(t)
	mockGen.EXPECT().Generate(3, mock.AnythingOfType("int64"), 1, "").Return(
		`{"results":[{"name":{"first":"Test1"}},{"name":{"first":"Test2"}},{"name":{"first":"Test3"}}]}`, nil,
	).Times(2) // 2回呼ばれることを期待

	// エンドポイントの設定
	r.GET("/api/user", func(c *gin.Context) {
		GenerateUser(c, mockGen, cfg)
	})

	// 最初のリクエスト - 3ユーザー
	req1, _ := http.NewRequest("GET", "/api/user?results=3", nil)
	r.ServeHTTP(w, req1)
	assert.Equal(t, http.StatusOK, w.Code)
	w = httptest.NewRecorder() // レコーダーをリセット

	// 2回目のリクエスト - 3ユーザー
	req2, _ := http.NewRequest("GET", "/api/user?results=3", nil)
	r.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	w = httptest.NewRecorder() // レコーダーをリセット

	// 3回目のリクエスト - 制限超過
	req3, _ := http.NewRequest("GET", "/api/user?results=1", nil)
	r.ServeHTTP(w, req3)

	// 429 Too Many Requestsが返されることを確認
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	// エラーメッセージの確認
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "制限超過")
}
