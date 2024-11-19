<template>
  <div class="tab-bar">
    <div
        v-for="(tab, index) in tabs"
        :key="tab.id"
        class="tab"
        :class="{ active: tab.id === activeTab }"
    >
      <span @click="switchTab(tab.id)">{{ tab.name }}</span>
<!--      <button class="close-btn" @click="closeTab(tab.id)">x</button>-->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

// 使用 ref 包裹标签数组
const tabs = ref([
  { id: "home", name: "首页", route: "/" }, // 默认第一个标签
]);

const activeTab = ref("home");
let tabCounter = 1; // 确保标签 ID 唯一

// 切换标签
const switchTab = (id: string) => {
  const tab = tabs.value.find((t) => t.id === id);
  if (tab) {
    activeTab.value = id;
    router.push(tab.route).catch((err) => {
      if (err.name !== "NavigationDuplicated") {
        console.error(err);
      }
    });
  }
};

// 添加新标签
const addTab = (name: string, route: string) => {
  const existingTab = tabs.value.find((t) => t.route === route);
  if (existingTab) {
    switchTab(existingTab.id);
    return;
  }

  const id = `tab-${tabCounter++}`;
  tabs.value.push({ id, name, route });
  switchTab(id); // 自动切换到新标签
};

// 关闭标签
const closeTab = (id: string) => {
  const index = tabs.value.findIndex((t) => t.id === id);
  if (index !== -1) {
    tabs.value.splice(index, 1);
    if (activeTab.value === id) {
      const nextTab = tabs.value[index] || tabs.value[index - 1];
      switchTab(nextTab?.id || "home");
    }
  }
};

// 暴露方法供外部调用
defineExpose({
  addTab,
});
</script>

<style scoped>
.tab-bar {
  display: flex;
  background-color: #f8f8f8;
  border-bottom: 1px solid #ddd;
  padding: 5px;
}

.tab {
  display: flex;
  align-items: center;
  padding: 5px 10px;
  margin-right: 5px;
  border: 1px solid #ddd;
  border-radius: 3px;
  background-color: white;
  cursor: pointer;
}

.tab.active {
  background-color: #007bff;
  color: white;
}

.close-btn {
  margin-left: 5px;
  border: none;
  background: none;
  cursor: pointer;
  color: red;
}
</style>
