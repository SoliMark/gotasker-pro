package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTaskCRUDWithContainers 測試任務 CRUD 操作
func (ts *ContainerTestSuite) TestTaskCRUDWithContainers() {
	// 創建用戶並獲取 token
	token := ts.createTestUserAndLogin("taskuser@example.com", "password123", "taskuser")
	require.NotEmpty(ts.T(), token)

	// 創建任務
	createBody := map[string]interface{}{
		"title":   "Test Task",
		"content": "Test Content",
		"status":  "pending",
	}
	createJSON, _ := json.Marshal(createBody)
	createReq, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	require.Equal(ts.T(), http.StatusCreated, createW.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResponse)
	taskID := createResponse["id"].(float64)

	// 獲取任務
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/tasks/%d", int(taskID)), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	require.Equal(ts.T(), http.StatusOK, getW.Code)

	// 更新任務
	updateBody := map[string]interface{}{
		"title":   "Updated Task",
		"content": "Updated Content",
		"status":  "done",
	}
	updateJSON, _ := json.Marshal(updateBody)
	updateReq, _ := http.NewRequest("PUT", fmt.Sprintf("/api/tasks/%d", int(taskID)), bytes.NewBuffer(updateJSON))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateW := httptest.NewRecorder()
	ts.router.ServeHTTP(updateW, updateReq)

	require.Equal(ts.T(), http.StatusOK, updateW.Code)

	// 獲取任務列表
	listReq, _ := http.NewRequest("GET", "/api/tasks", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listW := httptest.NewRecorder()
	ts.router.ServeHTTP(listW, listReq)

	require.Equal(ts.T(), http.StatusOK, listW.Code)

	// 刪除任務
	deleteReq, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/tasks/%d", int(taskID)), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteW := httptest.NewRecorder()
	ts.router.ServeHTTP(deleteW, deleteReq)

	require.Equal(ts.T(), http.StatusNoContent, deleteW.Code)
}

// TestTaskAuthorizationWithContainers 測試任務授權
func (ts *ContainerTestSuite) TestTaskAuthorizationWithContainers() {
	// 創建兩個用戶
	token1 := ts.createTestUserAndLogin("user1@example.com", "password123", "user1")
	token2 := ts.createTestUserAndLogin("user2@example.com", "password123", "user2")
	require.NotEmpty(ts.T(), token1)
	require.NotEmpty(ts.T(), token2)

	// 用戶1創建任務
	createBody := map[string]interface{}{
		"title":   "User1 Task",
		"content": "User1 Content",
		"status":  "pending",
	}
	createJSON, _ := json.Marshal(createBody)
	createReq, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token1)
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	require.Equal(ts.T(), http.StatusCreated, createW.Code)

	var createResponse map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createResponse)
	taskID := createResponse["id"].(float64)

	// 用戶2嘗試訪問用戶1的任務
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/tasks/%d", int(taskID)), nil)
	getReq.Header.Set("Authorization", "Bearer "+token2)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	require.Equal(ts.T(), http.StatusForbidden, getW.Code)
}

// TestTaskValidationWithContainers 測試任務驗證
func (ts *ContainerTestSuite) TestTaskValidationWithContainers() {
	token := ts.createTestUserAndLogin("validation@example.com", "password123", "validationuser")
	require.NotEmpty(ts.T(), token)

	// 測試缺少標題
	invalidBody := map[string]interface{}{
		"content": "No title",
		"status":  "pending",
	}
	invalidJSON, _ := json.Marshal(invalidBody)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	require.Equal(ts.T(), http.StatusBadRequest, w.Code)
}

// TestTaskCachingWithContainers 測試任務快取
func (ts *ContainerTestSuite) TestTaskCachingWithContainers() {
	token := ts.createTestUserAndLogin("cache@example.com", "password123", "cacheuser")
	require.NotEmpty(ts.T(), token)

	// 創建任務
	createBody := map[string]interface{}{
		"title":   "Cache Test Task",
		"content": "Cache Test Content",
		"status":  "pending",
	}
	createJSON, _ := json.Marshal(createBody)
	createReq, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	require.Equal(ts.T(), http.StatusCreated, createW.Code)

	// 第一次獲取任務列表（應該從數據庫）
	start := time.Now()
	listReq1, _ := http.NewRequest("GET", "/api/tasks", nil)
	listReq1.Header.Set("Authorization", "Bearer "+token)
	listW1 := httptest.NewRecorder()
	ts.router.ServeHTTP(listW1, listReq1)
	firstRequestTime := time.Since(start)

	require.Equal(ts.T(), http.StatusOK, listW1.Code)

	// 第二次獲取任務列表（應該從快取）
	start = time.Now()
	listReq2, _ := http.NewRequest("GET", "/api/tasks", nil)
	listReq2.Header.Set("Authorization", "Bearer "+token)
	listW2 := httptest.NewRecorder()
	ts.router.ServeHTTP(listW2, listReq2)
	secondRequestTime := time.Since(start)

	require.Equal(ts.T(), http.StatusOK, listW2.Code)

	// 驗證第二次請求應該更快（從快取）
	require.Less(ts.T(), secondRequestTime, firstRequestTime)
}
