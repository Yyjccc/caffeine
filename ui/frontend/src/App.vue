<script lang="ts" setup>


import AsideMenu from "./components/aside/AsideMenu.vue";</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header>Header</el-header>
      <el-container>
        <el-aside width="200px">
        <AsideMenu/>
        </el-aside>
        <el-main class="main">

          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style>
.main{
  margin: 0;
  padding: 0;
}
</style>
