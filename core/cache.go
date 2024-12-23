package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CacheManager struct {
	db              *gorm.DB
	mu              sync.RWMutex
	sessionCache    map[int64]*Session
	systemInfoCache map[int64]*SystemInfo
	fileSystemCache map[string]*FileSystemCache
	directories     map[string]*Directory
}

// Cache models
type SessionCache struct {
	ID             int64  `gorm:"primarykey"`
	OperateHistory string // JSON serialized
	OutputHistory  string // JSON serialized
	StartTime      time.Time
	LastActive     time.Time
	TargetURL      string
	SystemInfoID   int64
}

type SystemInfoCache struct {
	ID            int64 `gorm:"primarykey"`
	FileRoot      string
	CurrentDir    string
	CurrentUser   string
	ProcessArch   int
	TempDirectory string
	IpList        string // JSON serialized
	OsInfo        string // JSON serialized
	Env           string // JSON serialized
}

// HttpCache 数据库模型
type HttpCache struct {
	ID         int64 `gorm:"primarykey"`
	RequestID  int64 `gorm:"index"`
	Method     string
	URL        string
	ReqHeader  json.RawMessage
	ReqBody    []byte
	RespCode   int
	RespHeader json.RawMessage
	RespBody   []byte
	CreatedAt  time.Time
}

var (
	cacheManager *CacheManager
	once         sync.Once
)

func GetCacheManager() *CacheManager {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open("cache.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Auto migrate schemas
		db.AutoMigrate(&SessionCache{}, &SystemInfoCache{}, &HttpCache{})

		cacheManager = &CacheManager{
			db:              db,
			sessionCache:    make(map[int64]*Session),
			systemInfoCache: make(map[int64]*SystemInfo),
			fileSystemCache: make(map[string]*FileSystemCache),
			directories:     make(map[string]*Directory),
		}
	})
	return cacheManager
}

// Session cache methods
func (cm *CacheManager) SaveSession(session *Session) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Cache in memory
	cm.sessionCache[session.ID] = session

	// Save to database
	operateHistory, _ := json.Marshal(session.OperateHistory)
	outputHistory, _ := json.Marshal(session.OutputHistory)

	sessionCache := SessionCache{
		ID:             session.ID,
		OperateHistory: string(operateHistory),
		OutputHistory:  string(outputHistory),
		StartTime:      session.StartTime,
		LastActive:     session.LastActive,
		TargetURL:      session.Target.ShellURL,
		SystemInfoID:   session.Info.ID,
	}

	return cm.db.Save(&sessionCache).Error
}

func (cm *CacheManager) GetSession(id int64) (*Session, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Try memory cache first
	if session, ok := cm.sessionCache[id]; ok {
		return session, nil
	}

	// Load from database
	var sessionCache SessionCache
	if err := cm.db.First(&sessionCache, id).Error; err != nil {
		return nil, err
	}

	// Deserialize and construct Session
	var operateHistory []Operate
	var outputHistory []string
	json.Unmarshal([]byte(sessionCache.OperateHistory), &operateHistory)
	json.Unmarshal([]byte(sessionCache.OutputHistory), &outputHistory)

	session := &Session{
		ID:             sessionCache.ID,
		OperateHistory: operateHistory,
		OutputHistory:  outputHistory,
		StartTime:      sessionCache.StartTime,
		LastActive:     sessionCache.LastActive,
		Target: Target{
			ShellURL: sessionCache.TargetURL,
		},
	}

	// Cache in memory
	cm.sessionCache[id] = session
	return session, nil
}

// SystemInfo cache methods
func (cm *CacheManager) SaveSystemInfo(info *SystemInfo) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Cache in memory
	cm.systemInfoCache[info.ID] = info

	// Convert env map[string]string to map[string]interface{} for serialization
	envInterface := make(map[string]interface{})
	for k, v := range info.Env {
		envInterface[k] = v
	}

	// Save to database
	ipList, _ := json.Marshal(info.IpList)
	osInfo, _ := json.Marshal(info.Os)
	env, _ := json.Marshal(envInterface)

	systemInfoCache := SystemInfoCache{
		ID:            info.ID,
		FileRoot:      info.FileRoot,
		CurrentDir:    info.CurrentDir,
		CurrentUser:   info.CurrentUser,
		ProcessArch:   info.ProcessArch,
		TempDirectory: info.TempDirectory,
		IpList:        string(ipList),
		OsInfo:        string(osInfo),
		Env:           string(env),
	}

	return cm.db.Save(&systemInfoCache).Error
}

func (cm *CacheManager) GetSystemInfo(id int64) (*SystemInfo, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Try memory cache first
	if info, ok := cm.systemInfoCache[id]; ok {
		return info, nil
	}

	// Load from database
	var systemInfoCache SystemInfoCache
	if err := cm.db.First(&systemInfoCache, id).Error; err != nil {
		return nil, err
	}

	// Deserialize data
	var ipList []string
	var osInfo OSInfo
	var envTemp map[string]interface{}
	json.Unmarshal([]byte(systemInfoCache.IpList), &ipList)
	json.Unmarshal([]byte(systemInfoCache.OsInfo), &osInfo)
	json.Unmarshal([]byte(systemInfoCache.Env), &envTemp)

	// Convert env map to correct type
	env := make(map[string]string)
	for k, v := range envTemp {
		if str, ok := v.(string); ok {
			env[k] = str
		}
	}

	// Convert back to map[string]interface{} for struct assignment
	envInterface := make(map[string]interface{})
	for k, v := range env {
		envInterface[k] = v
	}

	// Construct SystemInfo
	info := &SystemInfo{
		ID:            systemInfoCache.ID,
		FileRoot:      systemInfoCache.FileRoot,
		CurrentDir:    systemInfoCache.CurrentDir,
		CurrentUser:   systemInfoCache.CurrentUser,
		ProcessArch:   systemInfoCache.ProcessArch,
		TempDirectory: systemInfoCache.TempDirectory,
		IpList:        ipList,
		Os:            osInfo,
		Env:           envInterface, // Use the converted map
	}

	// Cache in memory
	cm.systemInfoCache[id] = info
	return info, nil
}

// Add this method to the CacheManager struct
func (c *CacheManager) GetDirectory(path string) (*Directory, error) {
	dir, exists := c.directories[path]
	if !exists {
		return nil, fmt.Errorf("directory not found in cache: %s", path)
	}
	return dir, nil
}

// Add this method to the CacheManager struct
func (c *CacheManager) SaveDirectory(dir *Directory) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.directories[dir.Path] = dir
	return nil
}

// SaveToCache 保存 HTTP 请求到缓存
func (cm *CacheManager) SaveToCache(req *HttpRequest) error {
	reqHeaders, _ := json.Marshal(req.Headers)
	respHeaders, _ := json.Marshal(req.Response.Headers)

	cache := HttpCache{
		RequestID:  req.ID,
		Method:     req.Method,
		URL:        req.URL,
		ReqHeader:  reqHeaders,
		ReqBody:    req.Body,
		RespCode:   req.Response.code,
		RespHeader: respHeaders,
		RespBody:   req.Response.Body,
		CreatedAt:  time.Now(),
	}

	return cm.db.Create(&cache).Error
}

// GetFromCache 从缓存获取 HTTP 请求
func (cm *CacheManager) GetFromCache(requestID int64) (*HttpRequest, error) {
	var cache HttpCache
	if err := cm.db.Where("request_id = ?", requestID).First(&cache).Error; err != nil {
		return nil, err
	}

	var headers map[string]string
	json.Unmarshal(cache.ReqHeader, &headers)

	req := &HttpRequest{
		ID:      cache.RequestID,
		Method:  cache.Method,
		URL:     cache.URL,
		Headers: headers,
		Body:    cache.ReqBody,
		Response: &HttpResponse{
			code: cache.RespCode,
			Body: cache.RespBody,
		},
	}
	json.Unmarshal(cache.RespHeader, &req.Response.Headers)

	return req, nil
}

// CloseAndCleanup closes the database connection and cleans up resources
func (cm *CacheManager) CloseAndCleanup() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Clear in-memory caches
	cm.sessionCache = make(map[int64]*Session)
	cm.systemInfoCache = make(map[int64]*SystemInfo)
	cm.fileSystemCache = make(map[string]*FileSystemCache)
	cm.directories = make(map[string]*Directory)

	// Close database connection
	sqlDB, err := cm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	return sqlDB.Close()
}
