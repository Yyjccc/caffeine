<script setup lang="ts">

import ShellList from "./view/ShellList.vue";
import AsideMenu from "./aside/AsideMenu.vue";
</script>

<template>


  <div class="common-layout">
    <el-container>
      <el-header>Header</el-header>
      <el-container>
        <el-aside width="200px">
          <AsideMenu/>
        </el-aside>
        <el-main class="main">

          <ShellList/>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped>

</style>