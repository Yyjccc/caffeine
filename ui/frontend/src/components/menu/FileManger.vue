<template>
  <div class="file-manager">
    <aside class="directory-tree" ref="directoryTree">
      <!-- 添加目录树标题 -->
      <div class="section-header">
        <el-icon><Folder /></el-icon>
        <span>目录列表</span>
        <el-icon class="collapse-icon"><ArrowLeft /></el-icon>
      </div>
      <div class="tree-content-wrapper">
        <div class="tree-content" ref="treeContent">
      <el-tree
          :data="directoryData"
          :props="treeProps"
          accordion
          :default-checked-keys="[]"
          :default-expanded-keys="[]"
          node-key="path"
          :expand-on-click-node="false"
          @node-click="handleDirectorySelect"
          @update:expanded-keys="onExpandChange"
      >
        <!-- 自定义树节点图标 -->
        <template #default="{ node, data }">
          <span class="custom-tree-node">
            <el-icon>
              <component :is="getNodeIcon(data.type)" />
            </el-icon>
            <span>{{ node.label }}</span>
          </span>
        </template>
      </el-tree>
        </div>
      </div>
      <!-- 添加水平滚动条 -->
      <div class="horizontal-scrollbar">
        <div
            class="scrollbar-thumb"
            ref="scrollThumb"
            @mousedown="startScrolling"
            :style="{
            width: thumbWidth + 'px',
            left: thumbPosition + 'px'
          }"
        ></div>
      </div>
    </aside>

    <!-- 拖动条 -->
    <div class="resize-bar" @mousedown="startResizing"></div>

    <main class="file-list">
      <!-- 添加文件列表标题 -->
      <div class="section-header">
        <el-icon><Document /></el-icon>
        <span>文件列表 ({{ currentFiles.length }})</span>
      </div>
      <div class="toolbar">
        <div class="toolbar-buttons">
          <div class="toolbar-item" @click="handleNew">
            <el-icon><Plus /></el-icon>
            <span>新建</span>
          </div>
          <div class="divider">|</div>

          <div class="toolbar-item" @click="handleParent">
            <el-icon><ArrowUp /></el-icon>
            <span>上层</span>
          </div>
          <div class="divider">|</div>

          <div class="toolbar-item" @click="handleRefresh">
            <el-icon><Refresh /></el-icon>
            <span>刷新</span>
          </div>
          <div class="divider">|</div>

          <div class="toolbar-item" @click="handleHome">
            <el-icon><House /></el-icon>
            <span>主目录</span>
          </div>
          <div class="divider">|</div>

          <div class="toolbar-item" @click="handleBookmark">
            <el-icon><Star /></el-icon>
            <span>书签</span>
          </div>
        </div>
        <div class="path-navigator">
          <el-input
              v-model="currentPath"
              placeholder="输入路径"
              class="path-input"
          >
            <template #append>
              <el-button @click="handleNavigate">
                <el-icon><Right /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>
      <table @contextmenu.prevent="showContextMenu">
        <thead>
        <tr>
          <th>名称</th>
          <th>大小</th>
          <th>日期</th>
          <th>权限</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(file, index) in currentFiles"
            :key="file.name"
            :class="{ 'row-striped': index % 2 === 1, 'selected': selectedFiles.includes(file) }"
            @click="selectFile(file, $event)"
            @contextmenu.prevent="showContextMenu($event, file)">
          <td>
            <div class="file-item">
              <el-icon><component :is="getFileIcon(file.type, file.name)" /></el-icon>
              <span>{{ file.name }}</span>
            </div>
          </td>
          <td>{{ file.size }}</td>
          <td>{{ file.date }}</td>
          <td>{{ file.permissions }}</td>
        </tr>
        </tbody>
      </table>

      <!-- 右键菜单 -->
      <div v-show="showMenu"
           class="context-menu"
           :style="{ top: menuPosition.y + 'px', left: menuPosition.x + 'px' }">
        <div class="menu-item" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          <span>刷新目录</span>
        </div>
        <div class="menu-divider"></div>
        <div class="menu-item" :class="{ disabled: !canManageFiles }" @click="handleUpload">
          <el-icon><Upload /></el-icon>
          <span>上传文件</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleDownload">
          <el-icon><Download /></el-icon>
          <span>下载文件</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleCopy">
          <el-icon><CopyDocument /></el-icon>
          <span>复制</span>
        </div>
        <div class="menu-item" :class="{ disabled: !canPaste }" @click="handlePaste">
          <el-icon><DocumentAdd /></el-icon>
          <span>粘贴</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleDelete">
          <el-icon><Delete /></el-icon>
          <span>删除</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleRename">
          <el-icon><Edit /></el-icon>
          <span>重命名</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleChangeTime">
          <el-icon><Timer /></el-icon>
          <span>更改文件时间</span>
        </div>
        <div class="menu-item" :class="{ disabled: !hasSelection }" @click="handleChangePermissions">
          <el-icon><Lock /></el-icon>
          <span>更改权限</span>
        </div>
        <div class="menu-divider"></div>
        <div class="menu-item" @click="handleNew">
          <el-icon><Plus /></el-icon>
          <span>新建</span>
        </div>
        <div class="menu-item" @click="handleOpenTerminal">
          <el-icon><Operation /></el-icon>
          <span>在此处打开终端</span>
        </div>
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, onUnmounted, ref} from "vue";

import {
  Folder,
  Document,
  Plus,
  House,
  ArrowUp,
  ArrowLeft,
  Refresh,
  Star,
  Right,
  Picture,
  VideoCamera,
  Headset,
  Files,
  Platform,
  Upload,
  Download,
  CopyDocument,
  DocumentAdd,
  Delete,
  Edit,
  Timer,
  Lock,
  Operation,
} from '@element-plus/icons-vue'

const treeContent = ref<HTMLElement | null>(null);
const scrollThumb = ref<HTMLElement | null>(null);
let isScrolling = false;
let startScrollX = 0;
let startThumbPosition = 0;
let thumbWidth = ref(50); // 初始滑块宽度
let thumbPosition = ref(0);
const currentPath = ref('C:/');

// 右键菜单状态
const showMenu = ref(false);
const menuPosition = ref({ x: 0, y: 0 });
const selectedFiles = ref<any[]>([]);
const clipboard = ref<any[]>([]);

// 计算属性
const hasSelection = computed(() => selectedFiles.value.length > 0);
const canManageFiles = computed(() => true); // 可以基于权限判断
const canPaste = computed(() => clipboard.value.length > 0);

// 添加工具栏按钮处理函数
const handleNew = () => {
  console.log('新建');
};

const handleParent = () => {
  console.log('跳转到上层目录');
};

const handleRefresh = () => {
  console.log('刷新当前目录');
};

const handleHome = () => {
  console.log('跳转到主目录');
};

const handleBookmark = () => {
  console.log('打开书签');
};

const handleNavigate = () => {
  console.log('导航到:', currentPath.value);
};

// 文件选择
const selectFile = (file: any, event: MouseEvent) => {
  if (event.ctrlKey) {
    const index = selectedFiles.value.indexOf(file);
    if (index === -1) {
      selectedFiles.value.push(file);
    } else {
      selectedFiles.value.splice(index, 1);
    }
  } else {
    selectedFiles.value = [file];
  }
};

// 显示右键菜单
const showContextMenu = (event: MouseEvent, file?: any) => {
  event.preventDefault();
  if (file && !selectedFiles.value.includes(file)) {
    selectedFiles.value = [file];
  }
  menuPosition.value = { x: event.clientX, y: event.clientY };
  showMenu.value = true;
};

// 隐藏右键菜单
const hideContextMenu = () => {
  showMenu.value = false;
};

// 添加全局点击事件来关闭菜单
onMounted(() => {
  document.addEventListener('click', hideContextMenu);
});

onUnmounted(() => {
  document.removeEventListener('click', hideContextMenu);
});

// 菜单项处理函数
const handleUpload = () => {
  if (!canManageFiles.value) return;
  console.log('上传文件');
};

const handleDownload = () => {
  if (!hasSelection.value) return;
  console.log('下载文件', selectedFiles.value);
};

const handleCopy = () => {
  if (!hasSelection.value) return;
  clipboard.value = [...selectedFiles.value];
  console.log('复制文件', selectedFiles.value);
};

const handlePaste = () => {
  if (!canPaste.value) return;
  console.log('粘贴文件', clipboard.value);
};

const handleDelete = () => {
  if (!hasSelection.value) return;
  console.log('删除文件', selectedFiles.value);
};

const handleRename = () => {
  if (!hasSelection.value) return;
  console.log('重命名文件', selectedFiles.value[0]);
};

const handleChangeTime = () => {
  if (!hasSelection.value) return;
  console.log('更改文件时间', selectedFiles.value);
};

const handleChangePermissions = () => {
  if (!hasSelection.value) return;
  console.log('更改权限', selectedFiles.value);
};

const handleOpenTerminal = () => {
  console.log('打开终端');
};
// 获取树节点图标
const getNodeIcon = (type: string) => {
  switch (type) {
    case 'disk':
      return Platform
    case 'folder':
      return Folder
    case 'file':
      return Document
    default:
      return Document
  }
}

// 获取文件图标
const getFileIcon = (type: string, fileName: string) => {
  if (type === 'folder') return Folder

  const extension = fileName.split('.').pop()?.toLowerCase()
  switch (extension) {
    case 'jpg':
    case 'png':
    case 'gif':
      return Picture
    case 'mp4':
    case 'avi':
    case 'mkv':
      return VideoCamera
    case 'mp3':
    case 'wav':
      return Headset
    case 'zip':
    case 'rar':
    case '7z':
      return Files
    default:
      return Document
  }
}

// 计算滑块宽度和位置
const updateScrollbar = () => {
  if (!treeContent.value) return;

  const container = treeContent.value.parentElement;
  if (!container) return;  // 添加空检查

  const contentWidth = treeContent.value.scrollWidth;
  const containerWidth = container.clientWidth;

  // 计算滑块宽度
  thumbWidth.value = Math.max(30, (containerWidth / contentWidth) * containerWidth);

  // 更新滑块位置
  const scrollPercent = treeContent.value.scrollLeft / (contentWidth - containerWidth);
  thumbPosition.value = scrollPercent * (containerWidth - thumbWidth.value);
};

// 开始滚动
const startScrolling = (e: MouseEvent) => {
  if (!scrollThumb.value || !treeContent.value) return;

  isScrolling = true;
  startScrollX = e.clientX;
  startThumbPosition = thumbPosition.value;

  window.addEventListener('mousemove', onScrolling);
  window.addEventListener('mouseup', stopScrolling);
};

// 滚动过程
const onScrolling = (e: MouseEvent) => {
  if (!isScrolling || !treeContent.value) return;

  const container = treeContent.value.parentElement;
  if (!container) return;  // 添加空检查

  const contentWidth = treeContent.value.scrollWidth;
  const containerWidth = container.clientWidth;

  const dx = e.clientX - startScrollX;
  const maxThumbPosition = containerWidth - thumbWidth.value;

  // 更新滑块位置
  thumbPosition.value = Math.max(0, Math.min(maxThumbPosition, startThumbPosition + dx));

  // 更新内容滚动位置
  const scrollPercent = thumbPosition.value / maxThumbPosition;
  treeContent.value.scrollLeft = scrollPercent * (contentWidth - containerWidth);
};

// 停止滚动
const stopScrolling = () => {
  isScrolling = false;
  window.removeEventListener('mousemove', onScrolling);
  window.removeEventListener('mouseup', stopScrolling);
};

// 监听内容变化
onMounted(() => {
  document.addEventListener('click', hideContextMenu);
  if (treeContent.value) {
    const observer = new ResizeObserver(() => {
      updateScrollbar();
    });
    observer.observe(treeContent.value);
  }
});
// 定义节点类型接口
interface DirectoryNode {
  name: string;
  type: "disk" | "folder" | "file";
  path: string; // 文件或目录路径
  children?: DirectoryNode[]; // 子节点
}

// 初始化目录树数据
const directoryData = ref<DirectoryNode[]>([
  {
    name: "C:/",
    type: "disk",
    path: "C:/",
    children: [
      {
        name: "Android-Tool",
        type: "folder",
        path: "C:/Android-Tool",
        children: [
          { name: "AndroidKiller4J-Windows", type: "folder", path: "C:/Android-Tool/AndroidKiller4J-Windows" },
          { name: "android-killer-main", type: "folder", path: "C:/Android-Tool/android-killer-main" },
          { name: "fridaUiTools_for_window", type: "folder", path: "C:/Android-Tool/fridaUiTools_for_window" },
        ],
      },
      { name: "Docker", type: "folder", path: "C:/Docker" },
    ],
  },
  { name: "D:/", type: "disk", path: "D:/", children: [] },
]);

// 定义文件类型接口
interface FileItem {
  name: string;
  size: string;
  date: string;
  permissions: string;
  type: 'file' | 'folder';
}

// 当前显示的文件列表
const currentFiles = ref<FileItem[]>([
  {
    name: "AndroidKiller4J-Windows",
    size: "0 b",
    date: "2024-05-14 07:59:20",
    permissions: "0777",
    type: "folder",
  },
  {
    name: "android-killer-main",
    size: "128 Kb",
    date: "2024-05-14 08:00:15",
    permissions: "0777",
    type: "folder",
  },
  {
    name: "fridaUiTools_for_window.zip",
    size: "440.83 Mb",
    date: "2024-05-14 08:23:12",
    permissions: "0666",
    type: "file",
  },
  {
    name: "example.txt",
    size: "2.5 Kb",
    date: "2024-05-14 09:15:00",
    permissions: "0644",
    type: "file",
  }
]);
// 定义树形组件的属性配置
const treeProps = {
  children: 'children',  // 指定子节点字段
  label: 'name',         // 指定显示名称字段
  disabled: 'disabled',  // 指定禁用状态字段
};

// 处理目录点击事件
const handleDirectorySelect = (node: DirectoryNode) => {
  console.log("选择的目录路径:", node.path);

  // 根据选择的目录更新文件列表
  if (node.path === "C:/Android-Tool") {
    currentFiles.value = [
      {
        name: "AndroidKiller4J-Windows",
        size: "0 b",
        date: "2024-05-14 07:59:20",
        permissions: "0777",
        type: "folder",
      },
      {
        name: "android-killer-main",
        size: "128 Kb",
        date: "2024-05-14 08:00:15",
        permissions: "0777",
        type: "folder",
      },
      {
        name: "fridaUiTools_for_window.zip",
        size: "440.83 Mb",
        date: "2024-05-14 08:23:12",
        permissions: "0666",
        type: "file",
      }
    ];
  } else {
    // 其他目录的默认文件列表
    currentFiles.value = [
      {
        name: "example.txt",
        size: "2.5 Kb",
        date: "2024-05-14 09:15:00",
        permissions: "0644",
        type: "file",
      },
      {
        name: "example-folder",
        size: "--",
        date: "2024-05-14 09:30:00",
        permissions: "0755",
        type: "folder",
      }
    ];
  }
};


// 扩展/收起时更新树的状态
const onExpandChange = (keys: string[]) => {
  console.log("当前展开的节点路径:", keys);
};

// 拖动条的拖动逻辑
let isResizing = false;
let startX = 0;
let startWidth = 0;

const startResizing = (e: MouseEvent) => {
  isResizing = true;
  startX = e.clientX;
  startWidth = (document.querySelector(".directory-tree") as HTMLElement).offsetWidth;
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseup", stopResizing);
};

const onMouseMove = (e: MouseEvent) => {
  if (!isResizing) return;
  const dx = e.clientX - startX;
  let newWidth = startWidth + dx;

  // 限制宽度的最小值和最大值
  const minWidth = window.innerWidth * 0.05;  // 最小 10%
  const maxWidth = window.innerWidth * 0.6;  // 最大 90%

  // 限制新宽度在最小和最大范围内
  newWidth = Math.max(minWidth, Math.min(newWidth, maxWidth));

  const directoryTree = document.querySelector(".directory-tree") as HTMLElement;
  directoryTree.style.width = `${newWidth}px`;
};

const stopResizing = () => {
  isResizing = false;
  window.removeEventListener("mousemove", onMouseMove);
  window.removeEventListener("mouseup", stopResizing);
};
</script>

<style scoped>
.file-manager {
  display: flex;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.directory-tree {
  display: flex;
  flex-direction: column;
  width: 200px;
  min-width: 100px;
  max-width: 50%;
  border-right: 1px solid #ddd;
  padding: 10px;
  height: 100%;
  overflow: hidden;
}

.resize-bar {
  cursor: ew-resize;
  width: 5px;
  background-color: #ccc;
  height: 100%;
  flex-shrink: 0;
}

.file-list {
  flex: 1;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 文件列表标题和工具栏 */
.section-header {
  flex-shrink: 0;
}

.toolbar {
  flex-shrink: 0;
}

/* 文件列表表格容器 */
.file-list-container {
  flex: 1;
  overflow: auto;
  padding: 10px;
}

/* 表格样式 */
table {
  width: 100%;
  border-collapse: collapse;
}

/* 确保表格列合理分配宽度 */
th:nth-child(1),
td:nth-child(1) {
  width: 40%;  /* 文件名列 */
}

th:nth-child(2),
td:nth-child(2) {
  width: 15%;  /* 大小列 */
}

th:nth-child(3),
td:nth-child(3) {
  width: 25%;  /* 日期列 */
}

th:nth-child(4),
td:nth-child(4) {
  width: 20%;  /* 权限列 */
}

/* 水平滚动条 */
.horizontal-scrollbar {
  height: 8px;
  background-color: #f0f0f0;
  position: relative;
  border-radius: 4px;
  margin: 4px 0;
  min-height: 8px;
}

.scrollbar-thumb {
  position: absolute;
  height: 100%;
  background-color: #ccc;
  border-radius: 4px;
  cursor: pointer;
  user-select: none;
  min-width: 20px;
  transition: background-color 0.2s;
}

.scrollbar-thumb:hover {
  background-color: #999;
}

.scrollbar-thumb:active {
  background-color: #666;
}

/* 分隔条 */
.resize-bar {
  cursor: ew-resize;
  width: 5px;
  background-color: #ccc;
  height: 100%;
}

/* 标题栏 */
.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background-color: rgb(245, 245, 245);
  border-bottom: 1px solid #ddd;
  font-family: SimSun, "宋体", serif;
  font-size: 14px;
  font-weight: bold;
}

.section-header .collapse-icon {
  margin-left: auto;
  cursor: pointer;
}

/* 工具栏 */
.toolbar {
  display: flex;
  padding: 8px 14px;
  gap: 14px;
  border-bottom: 1px solid #ddd;
  align-items: center;
  background-color: rgb(235, 235, 235);
}

.toolbar-buttons {
  display: flex;
  align-items: center;
  gap: 14px;
}

.toolbar-item {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 4px 8px;
  font-family: SimSun, "宋体", serif;
  font-size: 14px;
}

.toolbar-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.divider {
  color: #999;
  font-size: 14px;
  font-weight: lighter;
}

/* 路径导航 */
.path-navigator {
  flex: 1;
  margin-left: 14px;
}

.path-input {
  width: 100%;
  font-family: SimSun, "宋体", serif;
  font-size: 14px;
}

.path-input :deep(.el-input__inner) {
  font-size: 14px;
  font-family: SimSun, "宋体", serif;
}

/* 文件列表 */
.file-list {
  flex: 1;
  padding: 10px;
  overflow-x: auto;
}

.file-list table {
  width: 100%;
  border-collapse: collapse;
}

.file-list th,
.file-list td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.file-list th {
  background-color: #f4f4f4;
  font-weight: bold;
}

.file-list tr:hover {
  background-color: #f5f7fa;
}

.file-list tr.row-striped {
  background-color: #fafafa;
}

.file-list tr.row-striped:hover {
  background-color: #f5f7fa;
}

/* 单元格对齐 */
.file-list td:nth-child(2) { text-align: right; }  /* 大小列 */
.file-list td:nth-child(3) { text-align: center; } /* 日期列 */
.file-list td:nth-child(4) { text-align: center; } /* 权限列 */

/* 图标相关 */
.el-icon {
  font-size: 14px;
}

.custom-tree-node,
.file-item {
  display: flex;
  align-items: center;
}

.custom-tree-node {
  gap: 4px;
}

.file-item {
  gap: 8px;
}

.custom-tree-node .el-icon,
.file-item .el-icon {
  color: #666;
}

/* Element Plus 树组件覆盖样式 */
:deep(.el-tree-node__content) {
  white-space: nowrap;
  min-width: fit-content;
}

:deep(.el-tree) {
  min-width: fit-content;
  overflow: visible;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 4px 0;
  min-width: 200px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
  z-index: 1000;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
}

.menu-item:hover:not(.disabled) {
  background-color: #f5f7fa;
}

.menu-item.disabled {
  color: #c0c4cc;
  cursor: not-allowed;
}

.menu-divider {
  height: 1px;
  background-color: #ddd;
  margin: 4px 0;
}

/* 选中行的样式 */
.file-list tr.selected {
  background-color: #e6f3ff;
}

.file-list tr.selected:hover {
  background-color: #e6f3ff;
}
</style>
