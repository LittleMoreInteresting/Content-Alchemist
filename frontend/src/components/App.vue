<template>
  <div class="app">
    <h1>Content Alchemist</h1>
    <p v-if="loading">正在初始化...</p>
    <div v-else>
      <p>应用已加载成功！</p>
      <button @click="testDialog">测试打开文件对话框</button>
      <p v-if="message">{{ message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as App from '../../wailsjs/go/backend/App';

const loading = ref(true);
const message = ref('');

onMounted(async () => {
  try {
    // 测试后端连接
    await App.GetRecentArticles(5);
    loading.value = false;
  } catch (err) {
    message.value = '后端连接失败: ' + String(err);
    loading.value = false;
  }
});

const testDialog = async () => {
  try {
    message.value = '正在打开对话框...';
    const path = await App.OpenFileDialog();
    message.value = path ? `选中文件: ${path}` : '用户取消';
  } catch (err) {
    message.value = '错误: ' + String(err);
  }
};
</script>

<style>
.app {
  padding: 40px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  text-align: center;
}
button {
  padding: 12px 24px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 20px;
}
</style>
