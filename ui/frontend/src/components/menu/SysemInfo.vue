<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElCard, ElDescriptions, ElDescriptionsItem, ElTag } from 'element-plus'
import { useRoute } from "vue-router"
import {SystemInfo} from "../../../bindings/caffeine/core";
import {WebShellSession} from "../../utils/session";
import store from "../../utils/store";

const systemInfo = ref<SystemInfo>()
const route = useRoute()
const id = Number(route.params.id);
// 添加一个格式化 PATH 的函数
const formatPath = (path: string) => {
  return path.split(';').join(';\n')
}

onMounted(() => {
  try {
    // 如果是字符串，先解析成对象
    const data = typeof route.query.systemInfo === 'string' 
      ? JSON.parse(route.query.systemInfo)
      : route.query.systemInfo
    
    systemInfo.value = data as SystemInfo
    //创建session
   var session:WebShellSession={
      sessionId:id,
      SystemType:systemInfo.value.systemType,
      SystemInfo:systemInfo.value
   }
    store.dispatch('createSession', session);

  } catch (error) {
    console.error('Failed to parse system info:', error)
  }
})
</script>

<template>
  <div class="system-info-container">
    <!-- 操作系统信息卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>操作系统信息</span>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="系统名称">{{ systemInfo?.os?.name }}</el-descriptions-item>
        <el-descriptions-item label="系统版本">{{ systemInfo?.os?.version }}</el-descriptions-item>
        <el-descriptions-item label="系统架构">{{ systemInfo?.os?.arch }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 系统路径信息卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>系统路径信息</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="当前根目录">{{ systemInfo?.currentFileRoot }}</el-descriptions-item>
        <el-descriptions-item label="当前目录">{{ systemInfo?.currentDir }}</el-descriptions-item>
        <el-descriptions-item label="所有根目录">{{ systemInfo?.fileRoots }}</el-descriptions-item>
        <el-descriptions-item label="临时目录">{{ systemInfo?.tempDirectory }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 系统用户信息卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>系统用户信息</span>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="当前用户">{{ systemInfo?.currentUser }}</el-descriptions-item>
        <el-descriptions-item label="进程架构">{{ systemInfo?.processArch }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- IP列表卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>IP地址列表</span>
        </div>
      </template>
      <div class="ip-list">
        <el-tag
          v-for="(ip, index) in systemInfo?.ipList"
          :key="index.toString()"
          class="ip-tag"
          type="info"
        >
          {{ ip }}
        </el-tag>
      </div>
    </el-card>

    <!-- 环境变量卡片 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>环境变量</span>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item
          v-for="(value, key) in systemInfo?.env"
          :key="String(key)"
          :label="String(key)"
        >
          <pre v-if="key === 'PATH'" class="path-content">{{ formatPath(String(value)) }}</pre>
          <span v-else>{{ value }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<style scoped>
.system-info-container {
  height: 100%;
  width: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-sizing: border-box;
  overflow: auto;
}

.info-card {
  width: 100%;
  flex-shrink: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ip-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.ip-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

:deep(.el-descriptions__body) {
  width: 100%;
}

:deep(.el-card__body) {
  padding: 20px;
}

.path-content {
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  font-family: inherit;
  max-height: 300px;
  overflow-y: auto;
}

:deep(.el-descriptions__cell) {
  word-break: break-all;
}

:deep(.el-card) {
  margin-bottom: 0;
}

:deep(.el-card__body) {
  padding: 15px;
}
</style>