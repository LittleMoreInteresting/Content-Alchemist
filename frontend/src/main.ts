import { createApp } from 'vue';
import App from './components/App.vue';

// 导入全局样式
import './styles/global.css';

// 创建 Vue 应用
const app = createApp(App);

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue Error:', err);
  console.error('Component:', instance);
  console.error('Info:', info);
};

// 挂载应用
app.mount('#app');
