package client

import (
	"fmt"
	"net"
	"syscall"

	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemMetric struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Time   int64   `json:"time"`
}

// NetworkInterface 网络接口信息
type NetworkInterface struct {
	Name          string `json:"name"`
	IPAddress     string `json:"ipAddress"`
	MACAddress    string `json:"macAddress"`
	BytesReceived uint64 `json:"bytesReceived"`
	BytesSent     uint64 `json:"bytesSent"`
}

// Port 端口信息
type Port struct {
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
	Process  string `json:"process"`
	State    string `json:"state"`
}

// Connection 连接信息
type Connection struct {
	LocalAddress  string `json:"localAddress"`
	RemoteAddress string `json:"remoteAddress"`
	State         string `json:"state"`
}

// GetNetworkInterfaces 获取网络接口信息
func (a *ClientApp) GetNetworkInterfaces() ([]NetworkInterface, error) {
	// 获取物理接口信息
	interfaces, err := psnet.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	// 获取网络IO统计信息
	ioStats, err := psnet.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get IO stats: %v", err)
	}

	// 创建名称到IO统计的映射
	ioStatsMap := make(map[string]psnet.IOCountersStat)
	for _, stat := range ioStats {
		ioStatsMap[stat.Name] = stat
	}

	var result []NetworkInterface
	for _, iface := range interfaces {
		// 跳过回环接口和没有运行的接口
		isLoopback := false
		isUp := false
		for _, flag := range iface.Flags {
			if flag == "loopback" {
				isLoopback = true
			}
			if flag == "up" {
				isUp = true
			}
		}

		// 跳过回环接口和没有运行的接口
		if isLoopback || !isUp {
			continue
		}

		// Get the addresses directly from the net package instead
		addrs, err := net.InterfaceByName(iface.Name)
		if err != nil {
			continue
		}
		ipAddrs, err := addrs.Addrs()
		if err != nil {
			continue
		}

		// Get the first IPv4 address
		var ipAddr string
		for _, addr := range ipAddrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ipAddr = ipNet.IP.String()
					break
				}
			}
		}

		// 获取IO统计信息
		var bytesRecv, bytesSent uint64
		if stat, ok := ioStatsMap[iface.Name]; ok {
			bytesRecv = stat.BytesRecv
			bytesSent = stat.BytesSent
		}

		result = append(result, NetworkInterface{
			Name:          iface.Name,
			IPAddress:     ipAddr,
			MACAddress:    iface.HardwareAddr,
			BytesReceived: bytesRecv,
			BytesSent:     bytesSent,
		})
	}

	return result, nil
}

// GetListeningPorts 获取监听端口信息
func (a *ClientApp) GetListeningPorts() ([]Port, error) {
	connections, err := psnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %v", err)
	}

	var result []Port
	seenPorts := make(map[string]bool)

	for _, conn := range connections {
		// 只关注监听状态的连接
		if conn.Status != "LISTEN" {
			continue
		}

		// 将 conn.Type (uint32) 转换为协议字符串
		protocol := "unknown"
		switch conn.Type {
		case syscall.SOCK_STREAM:
			protocol = "tcp"
		case syscall.SOCK_DGRAM:
			protocol = "udp"
		}

		// 创建端口标识符以避免重复
		portID := fmt.Sprintf("%d-%s", conn.Laddr.Port, protocol)
		if seenPorts[portID] {
			continue
		}
		seenPorts[portID] = true

		// 获取进程信息
		var processName string
		if conn.Pid > 0 {
			if proc, err := process.NewProcess(int32(conn.Pid)); err == nil {
				if name, err := proc.Name(); err == nil {
					processName = name
				}
			}
		}

		result = append(result, Port{
			Port:     uint16(conn.Laddr.Port),
			Protocol: protocol,
			Process:  processName,
			State:    conn.Status,
		})
	}

	return result, nil
}

// GetActiveConnections 获取活动连接信息
func (a *ClientApp) GetActiveConnections() ([]Connection, error) {
	connections, err := psnet.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %v", err)
	}

	var result []Connection
	for _, conn := range connections {
		// 跳过监听状态的连接
		if conn.Status == "LISTEN" {
			continue
		}

		// 格式化地址
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		remoteAddr := fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)

		result = append(result, Connection{
			LocalAddress:  localAddr,
			RemoteAddress: remoteAddr,
			State:         conn.Status,
		})
	}

	// 限制返回的连接数量，避免数据太多
	maxConnections := 100
	if len(result) > maxConnections {
		result = result[:maxConnections]
	}

	return result, nil
}
