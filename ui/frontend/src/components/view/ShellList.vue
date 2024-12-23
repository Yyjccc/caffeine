<script lang="ts" setup >
import {onBeforeUnmount, onMounted, ref} from 'vue';
import {GetShellID, GetShellList} from "../../../wailsjs/go/client/ClientApp";
import {client, webshell} from "../../../wailsjs/go/models";
import {StartWebShell} from '../../../wailsjs/go/main/App'
import WebClient = webshell.WebClient;
import {useRouter} from "vue-router";
let mode =0;
const empty = ref(true); // 控制是否显示 el-empty
const shellList =ref([] as client.ShellEntry[]);
// 控制菜单位置的状态
const showContextMenu = ref(false);
const menuX = ref(0);
const menuY = ref(0);
const router = useRouter();
// 处理右键点击事件
const handleRightClick = (event: MouseEvent) => {
  event.preventDefault(); // 阻止默认的右键菜单
  menuX.value = event.clientX;
  menuY.value = event.clientY;
  showContextMenu.value = true;
};
// 监听全局点击事件，点击菜单外部时关闭菜单
const handleClickOutside = (event: MouseEvent) => {
  const contextMenu = document.querySelector('.context-menu');
  if (contextMenu && !contextMenu.contains(event.target as Node)) {
    showContextMenu.value = false;
  }
};

// 处理菜单项点击事件
const handleMenuAction = async (action: string) => {
  switch (action) {
    case "1":

      console.log('进入shell');
      //启动一个shell;参数省略
      var shellID = await GetShellID()
      router.push('/webshell/'+shellID).then(() => {
        console.log('跳转成功');
      }).catch((err) => {
            console.error('跳转失败:', err);
          });
  }
  showContextMenu.value = false;
};



onMounted(() => {

  GetShellList(mode).then((res)=>{
    if(res.length!=0){
      empty.value=false
     shellList.value=res
    }

  })
  // 添加全局点击事件监听
  document.addEventListener('mousedown', handleClickOutside);
});
onBeforeUnmount(() => {
  // 组件销毁时移除全局点击事件监听
  document.removeEventListener('mousedown', handleClickOutside);
});

</script>

<template>

  <el-empty v-if="empty" description="No data available" />
  <el-table v-else :data="shellList" stripe style="width: 100%" @contextmenu="handleRightClick">
    <el-table-column prop="ID" label="ID"  width="60px"/>
    <el-table-column prop="ShellType" label="类型"  />
    <el-table-column  prop="URL" label="URL" width="180" />
    <el-table-column  prop="Location" label="位置" />
    <el-table-column  prop="Note" label="备注" />
    <el-table-column  prop="CreateTime" label="创建时间" />
    <el-table-column  prop="UpdateTime" label="更新时间" />
  </el-table>


  <!-- 右键菜单 -->
  <div
      v-if="showContextMenu"
      :style="{ top: `${menuY}px`, left: `${menuX}px` }"
      class="context-menu"
      @click="showContextMenu = false">
  <ul class="right_menu">
    <li @click="handleMenuAction('1')" >进入shell</li>
    <li @click="handleMenuAction('action2')">编辑数据</li>
    <li @click="handleMenuAction('action2')">虚拟终端</li>
    <li @click="handleMenuAction('action2')">复制URL</li>
    <li @click="handleMenuAction('action2')">进入缓存</li>
    <li @click="handleMenuAction('action2')">删除缓存</li>
  </ul>
  </div>
</template>

<style scoped>
.context-menu {
  position: absolute;
  background-color: white;
  border: 1px solid #ddd;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  padding: 4px;
  font-size: 12px;
  border-radius: 4px;
}
.context-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}
.context-menu li {
  padding: 8px 12px;
  cursor: pointer;
}
.context-menu li:hover {
  background-color: #f5f5f5;
  font-size: 14px;
}


</style>

