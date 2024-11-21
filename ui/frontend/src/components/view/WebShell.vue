<script setup lang="ts">
import {inject, onMounted, ref} from "vue";
import {useRoute, useRouter} from "vue-router";
import {InitShell, TestConnect} from "../../../wailsjs/go/client/ClientApp";
import { ElIcon } from 'element-plus';
import {
  Document,
  House,
  Cpu,
  FolderOpened,
    Notebook,
} from '@element-plus/icons-vue'
import {core} from "../../../wailsjs/go/models";
import SystemInfo = core.SystemInfo;
interface TabBar {
  addTab: (name: string, path: string) => void;
}

const tabBar = inject<{ addTab: (name: string, path: string) => void } | null>("tabBar", null);

const route = useRoute(); // 获取当前路由对象
const router =useRouter();
var systemInfo = ref<SystemInfo>()
var infoRef =ref("")
const id =Number(route.params.id)

// 在组件挂载时调用 addTab 方法
onMounted(async () => {
  var flag= await TestConnect(id)
  if (!flag){
    alert("连接失败")
    router.push('/')
    }
  if (tabBar && typeof tabBar.addTab === 'function') {
    tabBar.addTab("shell-"+id,route.path ); // 调用 addTab 方法
  } else {
    console.error('tabBar or addTab is not available');
  }
  systemInfo.value= await InitShell(Number(route.params.id))
  infoRef.value=systemInfo.value.currentDir
});

// 跳转到终端页面并传递SystemInfo对象
const goToTerminal = () => {
  router.push({ name: 'terminal', params: { id: id }, query: { systemInfo: JSON.stringify(systemInfo.value) } });
};


const goToFileManger = () => {
  router.push({ name: 'files', params: { id: id }, query: { systemInfo: JSON.stringify(systemInfo.value) } });
};

const goToHome = () => {
  router.push({ name: 'home', params: { id: id }, query: { systemInfo: JSON.stringify(systemInfo.value) } });
};

const goToMonitor = () => {
  router.push({ name: 'monitor', params: { id: id }, query: { systemInfo: JSON.stringify(systemInfo.value) } });
};

const goToNote = () => {
  router.push({ name: 'note', params: { id: id }, query: { systemInfo: JSON.stringify(systemInfo.value) } });
};


</script>

<template>
  <div class="webshell-container">
    <el-container class="webshell-content">
      <el-aside width="150px">
        <el-menu>
          <el-menu-item index="1" @click="goToHome">
            <template #title>
              <el-icon><House /></el-icon>
              <span>系统信息</span>
            </template>
          </el-menu-item>
          <el-menu-item index="2" @click="goToFileManger">
            <el-icon><FolderOpened /></el-icon>
            <template #title>文件管理</template>
          </el-menu-item>
          <el-menu-item index="3" @click="goToTerminal">
            <el-icon><document /></el-icon>
            <template #title>虚拟终端</template>
          </el-menu-item>
          <el-menu-item index="4" @click="goToMonitor">
            <el-icon><Cpu /></el-icon>
            <template #title>系统监测</template>
          </el-menu-item>
          <el-menu-item index="5" @click="goToNote">
            <el-icon><Notebook /></el-icon>
            <template #title>笔记</template>
          </el-menu-item>
        </el-menu>

      </el-aside>
      <el-container class="main-container">
        <router-view></router-view>
      </el-container>
    </el-container>
  </div>

</template>

<style scoped>
.webshell-container {
  height: 100%;
  width: 100%;
  overflow: hidden;
  position: relative;
}

.webshell-content {
  height: 100%;
  width: 100%;
}

.main-container {
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.el-menu {
  height: 100%;
  border-right: solid 1px #e6e6e6;
}
</style>