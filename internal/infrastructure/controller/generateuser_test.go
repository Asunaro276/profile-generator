package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ryuhei/randomuser-go/internal/config"
	"github.com/ryuhei/randomuser-go/internal/model"
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
			mockReturnJSON: `{"results":[{"gender":"male","name":{"title":"","first":"Test","last":"User"},"location":{"street":{"number":0,"name":""},"city":"","state":"","country":"","postcode":"","coordinates":{"latitude":"","longitude":""}},"email":"","login":{"uuid":"","username":"","password":"","salt":"","md5":"","sha1":"","sha256":""},"dob":{"date":"","age":0},"registered":{"date":"","age":0},"phone":"","cell":"","id":{"name":"","value":""},"picture":{"large":"","medium":"","thumbnail":""},"nat":""},{"gender":"male","name":{"title":"","first":"Test2","last":"User2"},"location":{"street":{"number":0,"name":""},"city":"","state":"","country":"","postcode":"","coordinates":{"latitude":"","longitude":""}},"email":"","login":{"uuid":"","username":"","password":"","salt":"","md5":"","sha1":"","sha256":""},"dob":{"date":"","age":0},"registered":{"date":"","age":0},"phone":"","cell":"","id":{"name":"","value":""},"picture":{"large":"","medium":"","thumbnail":""},"nat":""}],"info":{"seed":"12345","results":2,"page": 2}}`,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"results":[{"gender":"male","name":{"title":"","first":"Test","last":"User"},"location":{"street":{"number":0,"name":""},"city":"","state":"","country":"","postcode":"","coordinates":{"latitude":"","longitude":""}},"email":"","login":{"uuid":"","username":"","password":"","salt":"","md5":"","sha1":"","sha256":""},"dob":{"date":"","age":0},"registered":{"date":"","age":0},"phone":"","cell":"","id":{"name":"","value":""},"picture":{"large":"","medium":"","thumbnail":""},"nat":""},{"gender":"male","name":{"title":"","first":"Test2","last":"User2"},"location":{"street":{"number":0,"name":""},"city":"","state":"","country":"","postcode":"","coordinates":{"latitude":"","longitude":""}},"email":"","login":{"uuid":"","username":"","password":"","salt":"","md5":"","sha1":"","sha256":""},"dob":{"date":"","age":0},"registered":{"date":"","age":0},"phone":"","cell":"","id":{"name":"","value":""},"picture":{"large":"","medium":"","thumbnail":""},"nat":""}],"info":{"seed":"12347","results":2,"page":2}}`,
			setUpMock: func(m *MockUserGenerator) {
				m.EXPECT().Generate(2, mock.AnythingOfType("int64"), 2, "male").Return(
					[]model.User{
						{
							Gender: "male",
							Name: model.Name{
								First: "Test",
								Last:  "User",
							},
						},
						{
							Gender: "male",
							Name: model.Name{
								First: "Test2",
								Last:  "User2",
							},
						},
					},
					nil,
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
				m.EXPECT().Generate(1, mock.AnythingOfType("int64"), 1, "").Return(
					[]model.User{
						{
							Gender: "male",
							Name: model.Name{
								First: "Test",
								Last:  "User",
							},
						},
					},
					assert.AnError,
				)
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
	clients = make(map[string]int)

	mockGen := NewMockUserGenerator(t)
	mockGen.EXPECT().Generate(3, mock.AnythingOfType("int64"), 1, "").Return(
		[]model.User{
			{
				Gender: "male",
				Name: model.Name{
					First: "Test1",
				},
			},
			{
				Gender: "male",
				Name: model.Name{
					First: "Test2",
				},
			},
			{
				Gender: "male",
				Name: model.Name{
					First: "Test3",
				},
			},
		},
		nil,
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

	req3, _ := http.NewRequest("GET", "/api/user?results=2", nil)
	r.ServeHTTP(w, req3)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "制限超過")
}
