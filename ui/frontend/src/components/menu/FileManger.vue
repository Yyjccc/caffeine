<template>
  <div class="file-manager">
    <aside class="directory-tree" ref="directoryTree">
      <!-- 添加目录树标题 -->
      <div class="section-header">
        <el-icon><Folder /></el-icon>
        <span>目录列表</span>
        <el-icon class="collapse-icon" @click="toggleCollapse"><ArrowLeft /></el-icon>
      </div>
      <div class="tree-content-wrapper">
        <div class="tree-content" ref="treeContent">
          <el-tree
            :data="dirTreeRoot"
            :props="treeProps"
            accordion
            :default-checked-keys="[]"
            :default-expanded-keys="[]"
            node-key="path"
            :expand-on-click-node="true"
            @node-click="handleDirectorySelect"
            @update:expanded-keys="onExpandChange"
          >
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
              @contextmenu.prevent="showContextMenu($event, file)"
              @dblclick="handleDoubleClick(file)">  <!-- 添加双击事件 -->
            <td>
              <div class="file-item">
                <el-icon><component :is="getFileIcon(file.type, file.name)" /></el-icon>
                <span>{{ file.name }}</span>
              </div>
            </td>
            <td>{{ formattedSize(file.size)  }}</td>
            <td>{{ file.date }}</td>
            <td>{{ file.permissions !== -1 ? file.permissions : ''  }}</td>
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
import {OpenSelectFilePath} from "../../../bindings/ui/app"

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

import {useRoute} from "vue-router";
import {Directory, FileSystemCache} from "../../../bindings/caffeine/core";
import {DownloadFile, GetFileSystem, LoadDirInfo} from "../../../bindings/caffeine/client/clientapp";
import store from "../../utils/store";
import {WebShellSession} from "../../utils/session";



// 定义节点类型接口
interface DirectoryNode {
  name: string;
  type: "folder" | "file";
  path: string;
  loaded: boolean;
  children: DirectoryNode[];
}

// 定义文件类型接口
interface FileItem {
  name: string;
  size: number;
  date: string;
  permissions: number;
  type: 'file' | 'folder';
}

const route = useRoute();
const id = Number(route.params.id);

const fileSystem = ref<FileSystemCache | null>(null);
const currentPath = ref('C:/'); // 默认路径

const treeContent = ref<HTMLElement | null>(null);
const scrollThumb = ref<HTMLElement | null>(null);
let isScrolling = false;
let startScrollX = 0;
let startThumbPosition = 0;
let thumbWidth = ref(50); // 初始滑块宽度
let thumbPosition = ref(0);

// 右键菜单状态
const showMenu = ref(false);
const menuPosition = ref({ x: 0, y: 0 });
const selectedFiles = ref<any[]>([]);
const clipboard = ref<any[]>([]);

// 计算属性
const hasSelection = computed(() => selectedFiles.value.length > 0);
const canManageFiles = computed(() => true); // 可以基于权限判断
const canPaste = computed(() => clipboard.value.length > 0);


const session = computed(() => store.state.webShellCache.sessions[id] || null)

// 当前显示的文件列表
const currentFiles = ref<FileItem[]>([
]);
// 定义树形组件的属性配置
const treeProps = {
  children: 'children',  // 指定子节点字段
  label: 'name',         // 指定显示名称字段
  disabled: 'disabled',  // 指定禁用状态字段
};

//加载目录节点到目录树中(构建父目录节点)
const loadDirToTree= (currentNode: DirectoryNode)=> {

  const pathParts = currentNode.path
      .replace(/\\/g, '/')  // 将 Windows 路径中的反斜杠替换为正斜杠
      .split('/');  // 然后按正斜杠分割路径
  let dirTree = dirTreeRoot.value;
  if (!dirTree) {
    console.error("目录树未初始化");
    return;
  }
  console.log("目录树",dirTree)
  let parent:DirectoryNode = {
    name:"/",
    type:"folder",
    loaded:false,
    path:"/",
    children:[]
  };
  for (let i = 0; i < pathParts.length; i++) {
    //构建根目录
    const part = pathParts[i];
    // 查找当前路径部分是否已存在
    let existingNode:DirectoryNode | undefined
    if (i===0){
      existingNode = dirTree.find(node => node.name === part && node.type === "folder" );
    } else {
       existingNode = parent.children.find(node => node.name === part && node.type === "folder");
    }
    // 如果不存在，则创建新节点
    if (!existingNode) {

      existingNode = {
        name: part,
        type: "folder",
        path: (pathParts.slice(0, i + 1).join('/')),
        loaded:false,
        children: []
      };
      if (i==0){
        //如果是根节点直接添加
        dirTree.push(existingNode);
        console.log("添加节点 ",part,parent)
      }else {
        parent.children.push(existingNode)
      }

    }
    else {
      //如果是加载目录
      if (equalPath(currentNode.path,existingNode.path) && !existingNode.loaded){
        for (const sub of currentNode.children){
            existingNode.children.push(sub)
        }
       // existingNode.loaded=true
      }
      //如果存在，加载子目录信息
    }
    parent = existingNode
  }
  // 如果当前节点是文件，添加它到目录树的最后一级
  if (currentNode.type === "file") {
    dirTree.push(currentNode);
  }
  dirTreeRoot.value=dirTree
}

const GetFileInfoFromDir= (directory: Directory):FileItem[] =>{
  var res :FileItem[] = []

  //转化目录
  for (const sub of directory.sub) {
      let dir:FileItem = {
        name: sub.name,
        permissions:-1,
        date:"",
        size: -1,
        type: "folder",
      }
      res.push(dir)
  }
  //转化文件
  for (const item of directory.files) {
    let file:FileItem ={
      name:item.name,
      date:item.lastModified.replace("T"," ").replace("Z",""),
      size:item.size,
      permissions:item.permissions,
      type:"file",
    }
    res.push(file)
  }
  return res
}

// 初始化目录树数据
const dirTreeRoot = ref<DirectoryNode[]>([]);

// 将后端的 Directory 转换为前端的 DirectoryNode
const convertToDirectoryNode = (directory: Directory): DirectoryNode => {
  let children: DirectoryNode[] = [];
  console.log(directory)
  if (directory.sub){
    directory.sub.map(subDir => {
      children.push({
        name: subDir.name,
        path: subDir.path,
        type: "folder",
        children:[],
        loaded: false,
      })
    });
  }
  return {
    name: directory.name,
    type: 'folder',
    path: directory.path,
    children: children,
    loaded:true
  };
};

// ui切换到某个路径
function GotoPath(path:string){
  var dirNode = getDirNodeFromTree(path);
  if (dirNode==null){
    console.error("error dirNode is null,path:",path)
    return
  }
  handleDirectorySelect(dirNode)
  currentPath.value=normalizePath(path)
}

//通过路径获取目录节点,递归遍历
function  getDirNodeFromTree(path:string):DirectoryNode|null{
  let pathNormal = normalizePath(path)
  const pathParts = pathNormal.split('/');  // 然后按正斜杠分割路径
  let parent = dirTreeRoot.value[0];
  let curNode:DirectoryNode | undefined =undefined;
  for (let i = 0; i < pathParts.length; i++) {
    let currentPath = pathParts.slice(0, i + 1).join('/');
    if (i===0){
      curNode = dirTreeRoot.value.find(node => equalPath(node.path,currentPath) && node.type === "folder" );
    }else {
      curNode = parent.children.find(node => equalPath(node.path,currentPath) && node.type === "folder" );

    }
    if (!curNode){
      return null;
    }
    parent = curNode;
    //构建根目录
  }
  if (!curNode){
    return null;
  }
  return curNode;
}


function getParentDirectoryPath(filePath: string): string {
  if (!filePath) return '';

  // 标准化路径，处理不同操作系统的分隔符
  const normalizedPath = filePath.replace(/\\/g, '/');

  // 移除末尾的斜杠（如果存在）
  const trimmedPath = normalizedPath.replace(/\/+$/, '');

  // 获取最后一个斜杠的位置
  const lastSlashIndex = trimmedPath.lastIndexOf('/');

  // 如果没有找到斜杠，则返回空字符串（说明已是根目录）
  if (lastSlashIndex === -1) return '';

  // 截取到最后一个斜杠之前的部分，即上级目录路径
  return trimmedPath.substring(0, lastSlashIndex);
}

// 转换文件大小
const formattedSize = (bytes: number) => {
  if (bytes === -1) {
    return ''
  }
  if (bytes < 1024) {
    return `${bytes} B`
  } else if (bytes < 1024 ** 2) {
    return `${(bytes / 1024).toFixed(2)} KB`
  } else if (bytes < 1024 ** 3) {
    return `${(bytes / 1024 ** 2).toFixed(2)} MB`
  } else {
    return `${(bytes / 1024 ** 3).toFixed(2)} GB`
  }
}

const handleDoubleClick = (file:FileItem)=>{
  if(file.type==="folder"){
    GotoPath(normalizePath(currentPath.value+"/"+file.name));
  }
}

// 添加工具栏按钮处理函数
const handleNew = () => {
  console.log('新建');
};



// 跳转到父级目录
const handleParent = () => {
  var parentDirectoryPath = getParentDirectoryPath(currentPath.value);
  if (parentDirectoryPath===''){
    return
  }
  GotoPath(parentDirectoryPath);
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


const toggleCollapse=()=>{

}





// 菜单项处理函数
const handleUpload = () => {
  if (!canManageFiles.value) return;
  console.log('上传文件');
};

const handleDownload = async () => {
  if (!hasSelection.value) return;
  const file =selectedFiles.value[0]
  var savePath = await OpenSelectFilePath(file.name);
  const targetPath=normalizePath(currentPath.value+"/"+file.name)
  if( savePath!==""){
    DownloadFile(id,targetPath,savePath)
  }
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

// 比较路径是否相等
function equalPath(path1: string, path2: string): boolean {
  // 将路径标准化并转换为小写
  const normalizedPath1 = normalizePath(path1).toLowerCase();
  const normalizedPath2 = normalizePath(path2).toLowerCase();
  return normalizedPath1 === normalizedPath2;
}

//标准化路径
function normalizePath(inputPath: string): string {
  // 替换所有反斜杠为正斜杠
  let normalizedPath = inputPath.replace(/\\/g, "/");
  // 移除多余的斜杠（例如 "C://path//to//file"）
  normalizedPath = normalizedPath.replace(/\/+/g, "/");
  // 转换为小写以支持大小写不敏感
  return normalizedPath.toLowerCase();
}

// 处理目录点击事件
const handleDirectorySelect = async (node: DirectoryNode) => {
  var directory = await LoadDirInfo(id,node.path);
  var dirNode = convertToDirectoryNode(directory);
  if(!node.loaded){
    console.log("加载到目录树：",dirNode)
    loadDirToTree(dirNode)
    node.loaded=true
  }
  currentPath.value=node.path
  currentFiles.value=GetFileInfoFromDir(directory)
};

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
  if (!treeContent.value || !fileSystem.value) return; // 确保 fileSystem 已初始化

  const container = treeContent.value.parentElement;
  if (!container) return;  // 添加空检查

  const contentWidth = treeContent.value.scrollWidth;
  const containerWidth = container.clientWidth;

  // 确保参与运算的变量是 number 类型
  const calculatedThumbWidth = Math.max(30, (containerWidth / contentWidth) * containerWidth);
  thumbWidth.value = calculatedThumbWidth;

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


// 监听内容变化
onMounted(async () => {
  document.addEventListener('click', hideContextMenu);
  if (treeContent.value) {
    const observer = new ResizeObserver(() => {
      updateScrollbar();
    });
    observer.observe(treeContent.value);
  }
  //初始化文件系统
  fileSystem.value = await GetFileSystem(id)
  //获取session
  //var session =await store.dispatch("getSession",id);
  const session:WebShellSession = await store.dispatch('getSession', id)
  if (session) {
    //如果是Windows，添加所有盘符
    if (session.SystemType ===1){
      const drivers=session.SystemInfo.fileRoots
      for (const driver of drivers){
        dirTreeRoot.value.push(
            {
              type:"folder",
              name: driver,
              path: driver,
              loaded:false,
              children:[],
            }
        )
      }
    }
  } else {
    console.warn('Session not found for id:', id)
  }

  var currentDir = await LoadDirInfo(id,".");
  loadDirToTree(convertToDirectoryNode(currentDir))
  currentFiles.value= GetFileInfoFromDir(currentDir)
  currentPath.value = currentDir.path
});

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
  overflow-y: auto;  /* 启用垂直滚动条 */
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

.search-input {
  width: 200px;
  margin-right: 10px;
}

.directory-tree {
  overflow-y: auto;  /* 启用垂直滚动条 */
}

.tree-content {
  padding: 10px; /* 添加内边距 */
}
</style>
