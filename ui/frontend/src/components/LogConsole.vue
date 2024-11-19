<template>
  <div class="log-console">
    <div class="logs">
      <!-- 显示日志列表 -->
      <div v-for="(log, index) in logs" :key="index" class="log-entry">
         <span class="log-timestamp">
          {{ log.timestamp }}
        </span>
        <span :class="getLogLevelClass(log.level)">
          [{{ log.level }}]
        </span>
        <span class="log-funcname" :class="getFuncNameClass(log.funcName)">
          {{ log.funcName }}
        </span>
        <span class="log-message">
          {{ log.message }}
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, PropType } from 'vue';

// 接收父组件传递的 props
defineProps({
  logs: {
    type: Array as PropType<Array<{ level: string; timestamp: string; message: string,funcName:string }>>,
    required: true,
  },
  maxLines: {
    type: Number,
    default: 100,
  },
});

// 限制日志显示的最大行数（通过计算属性）
const displayedLogs =100;
// 获取函数名称样式
function getFuncNameClass(funcName: string): string {
  // 示例：不同函数名使用不同颜色
  if (funcName.includes('Log')) {
    return 'funcname-log';
  } else if (funcName.includes('Handle')) {
    return 'funcname-handle';
  } else {
    return 'funcname-default';
  }
}
// 获取日志级别样式
function getLogLevelClass(level: string): string {
  switch (level.toLowerCase()) {
    case 'debug':
      return 'log-level-debug';
    case 'info':
      return 'log-level-info';
    case 'warn':
      return 'log-level-warn';
    case 'error':
      return 'log-level-error';
    default:
      return 'log-level-default';
  }
}
</script>

<style scoped>
/* 容器样式 */
.log-console {
  position: fixed;
  bottom: 0;
  width: 100%;
  height: 150px;
  background-color: #1e1e1e;
  color: #dcdcdc;
  overflow-y: auto;
  font-family: monospace;
  padding: 10px;
  border-top: 1px solid #444;
}

.logs {
  display: flex;
  flex-direction: column;
}

.log-entry {
  margin: 2px 0;
}

/* 日志级别样式 */
.log-level-debug {
  color: #00bcd4; /* 蓝色 */
}

.log-level-info {
  color: #4caf50; /* 绿色 */
}

.log-level-warn {
  color: #ff9800; /* 橙色 */
}

.log-level-error {
  color: #f44336; /* 红色 */
}

.log-level-default {
  color: #dcdcdc; /* 默认白色 */
}

/* 日志内容样式 */
.log-message {
  color: #dcdcdc; /* 白色 */
}

.log-timestamp {
  color: #888; /* 灰色时间戳 */
  margin-right: 10px;
}
/* 函数名样式 */
.log-funcname {
  font-weight: bold;
  margin-right: 10px;
}

/* 特定函数名样式 */
.funcname-log {
  color: #ffcc00; /* 黄色 */
}

.funcname-handle {
  color: #009688; /* 青色 */
}

.funcname-default {
  color: #dcdcdc; /* 默认白色 */
}
</style>
