<script lang="ts" setup>
import { ref, onMounted, provide } from "vue";
import Header from "./components/aside/Header.vue";
import LogConsole from "./components/LogConsole.vue";
import { EventsOn } from "../wailsjs/runtime";

// 定义 tabBar 引用
const tabBar = ref<{ addTab: (name: string, path: string) => void } | null>(null);
provide("tabBar", tabBar);

// 定义日志数组，包含日志级别、时间戳、函数名和消息
const logs = ref<{ level: string; timestamp: string; funcName: string; message: string }[]>([]);

// 最大日志行数
const maxLines = 100;

// 监听后端日志事件
onMounted(() => {
  if (tabBar.value) {
    provide('tabBar', tabBar.value);
    console.log("tabBar  bound");
  }
  // 监听 "log" 事件，接收后端 JSON 格式的日志
  EventsOn("log", (logData: string) => {
    try {
      // 解析 JSON 数据
      const parsedLog = JSON.parse(logData);

      // 检查必需字段是否存在
      if (parsedLog.level && parsedLog.time && parsedLog.funcName && parsedLog.message) {
        // 如果日志超过最大行数，删除最旧的日志
        if (logs.value.length >= maxLines) {
          logs.value.shift();
        }
        // 将新日志添加到日志数组
        logs.value.push({
          level: parsedLog.level,
          timestamp: parsedLog.time,
          funcName: parsedLog.funcName,
          message: parsedLog.message,
        });
      }
    } catch (error) {
      console.error("日志解析失败:", error, "原始日志数据:", logData);
    }
  });
});
</script>

<template>
  <div class="app-container">
    <el-container class="main-layout">
      <el-header height="50px">
        <Header ref="tabBar" />
      </el-header>
      
      <el-container class="content-layout">
        <el-main class="main-content">
          <router-view></router-view>
        </el-main>
        
        <el-footer height="200px" class="log-footer">
          <LogConsole :logs="logs" :maxLines="100" />
        </el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.main-layout {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.content-layout {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0; /* 重要：允许flex子项收缩 */
}

.main-content {
  flex: 1;
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.log-footer {
  padding: 0;
  border-top: 1px solid #ddd;
  background-color: #f5f7fa;
}

:deep(.el-footer) {
  padding: 0;
  border-top: 1px solid #ddd;
  background-color: #f5f7fa;
}
</style>
