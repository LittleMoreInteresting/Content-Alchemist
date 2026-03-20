<template>
  <router-view />
  <CommandPalette ref="commandPaletteRef" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import CommandPalette from '@/components/ai/CommandPalette.vue'

const commandPaletteRef = ref<InstanceType<typeof CommandPalette>>()

// 监听打开命令面板事件
function handleOpenCommandPalette() {
  commandPaletteRef.value?.open()
}

onMounted(() => {
  window.addEventListener('open-command-palette', handleOpenCommandPalette)
})

onUnmounted(() => {
  window.removeEventListener('open-command-palette', handleOpenCommandPalette)
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100vh;
  overflow: hidden;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 选中文字样式 */
::selection {
  background: rgba(22, 119, 255, 0.2);
}
</style>
