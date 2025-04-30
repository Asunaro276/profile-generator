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
	cfg := &config.Config{
		MaxResults: 50,
		Limit:      100,
	}

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
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			mockGen := NewMockUserGenerator(t)
			tt.setUpMock(mockGen)

			req, _ := http.NewRequest("GET", "/api", nil)

			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			r.GET("/api", func(c *gin.Context) {
				GenerateUser(c, mockGen, cfg)
			})

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

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

func TestGenerateUserRateLimit(t *testing.T) {
	cfg := &config.Config{
		MaxResults: 50,
		Limit:      5,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	mockGen := NewMockUserGenerator(t)
	mockGen.EXPECT().Generate(3, mock.AnythingOfType("int64"), 1, "").Return(
		`{"results":[{"name":{"first":"Test1"}},{"name":{"first":"Test2"}},{"name":{"first":"Test3"}}]}`, nil,
	).Times(2)

	r.GET("/api/user", func(c *gin.Context) {
		GenerateUser(c, mockGen, cfg)
	})

	req1, _ := http.NewRequest("GET", "/api/user?results=3", nil)
	r.ServeHTTP(w, req1)
	assert.Equal(t, http.StatusOK, w.Code)
	w = httptest.NewRecorder()

	req2, _ := http.NewRequest("GET", "/api/user?results=3", nil)
	r.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusOK, w.Code)
	w = httptest.NewRecorder()

	req3, _ := http.NewRequest("GET", "/api/user?results=1", nil)
	r.ServeHTTP(w, req3)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "制限超過")
}
