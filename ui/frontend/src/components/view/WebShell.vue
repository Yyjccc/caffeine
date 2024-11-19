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


</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="150px">
        <el-menu>
          <el-menu-item index="1">
            <template #title>
              <el-icon><House /></el-icon>
              <span>系统信息</span>
            </template>
          </el-menu-item>
          <el-menu-item index="2">
            <el-icon><FolderOpened /></el-icon>
            <template #title>文件管理</template>
          </el-menu-item>
          <el-menu-item index="3" @click="goToTerminal">
            <el-icon><document /></el-icon>
            <template #title>虚拟终端</template>
          </el-menu-item>
          <el-menu-item index="4">
            <el-icon><Cpu /></el-icon>
            <template #title>系统监测</template>
          </el-menu-item>
          <el-menu-item index="4">
            <el-icon><Notebook /></el-icon>
            <template #title>笔记</template>
          </el-menu-item>
        </el-menu>

      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>

</template>

<style scoped>

</style>