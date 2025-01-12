<template>
  <div class="file-transfer-manager">
    <h3>File Transfer Manager</h3>

    <!-- 上传区域 -->
    <section class="upload-section">
      <h4>Upload Files</h4>
      <input type="file" multiple @change="handleFileUpload" />
    </section>

    <!-- 当前任务 -->
    <section v-if="activeTasks.length > 0" class="active-tasks">
      <h4>Active Transfers</h4>
      <ul>
        <li v-for="(task, index) in activeTasks" :key="task.id">
          <div class="task-info">
            <span>{{ task.fileName }}</span>
            <progress v-if="task.progress !== null" :value="task.progress" max="100"></progress>
            <span class="status">{{ task.type }} - {{ task.status }}</span>
          </div>
        </li>
      </ul>
    </section>

    <!-- 历史任务记录 -->
    <section v-if="completedTasks.length > 0" class="completed-tasks">
      <h4>Transfer History</h4>
      <ul>
        <li v-for="(task, index) in completedTasks" :key="task.id">
          <div class="task-info">
            <span>{{ task.fileName }}</span>
            <span>{{ task.completedAt }}</span>
            <template v-if="task.type === 'download'">
              <a :href="task.fileUrl" download>
                Download Again
              </a>
            </template>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from "vue";
import { useStore } from "vuex";

export default defineComponent({
  name: "FileTransferManager",
  setup() {
    const store = useStore();

    // 当前任务（上传和下载）
    const activeTasks = computed(() =>
        store.state.transferTasks.filter((task) => task.status === "transferring")
    );

    // 历史任务记录
    const completedTasks = computed(() =>
        store.state.transferTasks.filter((task) => task.status === "completed")
    );

    // 处理文件上传
    const handleFileUpload = (event: Event) => {
      const files = (event.target as HTMLInputElement).files;
      if (!files) return;
      Array.from(files).forEach((file) => {
        store.dispatch("addTransferTask", {
          id: Date.now() + Math.random(),
          fileName: file.name,
          file: file,
          type: "upload",
        });
      });
    };

    return {
      activeTasks,
      completedTasks,
      handleFileUpload,
    };
  },
});
</script>

<style scoped>
.file-transfer-manager {
  max-width: 600px;
  margin: auto;
  font-family: Arial, sans-serif;
  padding: 20px;
}

h3 {
  text-align: center;
}

ul {
  list-style: none;
  padding: 0;
}

.task-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.status {
  font-size: 0.9em;
  color: gray;
}

progress {
  width: 150px;
  margin: 0 10px;
}
</style>
