package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

// TestUserRegistrationWithContainers 測試用戶註冊
func (ts *ContainerTestSuite) TestUserRegistrationWithContainers() {
	body := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(ts.T(), err)
	require.Contains(ts.T(), response, "token")
}

// TestUserLoginWithContainers 測試用戶登入
func (ts *ContainerTestSuite) TestUserLoginWithContainers() {
	// 先註冊用戶
	registerBody := map[string]interface{}{
		"username": "loginuser",
		"email":    "login@example.com",
		"password": "password123",
	}
	registerJSON, _ := json.Marshal(registerBody)
	registerReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerJSON))
	registerReq.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	ts.router.ServeHTTP(registerW, registerReq)

	// 測試登入
	loginBody := map[string]interface{}{
		"email":    "login@example.com",
		"password": "password123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	loginReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	ts.router.ServeHTTP(loginW, loginReq)

	require.Equal(ts.T(), http.StatusOK, loginW.Code)

	var response map[string]interface{}
	err := json.Unmarshal(loginW.Body.Bytes(), &response)
	require.NoError(ts.T(), err)
	require.Contains(ts.T(), response, "token")
}

// TestUserEndToEndFlowWithContainers 測試完整的用戶流程
func (ts *ContainerTestSuite) TestUserEndToEndFlowWithContainers() {
	// 註冊並登入用戶
	token := ts.createTestUserAndLogin("e2e@example.com", "password123", "e2euser")
	require.NotEmpty(ts.T(), token)

	// 測試獲取用戶資料
	profileReq, _ := http.NewRequest("GET", "/api/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+token)
	profileW := httptest.NewRecorder()
	ts.router.ServeHTTP(profileW, profileReq)

	require.Equal(ts.T(), http.StatusOK, profileW.Code)
}

// TestUserWithoutTokenWithContainers 測試無 token 訪問
func (ts *ContainerTestSuite) TestUserWithoutTokenWithContainers() {
	// 測試無 token 訪問受保護的端點
	req, _ := http.NewRequest("GET", "/api/tasks", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusUnauthorized, w.Code)
}

// TestUserWithInvalidTokenWithContainers 測試無效 token
func (ts *ContainerTestSuite) TestUserWithInvalidTokenWithContainers() {
	// 測試無效 token
	req, _ := http.NewRequest("GET", "/api/tasks", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusUnauthorized, w.Code)
}
