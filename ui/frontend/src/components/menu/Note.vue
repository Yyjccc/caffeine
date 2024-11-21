<script setup lang="ts">
import { computed, ref, onMounted, watch, h } from 'vue'
import { useStore } from 'vuex'
import { key } from '../../utils/store'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import { ElMessage } from 'element-plus'
import {
  DocumentCopy,
  Minus,
  Link,
  Picture,
  List,
  Tickets,
  Check,
  Monitor,
  Document,
  View,
  ArrowDown,
  Promotion,
  Grid,
  Reading,
  Edit
} from '@element-plus/icons-vue'

// 添加光标位置状态
const currentLine = ref(1)
const currentColumn = ref(1)

// 添加右键菜单状态和方法
const closeContextMenu = () => {
  showContextMenu.value = false
}

const handleMenuItemClick = (action: () => void) => {
  action()
  closeContextMenu()
}

// 明确指定 MarkdownIt 类型
const md: MarkdownIt = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true,
  highlight: function (str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        const highlighted = hljs.highlight(str, { language: lang, ignoreIllegals: true }).value
        return `
          <div class="code-block-wrapper">
            <div class="code-block-header">
              <span class="code-lang">${lang}</span>
              <button class="copy-button" onclick="window.copyCode('${encodeURIComponent(str)}')">
                复制代码
              </button>
            </div>
            <pre class="code-block"><code class="hljs language-${lang}">${highlighted}</code></pre>
          </div>`
      } catch (__) {}
    }
    return `<pre class="code-block"><code class="hljs">${md.utils.escapeHtml(str)}</code></pre>`
  }
})

// 扩展 Window 接口
declare global {
  interface Window {
    copyCode: (code: string) => void;
  }
}

// 自动保存相关
const saveTimeout = ref<ReturnType<typeof setTimeout> | null>(null)

// 工具栏操作类型定义
type ToolbarAction = keyof typeof toolbarActions
const handleCommand = (command: ToolbarAction) => {
  toolbarActions[command]()
}

// 更新光标位置
const updateCursorPosition = (e: Event) => {
  const textarea = e.target as HTMLTextAreaElement
  const value = textarea.value.substr(0, textarea.selectionStart)
  const lines = value.split('\n')
  currentLine.value = lines.length
  currentColumn.value = lines[lines.length - 1].length + 1
}

const store = useStore(key)
const editorContent = ref('')
const showContextMenu = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const isSaved = ref(true)

// 实时渲染的内容
const renderedContent = computed(() => {
  return md.render(editorContent.value || '')
})

// 自动保存延迟时间（毫秒）
const AUTO_SAVE_DELAY = 1000

// 保存笔记
const saveNote = () => {
  store.dispatch('saveNote', editorContent.value)
  isSaved.value = true
}

// 自动保存（防抖）
const autoSave = () => {
  isSaved.value = false
  if (saveTimeout.value) {
    clearTimeout(saveTimeout.value)
  }
  saveTimeout.value = setTimeout(() => {
    saveNote()
  }, AUTO_SAVE_DELAY)
}

// 监听内容变化
watch(editorContent, () => {
  autoSave()
}, { deep: true })

// 组件挂载时加载笔记
onMounted(() => {
  store.dispatch('loadNote')
  editorContent.value = store.state.noteContent
})

const editorRef = ref<HTMLTextAreaElement | null>(null)

// 插入文本的通用函数
const insertText = (before: string, after: string = '') => {
  if (!editorRef.value) return
  
  const textarea = editorRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = textarea.value
  
  const selectedText = text.substring(start, end)
  const replacement = before + selectedText + after
  
  textarea.value = text.substring(0, start) + replacement + text.substring(end)
  
  // 更新光标位置
  const newCursorPos = selectedText ? start + replacement.length : start + before.length
  textarea.focus()
  textarea.setSelectionRange(newCursorPos, newCursorPos)
  
  // 触发内容更新
  editorContent.value = textarea.value
}

// 工具栏操作
const toolbarActions = {
  // 基础格式
  bold: () => insertText('**', '**'),
  italic: () => insertText('*', '*'),
  strikethrough: () => insertText('~~', '~~'),
  
  // 标题
  h1: () => insertText('# '),
  h2: () => insertText('## '),
  h3: () => insertText('### '),
  
  // 列表
  bulletList: () => insertText('- '),
  numberList: () => insertText('1. '),
  taskList: () => insertText('- [ ] '),
  
  // 代码
  inlineCode: () => insertText('`', '`'),
  codeBlock: () => insertText('```\n', '\n```'),
  
  // 引用和分割线
  quote: () => insertText('> '),
  divider: () => insertText('\n---\n'),
  
  // 链接和图片
  link: () => insertText('[', '](url)'),
  image: () => insertText('![alt text](', ')'),
  
  // 表格
  table: () => insertText('\n| Header 1 | Header 2 |\n| --------- | --------- |\n| Cell 1 | Cell 2 |\n'),
  
  // 保存和预览
  save: () => saveNote(),
}

// 右键菜单
const contextMenuItems = [
  { label: '撤销', action: () => document.execCommand('undo') },
  { label: '重做', action: () => document.execCommand('redo') },
  { label: '剪切', action: () => document.execCommand('cut') },
  { label: '复制', action: () => document.execCommand('copy') },
  { label: '粘贴', action: () => document.execCommand('paste') },
]
</script>

<template>
  <div class="note-container" @click="closeContextMenu">
    <!-- 工具栏 -->
    <div class="toolbar">
      <!-- 文本格式 -->
      <div class="toolbar-group">
        <el-tooltip content="粗体 (Ctrl+B)" placement="bottom">
          <el-button @click="toolbarActions.bold" size="small">
            <el-icon><Edit /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="斜体 (Ctrl+I)" placement="bottom">
          <el-button @click="toolbarActions.italic" size="small">
            <el-icon><Reading /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="删除线" placement="bottom">
          <el-button @click="toolbarActions.strikethrough" size="small">
            <el-icon><Minus /></el-icon>
          </el-button>
        </el-tooltip>
      </div>

      <!-- 标题 -->
      <div class="toolbar-group">
        <el-dropdown trigger="click" @command="handleCommand">
          <el-button size="small">
            标题 <el-icon><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="h1">H1</el-dropdown-item>
              <el-dropdown-item command="h2">H2</el-dropdown-item>
              <el-dropdown-item command="h3">H3</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>

      <!-- 列表 -->
      <div class="toolbar-group">
        <el-tooltip content="无序列表" placement="bottom">
          <el-button @click="toolbarActions.bulletList" size="small">
            <el-icon><List /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="有序列表" placement="bottom">
          <el-button @click="toolbarActions.numberList" size="small">
            <el-icon><Tickets /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="任务列表" placement="bottom">
          <el-button @click="toolbarActions.taskList" size="small">
            <el-icon><Check /></el-icon>
          </el-button>
        </el-tooltip>
      </div>

      <!-- 代码 -->
      <div class="toolbar-group">
        <el-tooltip content="行内代码" placement="bottom">
          <el-button @click="toolbarActions.inlineCode" size="small">
            <el-icon><Document /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="代码块" placement="bottom">
          <el-button @click="toolbarActions.codeBlock" size="small">
            <el-icon><Monitor /></el-icon>
          </el-button>
        </el-tooltip>
      </div>

      <!-- 其他工具 -->
      <div class="toolbar-group">
        <el-tooltip content="引用" placement="bottom">
          <el-button @click="toolbarActions.quote" size="small">
            <el-icon><Promotion /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="分割线" placement="bottom">
          <el-button @click="toolbarActions.divider" size="small">
            <el-icon><Minus /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="链接" placement="bottom">
          <el-button @click="toolbarActions.link" size="small">
            <el-icon><Link /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="图片" placement="bottom">
          <el-button @click="toolbarActions.image" size="small">
            <el-icon><Picture /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="表格" placement="bottom">
          <el-button @click="toolbarActions.table" size="small">
            <el-icon><Grid /></el-icon>
          </el-button>
        </el-tooltip>
      </div>

      <!-- 替换保存按钮为保存状态指示器 -->
      <div class="toolbar-group right">
        <span :class="['save-status', isSaved ? 'saved' : 'saving']">
          {{ isSaved ? '已保存' : '编辑中...' }}
        </span>
      </div>
    </div>

    <!-- 修改编辑器容器 -->
    <div class="editor-container split-view">
      <textarea
        ref="editorRef"
        v-model="editorContent"
        class="editor"
        @input="autoSave"
        @keyup="updateCursorPosition"
        @click="updateCursorPosition"
      ></textarea>
      <div
        class="preview markdown-body"
        v-html="renderedContent"
      ></div>
    </div>

    <!-- 状态栏 -->
    <div class="status-bar">
      <span>行 {{ currentLine }}, 列 {{ currentColumn }}</span>
      <span class="save-status">{{ isSaved ? '已保存' : '编辑中...' }}</span>
    </div>

    <!-- 右键菜单 -->
    <div
      v-show="showContextMenu"
      class="context-menu"
      :style="{
        left: contextMenuPosition.x + 'px',
        top: contextMenuPosition.y + 'px'
      }"
    >
      <div
        v-for="item in contextMenuItems"
        :key="item.label"
        class="context-menu-item"
        @click.stop="handleMenuItemClick(item.action)"
      >
        {{ item.label }}
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 容器样式 */
.note-container {
  display: flex;
  flex-direction: column;
  height: 100%; /* 改为100%，适应父容器高度 */
  width: 100%; /* 改为100%，适应父容器宽度 */
  margin: 0;
  padding: 0;
  overflow: hidden;
  background: white;
}

/* 工具栏样式 */
.toolbar {
  padding: 8px;
  background: #f5f5f5;
  border-bottom: 1px solid #ddd;
  display: flex;
  gap: 8px;
  align-items: center;
}

.toolbar-group {
  display: flex;
  gap: 4px;
  align-items: center;
  padding: 0 4px;
  border-right: 1px solid #ddd;
}

.toolbar-group:last-child {
  border-right: none;
}

.toolbar-group.right {
  margin-left: auto;
  border-right: none;
  padding-right: 12px;
}

:deep(.el-button) {
  padding: 6px 8px;
  border-radius: 4px;
}

:deep(.el-button:hover) {
  background-color: #e6e6e6;
}

:deep(.el-dropdown) {
  margin: 0 4px;
}

/* 编辑器容器 */
.editor-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  min-height: 0;
}

/* 分屏视图 */
.split-view {
  display: flex;
  flex: 1;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* 编辑器和预览区域共同样式 */
.editor,
.preview {
  flex: 1;
  width: 50%;
  height: 100%;
  overflow-y: auto;
  padding: 20px;
  margin: 0;
  box-sizing: border-box;
}

/* 编辑器特有样式 */
.editor {
  width: 100%;
  height: 100%;
  padding: 16px;
  border: none;
  resize: none;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
  background: #fff;
}

.editor:focus {
  outline: none;
}

/* 预览区域特有样式 */
.preview {
  border-left: 1px solid #ddd;
  background: white;
}

/* 状态栏样式 */
.status-bar {
  padding: 5px 16px;
  background: #f5f5f5;
  border-top: 1px solid #ddd;
  font-size: 12px;
  color: #666;
  display: flex;
  justify-content: space-between;
  flex-shrink: 0; /* 防止状态栏被压缩 */
}

/* 工具栏按钮样式 */
.toolbar button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.toolbar button:hover {
  background: #f0f0f0;
}

.save-btn {
  margin-left: auto;
  background: #4CAF50 !important;
  color: white;
  border: none !important;
}

.save-btn.unsaved {
  background: #ff9800 !important;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  z-index: 1000;
  min-width: 150px;
}

.context-menu-item {
  padding: 8px 12px;
  cursor: pointer;
}

.context-menu-item:hover {
  background: #f5f5f5;
}

/* Markdown 预览样式 */
:deep(.markdown-body) {
  font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  word-wrap: break-word;
  max-width: none; /* 移除最大宽度限制 */
  padding: 0; /* 移除内边距 */
}

:deep(.markdown-body h1) {
  padding-bottom: 0.3em;
  font-size: 2em;
  border-bottom: 1px solid #eaecef;
}

:deep(.markdown-body pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 6px;
}

:deep(.markdown-body code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27,31,35,0.05);
  border-radius: 6px;
}

/* 滚动条样式优化 */
.editor::-webkit-scrollbar,
.preview::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.editor::-webkit-scrollbar-thumb,
.preview::-webkit-scrollbar-thumb {
  background: #ccc;
  border-radius: 4px;
}

.editor::-webkit-scrollbar-track,
.preview::-webkit-scrollbar-track {
  background: #f1f1f1;
}

/* 代码块包装器 */
:deep(.code-block-wrapper) {
  margin: 0.8em 0;
  background: #2d333b;  /* 稍微浅一点的背景色 */
  border-radius: 4px;
  overflow: hidden;
}

/* 代码块头部 */
:deep(.code-block-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 12px;
  background: #22272e;  /* 深色背景 */
}

/* 语言标识 */
:deep(.code-lang) {
  font-size: 0.75em;
  color: #768390;  /* 柔和的文本颜色 */
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 复制按钮 */
:deep(.copy-button) {
  padding: 3px 8px;
  font-size: 0.75em;
  color: #768390;
  background: #2d333b;
  border: 1px solid #444c56;
  border-radius: 3px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.2s ease;
}

:deep(.copy-button:hover) {
  background: #373e47;
  color: #adbac7;
  border-color: #768390;
}

:deep(.copy-button:active) {
  transform: scale(0.96);
}

/* 代码块主体 */
:deep(.code-block) {
  margin: 0;
  padding: 12px;
  background: #2d333b;  /* 与包装器相同的背景色 */
  overflow-x: auto;
}

:deep(.code-block code.hljs) {
  background: transparent !important;
  padding: 0;
  margin: 0;
  font-family: 'Fira Code', 'SFMono-Regular', Consolas, monospace;
  font-size: 13px;
  line-height: 1.4;
  color: #adbac7;  /* 更亮的文本颜色 */
}

/* 滚动条样式 */
:deep(.code-block::-webkit-scrollbar) {
  width: 4px;
  height: 4px;
}

:deep(.code-block::-webkit-scrollbar-thumb) {
  background: #444c56;
  border-radius: 2px;
}

:deep(.code-block::-webkit-scrollbar-track) {
  background: #2d333b;
}

/* 代码高亮颜色主题 - GitHub Dark 风格 */
:deep(.hljs-keyword) { color: #f47067; }
:deep(.hljs-string) { color: #96d0ff; }
:deep(.hljs-comment) { color: #768390; font-style: italic; }
:deep(.hljs-function) { color: #dcbdfb; }
:deep(.hljs-number) { color: #6cb6ff; }
:deep(.hljs-operator) { color: #f47067; }
:deep(.hljs-class) { color: #dcbdfb; }
:deep(.hljs-builtin) { color: #6cb6ff; }
:deep(.hljs-params) { color: #adbac7; }
:deep(.hljs-variable) { color: #f47067; }
:deep(.hljs-title) { color: #dcbdfb; }
:deep(.hljs-attr) { color: #6cb6ff; }
:deep(.hljs-symbol) { color: #8ddb8c; }
:deep(.hljs-bullet) { color: #8ddb8c; }
:deep(.hljs-link) { color: #96d0ff; }
:deep(.hljs-meta) { color: #768390; }
:deep(.hljs-literal) { color: #6cb6ff; }

/* 行内代码样式 */
:deep(.markdown-body p code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: #2d333b;
  border-radius: 3px;
  font-family: 'Fira Code', monospace;
  color: #adbac7;
}

.save-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  transition: all 0.3s ease;
}

.save-status.saved {
  color: #67C23A;
  background-color: #f0f9eb;
}

.save-status.saving {
  color: #E6A23C;
  background-color: #fdf6ec;
}
</style>