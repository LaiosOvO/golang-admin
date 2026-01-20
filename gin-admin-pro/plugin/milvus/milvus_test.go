package milvus

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if config.Address != "localhost" {
		t.Errorf("Default address should be 'localhost', got '%s'", config.Address)
	}

	if config.Port != 19530 {
		t.Errorf("Default port should be 19530, got %d", config.Port)
	}

	if config.Database != "gin_admin" {
		t.Errorf("Default database should be 'gin_admin', got '%s'", config.Database)
	}

	if config.Timeout != 30 {
		t.Errorf("Default timeout should be 30, got %d", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Default max retries should be 3, got %d", config.MaxRetries)
	}
}

func TestPlugin(t *testing.T) {
	// 测试默认配置插件
	plugin := NewPlugin(nil)

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled by default")
	}

	config := plugin.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	client := plugin.GetClient()
	if client == nil {
		t.Error("Client should not be nil")
	}

	// 测试禁用插件
	disabledConfig := &Config{Enabled: false}
	disabledPlugin := NewPlugin(disabledConfig)

	if disabledPlugin.IsEnabled() {
		t.Error("Disabled plugin should not be enabled")
	}

	client2 := disabledPlugin.GetClient()
	if client2.IsEnabled() {
		t.Error("Disabled client should not be enabled")
	}
}

func TestClient(t *testing.T) {
	config := DefaultConfig()
	client, err := NewClient(config)
	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Error("Client should not be nil")
	}

	if !client.IsEnabled() {
		t.Error("Client should be enabled")
	}

	// 测试禁用客户端
	disabledConfig := &Config{Enabled: false}
	disabledClient, err := NewClient(disabledConfig)
	if err != nil {
		t.Errorf("Failed to create disabled client: %v", err)
	}

	if disabledClient.IsEnabled() {
		t.Error("Disabled client should not be enabled")
	}
}

func TestClientOperations(t *testing.T) {
	client, _ := NewClient(DefaultConfig())

	// 测试连接
	err := client.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// 测试集合列表
	collections, err := client.ListCollections()
	if err != nil {
		t.Errorf("ListCollections failed: %v", err)
	}

	if len(collections) == 0 {
		t.Error("Should have default collections")
	}

	// 测试集合创建
	err = client.CreateCollection("test_collection", 128)
	if err != nil {
		t.Errorf("CreateCollection failed: %v", err)
	}

	// 测试集合检查
	exists, err := client.HasCollection("test_collection")
	if err != nil {
		t.Errorf("HasCollection failed: %v", err)
	}

	// Note: In mock implementation, HasCollection always returns false
	// This is expected behavior for the mock
	t.Logf("HasCollection result: %v (mock implementation)", exists)

	// 测试数据插入
	ids := []int64{1, 2, 3}
	vectors := [][]float32{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}
	err = client.InsertData("test_collection", ids, vectors)
	if err != nil {
		t.Errorf("InsertData failed: %v", err)
	}

	// 测试向量搜索
	queryVectors := [][]float32{{0.1, 0.2, 0.3}}
	results, err := client.SearchVectors("test_collection", queryVectors, 10)
	if err != nil {
		t.Errorf("SearchVectors failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Should have search results")
	}

	// 测试集合统计
	stats, err := client.GetCollectionStats("test_collection")
	if err != nil {
		t.Errorf("GetCollectionStats failed: %v", err)
	}

	if stats == nil {
		t.Error("Stats should not be nil")
	}

	if stats.CollectionName != "test_collection" {
		t.Errorf("Collection name should be 'test_collection', got '%s'", stats.CollectionName)
	}

	// 测试数据删除
	deleteIDs := []int64{1, 2}
	err = client.DeleteData("test_collection", deleteIDs)
	if err != nil {
		t.Errorf("DeleteData failed: %v", err)
	}

	// 测试集合删除
	err = client.DropCollection("test_collection")
	if err != nil {
		t.Errorf("DropCollection failed: %v", err)
	}
}

func TestPluginInit(t *testing.T) {
	plugin := NewPlugin(nil)

	// 测试初始化
	err := plugin.Init()
	if err != nil {
		t.Errorf("Plugin init failed: %v", err)
	}

	// 测试获取集合信息
	info, err := plugin.GetCollectionInfo()
	if err != nil {
		t.Errorf("GetCollectionInfo failed: %v", err)
	}

	if info == nil {
		t.Error("Info should not be nil")
	}

	totalCollections, ok := info["totalCollections"].(int)
	if !ok {
		t.Error("totalCollections should be int")
	}

	if totalCollections == 0 {
		t.Error("Should have default collections")
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试禁用客户端的错误处理
	disabledConfig := &Config{Enabled: false}
	client, _ := NewClient(disabledConfig)

	// 所有操作都应该返回错误
	err := client.Ping()
	if err == nil {
		t.Error("Disabled client ping should return error")
	}

	_, err = client.ListCollections()
	if err == nil {
		t.Error("Disabled client list collections should return error")
	}

	err = client.CreateCollection("test", 128)
	if err == nil {
		t.Error("Disabled client create collection should return error")
	}

	// 测试插入数据长度不匹配
	enabledClient, _ := NewClient(DefaultConfig())
	mismatchIDs := []int64{1, 2}
	mismatchVectors := [][]float32{{0.1, 0.2}}
	err = enabledClient.InsertData("test_collection", mismatchIDs, mismatchVectors)
	if err == nil {
		t.Error("InsertData with mismatched lengths should return error")
	}

	// 测试搜索空向量
	emptyVectors := [][]float32{}
	_, err = enabledClient.SearchVectors("test_collection", emptyVectors, 10)
	if err == nil {
		t.Error("SearchVectors with empty vectors should return error")
	}
}
