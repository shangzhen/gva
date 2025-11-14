package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// 测试配置
const (
	baseURL = "http://localhost:8888"
	token   = "your-jwt-token-here" // 需要先登录获取token
)

// 测试用例
func main() {
	fmt.Println("====== Elasticsearch 用户操作日志测试 ======\n")

	// 1. 初始化索引
	fmt.Println("1. 初始化索引...")
	if err := initIndex(); err != nil {
		fmt.Printf("❌ 初始化索引失败: %v\n", err)
		return
	}
	fmt.Println("✅ 初始化索引成功\n")
	time.Sleep(2 * time.Second)

	// 2. 创建单条日志
	fmt.Println("2. 创建单条日志...")
	logID := ""
	if id, err := createLog(); err != nil {
		fmt.Printf("❌ 创建日志失败: %v\n", err)
		return
	} else {
		logID = id
		fmt.Printf("✅ 创建日志成功，ID: %s\n\n", logID)
	}
	time.Sleep(2 * time.Second)

	// 3. 批量创建日志
	fmt.Println("3. 批量创建测试数据...")
	if err := batchCreateLogs(50); err != nil {
		fmt.Printf("❌ 批量创建日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 批量创建50条日志成功\n")
	}
	time.Sleep(3 * time.Second)

	// 4. 获取单条日志
	if logID != "" {
		fmt.Printf("4. 获取单条日志 (ID: %s)...\n", logID)
		if err := getLog(logID); err != nil {
			fmt.Printf("❌ 获取日志失败: %v\n", err)
		} else {
			fmt.Println("✅ 获取日志成功\n")
		}
	}

	// 5. 搜索所有日志
	fmt.Println("5. 搜索所有日志（第1页，每页10条）...")
	if err := searchLogs(map[string]interface{}{
		"page":     1,
		"pageSize": 10,
	}); err != nil {
		fmt.Printf("❌ 搜索日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 搜索日志成功\n")
	}

	// 6. 按用户ID搜索
	fmt.Println("6. 按用户ID搜索（user_id=1）...")
	if err := searchLogs(map[string]interface{}{
		"page":     1,
		"pageSize": 10,
		"user_id":  1,
	}); err != nil {
		fmt.Printf("❌ 搜索日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 按用户ID搜索成功\n")
	}

	// 7. 按操作动作搜索
	fmt.Println("7. 按操作动作搜索（action=login）...")
	if err := searchLogs(map[string]interface{}{
		"page":     1,
		"pageSize": 10,
		"action":   "login",
	}); err != nil {
		fmt.Printf("❌ 搜索日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 按操作动作搜索成功\n")
	}

	// 8. 按模块搜索
	fmt.Println("8. 按模块搜索（module=user）...")
	if err := searchLogs(map[string]interface{}{
		"page":     1,
		"pageSize": 10,
		"module":   "user",
	}); err != nil {
		fmt.Printf("❌ 搜索日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 按模块搜索成功\n")
	}

	// 9. 按时间范围搜索
	fmt.Println("9. 按时间范围搜索（最近1小时）...")
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour)
	if err := searchLogs(map[string]interface{}{
		"page":       1,
		"pageSize":   10,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}); err != nil {
		fmt.Printf("❌ 搜索日志失败: %v\n", err)
	} else {
		fmt.Println("✅ 按时间范围搜索成功\n")
	}

	// 10. 获取统计数据
	fmt.Println("10. 获取统计数据（按action分组）...")
	if err := getStats("action"); err != nil {
		fmt.Printf("❌ 获取统计失败: %v\n", err)
	} else {
		fmt.Println("✅ 获取统计成功\n")
	}

	fmt.Println("====== 测试完成 ======")
}

// 初始化索引
func initIndex() error {
	url := fmt.Sprintf("%s/userActionLog/initIndex", baseURL)
	_, err := makeRequest("POST", url, nil)
	return err
}

// 创建日志
func createLog() (string, error) {
	url := fmt.Sprintf("%s/userActionLog/createLog", baseURL)
	data := map[string]interface{}{
		"user_id":    1,
		"username":   "admin",
		"action":     "login",
		"module":     "user",
		"method":     "POST",
		"path":       "/base/login",
		"ip":         "127.0.0.1",
		"user_agent": "Mozilla/5.0",
		"status":     200,
		"latency":    123,
	}

	resp, err := makeRequest("POST", url, data)
	if err != nil {
		return "", err
	}

	// 这里简化处理，实际应该从响应中提取ID
	return "test-id", nil
}

// 批量创建日志
func batchCreateLogs(count int) error {
	actions := []string{"login", "logout", "create", "update", "delete", "view", "export"}
	modules := []string{"user", "role", "menu", "api", "system", "log"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	usernames := []string{"admin", "user1", "user2", "test"}

	for i := 0; i < count; i++ {
		action := actions[rand.Intn(len(actions))]
		module := modules[rand.Intn(len(modules))]
		method := methods[rand.Intn(len(methods))]
		username := usernames[rand.Intn(len(usernames))]
		userID := rand.Intn(4) + 1
		status := []int{200, 201, 400, 401, 403, 404, 500}[rand.Intn(7)]

		url := fmt.Sprintf("%s/userActionLog/createLog", baseURL)
		data := map[string]interface{}{
			"user_id":    userID,
			"username":   username,
			"action":     action,
			"module":     module,
			"method":     method,
			"path":       fmt.Sprintf("/%s/%s", module, action),
			"ip":         fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			"user_agent": "Mozilla/5.0 Test",
			"status":     status,
			"latency":    rand.Int63n(1000),
		}

		if _, err := makeRequest("POST", url, data); err != nil {
			return fmt.Errorf("创建第%d条日志失败: %v", i+1, err)
		}

		// 避免请求过快
		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

// 获取日志
func getLog(id string) error {
	url := fmt.Sprintf("%s/userActionLog/getLog/%s", baseURL, id)
	resp, err := makeRequest("GET", url, nil)
	if err != nil {
		return err
	}

	fmt.Printf("日志详情: %s\n", string(resp))
	return nil
}

// 搜索日志
func searchLogs(params map[string]interface{}) error {
	url := fmt.Sprintf("%s/userActionLog/searchLogs", baseURL)
	resp, err := makeRequest("POST", url, params)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		total := data["total"]
		fmt.Printf("搜索结果: 共 %v 条记录\n", total)
	}

	return nil
}

// 获取统计
func getStats(groupBy string) error {
	url := fmt.Sprintf("%s/userActionLog/getStats", baseURL)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	data := map[string]interface{}{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
		"group_by":   groupBy,
	}

	resp, err := makeRequest("POST", url, data)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		total := data["total"]
		fmt.Printf("统计结果: 总数 %v\n", total)
		if stats, ok := data["stats"].([]interface{}); ok {
			fmt.Printf("分组统计: %d 个分组\n", len(stats))
			for i, stat := range stats {
				if s, ok := stat.(map[string]interface{}); ok {
					key := s["key"]
					count := s["doc_count"]
					fmt.Printf("  %d. %s: %v 条\n", i+1, key, count)
				}
			}
		}
	}

	return nil
}

// 发送HTTP请求
func makeRequest(method, url string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("x-token", token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
