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
  <div class="common-layout">
    <el-container>
      <el-header>
        <Header ref="tabBar" />
      </el-header>
      <el-main>
        <router-view></router-view>
        <!-- 传递日志到 LogConsole 组件 -->
        <LogConsole :logs="logs" :maxLines="100" />
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
.main {
  margin: 0;
  padding: 0;
}
</style>
