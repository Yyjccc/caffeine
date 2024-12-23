package client

import (
	"caffeine/client/c2"
	"caffeine/client/webshell"
	"caffeine/core"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type WebShellManger struct {
	alive   []int64 // 存放在线shell的ID
	clients map[int64]*webshell.WebClient
	entries map[int64]*ShellEntry // 存储 ShellEntry 实体
	db      *gorm.DB              // 添加数据库实例
}

// NewWebShellManager 创建管理器实例
func NewWebShellManager(db *gorm.DB) *WebShellManger {
	return &WebShellManger{
		alive:   make([]int64, 0),
		clients: make(map[int64]*webshell.WebClient),
		entries: make(map[int64]*ShellEntry),
		db:      db,
	}
}

// AddEntry 添加 ShellEntry
func (m *WebShellManger) AddEntry(entry *ShellEntry) {
	m.entries[entry.ID] = entry
	// 同时创建对应的 WebClient
	client := entry.ToWebClient()
	m.clients[client.ID] = client
}

// RemoveEntry 根据ID移除 ShellEntry
func (m *WebShellManger) RemoveEntry(id int64) {
	delete(m.entries, id)
	delete(m.clients, int64(id))
	// 同时从在线列表中移除
	for i, aliveID := range m.alive {
		if aliveID == id {
			m.alive = append(m.alive[:i], m.alive[i+1:]...)
			break
		}
	}
}

// GetEntry 根据ID获取 ShellEntry
func (m *WebShellManger) GetEntry(id int64) *ShellEntry {
	return m.entries[id]
}

// UpdateEntry 更新 ShellEntry
func (m *WebShellManger) UpdateEntry(entry *ShellEntry) {
	if _, exists := m.entries[entry.ID]; exists {
		m.entries[entry.ID] = entry
		// 同步更新 WebClient
		client := entry.ToWebClient()
		m.clients[client.ID] = client
	}
}

// SetEntryStatus 设置 ShellEntry 状态
func (m *WebShellManger) SetEntryStatus(id int64, status int) {
	if entry, exists := m.entries[id]; exists {
		entry.Status = status
		if status == 1 { // 在线
			m.SetAlive(int64(id))
		}
	}
}

// SetAlive adds the shell ID to the alive list
func (m *WebShellManger) SetAlive(id int64) {
	m.alive = append(m.alive, id)
}

func (*WebShellManger) name() {

}

func (m *WebShellManger) AddWebShell(client *webshell.WebClient) {
	m.clients[client.ID] = client
}

// ShellEntry 表示一条 webshell 基础数据记录
type ShellEntry struct {
	ID         int64  `gorm:"primaryKey"`
	Location   string // shell 所在位置
	ShellType  string // shell 类型(php/jsp等)
	IP         string
	CreateTime string
	UpdateTime string
	URL        string
	Note       string
	Password   string // shell 密码
	Encoding   string // 编码方式(如 base64)
	Status     int    // 状态: 0-离线 1-在线
}

// ToWebClient converts ShellEntry to WebClient
func (e *ShellEntry) ToWebClient() *webshell.WebClient {
	target := core.Target{
		ShellURL: e.URL,
	}
	// 使用默认的 C2 配置创建新的 WebClient
	config := c2.C2Yaml{} // 如果需要特定配置，可以从配置文件或其他地方加载
	return webshell.NewWebClient(target, config)
}

// AddNewShell 添加新的WebShell并保存到数据库
func (m *WebShellManger) AddNewShell(data map[string]interface{}) (int64, error) {
	// 创建新的ShellEntry
	entry := &ShellEntry{
		Location:   data["location"].(string),
		ShellType:  data["shellType"].(string),
		IP:         data["ip"].(string),
		URL:        data["url"].(string),
		Note:       data["note"].(string),
		Password:   data["password"].(string),
		Encoding:   data["encoding"].(string),
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     0, // 默认离线状态
	}

	// 保存到数据库
	result := m.db.Create(entry)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to save shell entry: %v", result.Error)
	}

	// 添加到内存管理
	m.AddEntry(entry)

	return entry.ID, nil
}

// GetShellList 从数据库获取shell列表
func (m *WebShellManger) GetShellList() ([]ShellEntry, error) {
	var entries []ShellEntry

	// 从数据库中查询所有记录
	result := m.db.Find(&entries)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch shell entries: %v", result.Error)
	}

	// 更新内存中的entries
	for _, entry := range entries {
		m.entries[entry.ID] = &entry
	}

	return entries, nil
}
