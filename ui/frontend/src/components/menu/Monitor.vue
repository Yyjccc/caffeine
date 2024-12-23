<template>
  <div class="monitor-container">
    <div class="control-panel-wrapper">
      <el-card class="control-panel">
        <div class="control-items">
          <el-radio-group v-model="monitorTarget" @change="handleTargetChange">
            <el-radio-button label="local">本地监控</el-radio-button>
            <el-radio-button label="server">服务器监控</el-radio-button>
          </el-radio-group>
          
          <el-select v-model="refreshInterval" placeholder="刷新间隔" @change="handleIntervalChange">
            <el-option label="5秒" :value="5000" />
            <el-option label="10秒" :value="10000" />
            <el-option label="30秒" :value="30000" />
            <el-option label="1分钟" :value="60000" />
            <el-option label="15分钟" :value="900000" />
            <el-option label="1小时" :value="3600000" />
            <el-option label="3小时" :value="10800000" />
            <el-option label="6小时" :value="21600000" />
            <el-option label="12小时" :value="43200000" />
          </el-select>
        </div>
      </el-card>
    </div>

    <div class="monitor-content">
      <!-- 系统资源监控 -->
      <el-row :gutter="20" class="resource-row">
        <!-- CPU使用率 -->
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>CPU使用率</span>
              </div>
            </template>
            <div ref="cpuChartRef" class="chart"></div>
          </el-card>
        </el-col>
        
        <!-- 内存使用率 -->
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>内存使用率</span>
              </div>
            </template>
            <div ref="memoryChartRef" class="chart"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 网络监控 -->
      <el-row :gutter="20" class="network-row">
        <!-- 网络接口信息 -->
        <el-col :span="24">
          <el-card class="network-card">
            <template #header>
              <div class="card-header">
                <span>网络接口信息</span>
                <el-button type="primary" size="small" @click="refreshNetworkData">
                  刷新
                </el-button>
              </div>
            </template>
            <el-table :data="networkInterfaces" style="width: 100%" stripe>
              <el-table-column prop="name" label="接口名称" min-width="120" show-overflow-tooltip />
              <el-table-column prop="ipAddress" label="IP地址" min-width="140" show-overflow-tooltip />
              <el-table-column prop="macAddress" label="MAC地址" min-width="140" show-overflow-tooltip />
              <el-table-column prop="bytesReceived" label="接收流量" min-width="120">
                <template #default="scope">
                  {{ formatBytes(scope.row.bytesReceived) }}
                </template>
              </el-table-column>
              <el-table-column prop="bytesSent" label="发送流量" min-width="120">
                <template #default="scope">
                  {{ formatBytes(scope.row.bytesSent) }}
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <!-- 监听端口 -->
        <el-col :span="12">
          <el-card class="port-card">
            <template #header>
              <div class="card-header">
                <span>监听端口</span>
                <el-button type="primary" size="small" @click="refreshPortData">
                  刷新
                </el-button>
              </div>
            </template>
            <el-table :data="listeningPorts" style="width: 100%" stripe>
              <el-table-column prop="port" label="端口" min-width="100" show-overflow-tooltip />
              <el-table-column prop="protocol" label="协议" min-width="100" show-overflow-tooltip />
              <el-table-column prop="process" label="进程" min-width="150" show-overflow-tooltip />
              <el-table-column prop="state" label="状态" min-width="100">
                <template #default="scope">
                  <el-tag :type="scope.row.state === 'LISTEN' ? 'success' : 'info'">
                    {{ scope.row.state }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <!-- 活动连接 -->
        <el-col :span="12">
          <el-card class="connection-card">
            <template #header>
              <div class="card-header">
                <span>活动连接</span>
                <el-button type="primary" size="small" @click="refreshConnectionData">
                  刷新
                </el-button>
              </div>
            </template>
            <el-table 
              :data="activeConnections" 
              style="width: 100%" 
              stripe
            >
              <el-table-column prop="localAddress" label="本地地址" min-width="180" show-overflow-tooltip />
              <el-table-column prop="remoteAddress" label="远程地址" min-width="180" show-overflow-tooltip />
              <el-table-column prop="state" label="状态" min-width="120">
                <template #default="scope">
                  <el-tag :type="getConnectionStateType(scope.row.state)">
                    {{ scope.row.state }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { EChartOption } from 'echarts'
import { GetLocalSystemMetrics, GetNetworkInterfaces, GetListeningPorts, GetActiveConnections } from '../../../wailsjs/go/client/ClientApp'
// import type { EChartsOption } from 'echarts'

// 响应式状态
const monitorTarget = ref<'local' | 'server'>('local')
const refreshInterval = ref<number>(5000)
const cpuData = ref<[number, number][]>([])
const memoryData = ref<[number, number][]>([])
const networkInterfaces = ref<any[]>([])
const listeningPorts = ref<any[]>([])
const activeConnections = ref<any[]>([])

// Chart references
const cpuChartRef = ref<HTMLElement | null>(null)
const memoryChartRef = ref<HTMLElement | null>(null)
let cpuChart: echarts.ECharts | null = null
let memoryChart: echarts.ECharts | null = null

// 修改定时器类型定义
let chartTimer: ReturnType<typeof setInterval> | null = null
let networkTimer: ReturnType<typeof setInterval> | null = null

// 处理监控目标变化
const handleTargetChange = async () => {
  // 清空现有数据
  cpuData.value = []
  memoryData.value = []
  networkInterfaces.value = []
  listeningPorts.value = []
  activeConnections.value = []
  
  // 重新初始化数据
  await updateCharts()
  await refreshNetworkData()
  await refreshPortData()
  await refreshConnectionData()
}

// 处理窗口大小变化
const handleResize = () => {
  if (cpuChart) {
    cpuChart.resize()
  }
  if (memoryChart) {
    memoryChart.resize()
  }
}

// 处理刷新间隔变化
const handleIntervalChange = () => {
  if (chartTimer) {
    clearInterval(chartTimer)
    chartTimer = setInterval(updateCharts, refreshInterval.value)
  }
  if (networkTimer) {
    clearInterval(networkTimer)
    const networkRefreshInterval = Math.min(refreshInterval.value * 2, 300000)
    networkTimer = setInterval(() => {
      refreshNetworkData()
      refreshPortData()
      refreshConnectionData()
    }, networkRefreshInterval)
  }
  
  // 清空现有数据，重新开始收集
  cpuData.value = []
  memoryData.value = []
  updateCharts()
}

// 组件挂载
onMounted(() => {
  nextTick(() => {
    initCPUChart()
    initMemoryChart()
    updateCharts()
    
    refreshNetworkData()
    refreshPortData()
    refreshConnectionData()
    
    chartTimer = setInterval(updateCharts, refreshInterval.value)
    networkTimer = setInterval(() => {
      refreshNetworkData()
      refreshPortData()
      refreshConnectionData()
    }, refreshInterval.value * 2)
    
    window.addEventListener('resize', handleResize)
  })
})

// 组件卸载
onUnmounted(() => {
  if (chartTimer) {
    clearInterval(chartTimer)
    chartTimer = null
  }
  if (networkTimer) {
    clearInterval(networkTimer)
    networkTimer = null
  }
  
  window.removeEventListener('resize', handleResize)
  
  if (cpuChart) {
    cpuChart.dispose()
    cpuChart = null
  }
  if (memoryChart) {
    memoryChart.dispose()
    memoryChart = null
  }
})

// 更新图表数据保留时间
const updateCharts = async () => {
  try {
    const metrics = await GetLocalSystemMetrics()
    const now = new Date().getTime()
    
    // 数据验证
    if (typeof metrics.cpu !== 'number' || typeof metrics.memory !== 'number') {
      console.error('Invalid metrics data:', metrics)
      return
    }

    // 更新数据数组
    cpuData.value.push([now, Number(metrics.cpu.toFixed(2))])
    memoryData.value.push([now, Number(metrics.memory.toFixed(2))])

    // 根据刷新间隔调整数据保留时间
    // 保留数据时长为刷新间隔的3倍
    const retentionPeriod = refreshInterval.value * 3
    const cutoffTime = now - retentionPeriod

    // 清理旧数据
    cpuData.value = cpuData.value.filter(item => item[0] > cutoffTime)
    memoryData.value = memoryData.value.filter(item => item[0] > cutoffTime)

    // 更新图表
    if (cpuChart) {
      cpuChart.setOption({
        series: [{
          data: cpuData.value
        }]
      }, false, true)
    }

    if (memoryChart) {
      memoryChart.setOption({
        series: [{
          data: memoryData.value
        }]
      }, false, true)
    }
  } catch (error) {
    console.error('获取系统指标失败:', error)
  }
}

// 修改图表配置以适应不同的时间间隔
const getChartBaseOption = (title: string): EChartOption => {
  return {
    grid: {
      top: 60,
      right: 30,
      bottom: 60,
      left: 50,
      containLabel: true
    },
    title: {
      text: title,
      left: 'center',
      textStyle: {
        color: '#333',
      },
    },
    tooltip: {
      trigger: 'axis',
      formatter: function (params: any) {
        const time = new Date(params[0].value[0]).toLocaleTimeString()
        return `${time}<br/>使用率: ${params[0].value[1].toFixed(2)}%`
      },
    },
    xAxis: {
      type: 'time',
      splitLine: {
        show: false,
      },
      axisLabel: {
        rotate: 30,
        formatter: (value: number) => {
          const date = new Date(value)
          // 根据刷新间隔调整时间显示格式
          if (refreshInterval.value >= 3600000) {
            // 1小时及以上显示日期和时间
            return date.toLocaleString('zh-CN', {
              month: '2-digit',
              day: '2-digit',
              hour: '2-digit',
              minute: '2-digit'
            })
          } else if (refreshInterval.value >= 60000) {
            // 1分钟及以上显示时间（时:分）
            return date.toLocaleTimeString('zh-CN', {
              hour: '2-digit',
              minute: '2-digit'
            })
          } else {
            // 小于1分钟显示时间（时:分:秒）
            return date.toLocaleTimeString('zh-CN', {
              hour: '2-digit',
              minute: '2-digit',
              second: '2-digit'
            })
          }
        }
      }
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 100,
      splitLine: {
        lineStyle: {
          type: 'dashed',
        },
      },
      axisLabel: {
        formatter: '{value}%',
      },
    },
    series: [
      {
        name: 'CPU',
        type: 'line',
        smooth: true,
        data: cpuData.value,
        areaStyle: {
          opacity: 0.3,
        },
      },
    ],
  }
}

// 更新初始化函数
const initCPUChart = () => {
  if (!cpuChartRef.value) return
  cpuChart = echarts.init(cpuChartRef.value)
  const option = getChartBaseOption('CPU使用率')
  cpuChart.setOption(option)
}

const initMemoryChart = () => {
  if (!memoryChartRef.value) return
  memoryChart = echarts.init(memoryChartRef.value)
  const option = getChartBaseOption('内存使用率')
  memoryChart.setOption(option)
}

// 添加网络监控相关的方法
const refreshNetworkData = async () => {
  try {
    const data = await GetNetworkInterfaces()
    networkInterfaces.value = data
  } catch (error) {
    console.error('Failed to fetch network interfaces:', error)
  }
}

const refreshPortData = async () => {
  try {
    const data = await GetListeningPorts()
    listeningPorts.value = data
  } catch (error) {
    console.error('Failed to fetch listening ports:', error)
  }
}

const refreshConnectionData = async () => {
  try {
    const data = await GetActiveConnections()
    activeConnections.value = data
  } catch (error) {
    console.error('Failed to fetch active connections:', error)
  }
}

// 工具方法
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

const getConnectionStateType = (state: string): string => {
  const stateMap: { [key: string]: string } = {
    'ESTABLISHED': 'success',
    'TIME_WAIT': 'warning',
    'CLOSE_WAIT': 'info',
    'FIN_WAIT': 'warning',
    'CLOSED': 'danger'
  }
  return stateMap[state] || 'info'
}
</script>

<style scoped>
.monitor-container {
  width: 100%;
  height: 100%;
  position: relative;
  padding: 20px;
  overflow-y: auto;
}

.control-panel-wrapper {
  margin-bottom: 20px;
}

.control-panel {
  height: 45px;
}

.control-panel :deep(.el-card) {
  height: 100%;
  border: 1px solid #DCDFE6;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.control-panel :deep(.el-card__body) {
  height: 100%;
  padding: 0 15px;
  display: flex;
  align-items: center;
}

.control-items {
  display: flex;
  gap: 20px;
  align-items: center;
  width: 100%;
}

.monitor-content {
  width: 100%;
}

.resource-row {
  margin-bottom: 20px;
}

.network-row {
  margin-bottom: 20px;
}

.el-row {
  width: 100%;
  margin: 0 !important;
}

.chart-card, .network-card, .port-card, .connection-card {
  height: auto;
  margin-bottom: 20px;
  min-height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.chart {
  height: 320px;
  width: 100%;
  position: relative;
}

.chart > div {
  position: absolute !important;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

:deep(.el-card__body) {
  height: auto;
  padding: 15px;
  overflow: hidden;
}

:deep(.el-card) {
  height: 100%;
}

:deep(.el-table) {
  width: 100% !important;
}

:deep(.el-table__body-wrapper) {
  overflow-x: auto !important;
}

:deep(.el-table .cell) {
  white-space: normal;
  line-height: 1.5;
  word-break: break-all;
}

:deep(.el-table) {
  flex: 1;
}

/* 表格内容过多时显示滚动条 */
:deep(.el-table__body-wrapper) {
  overflow-y: auto;
}

.control-panel {
  margin-bottom: 20px;
  height: 45px;
}

.control-panel :deep(.el-card) {
  height: 100%;
  border: 1px solid #DCDFE6;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.control-panel :deep(.el-card__body) {
  height: 100%;
  padding: 0 15px;
  display: flex;
  align-items: center;
}

.control-items {
  display: flex;
  gap: 20px;
  align-items: center;
  width: 100%;
}

.control-items :deep(.el-select) {
  width: 120px;
}

/* 调整radio按钮组的样式 */
.control-items :deep(.el-radio-group) {
  margin: 0;
  line-height: 1;
}

/* 调整select下拉框的样式 */
.control-items :deep(.el-input__wrapper) {
  line-height: 1;
}

/* 调整单选按钮组的样式 */
.control-items :deep(.el-radio-button__inner) {
  transition: all 0.3s;
}

/* 选中状态的样式 */
.control-items :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background-color: #409EFF;
  border-color: #409EFF;
  color: white;
  box-shadow: -1px 0 0 0 #409EFF;
}

/* 未选中状态的样式 */
.control-items :deep(.el-radio-button__inner) {
  background-color: white;
  border-color: #DCDFE6;
  color: #606266;
}

/* 悬停效果 */
.control-items :deep(.el-radio-button__inner:hover) {
  color: #409EFF;
}

/* 选中状态悬停效果 */
.control-items :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner:hover) {
  color: white;
  background-color: #66b1ff;
}
</style>
