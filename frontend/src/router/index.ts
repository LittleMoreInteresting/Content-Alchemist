import { createRouter, createWebHashHistory } from 'vue-router'
import { useConfigStore } from '@/stores'
import { checkHasConfig, fetchConfig } from '@/api'

const routes = [
  {
    path: '/',
    name: 'Home',
    redirect: '/loading'
  },
  {
    path: '/loading',
    name: 'Loading',
    component: () => import('@/views/LoadingView.vue'),
    meta: { public: true }
  },
  {
    path: '/welcome',
    name: 'Welcome',
    component: () => import('@/views/WelcomeView.vue'),
    meta: { public: true }
  },
  {
    path: '/editor',
    name: 'Editor',
    component: () => import('@/views/EditorView.vue')
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue')
  },
  {
    path: '/publish',
    name: 'Publish',
    component: () => import('@/views/PublishView.vue')
  },
  {
    path: '/topic',
    name: 'TopicCenter',
    component: () => import('@/views/topic/TopicCenter.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫 - 初始化配置检查
router.beforeEach(async (to, _from, next) => {
  const configStore = useConfigStore()
  
  // 如果已经检查过配置，直接放行
  if (configStore.hasConfig || to.meta.public) {
    next()
    return
  }
  
  // 检查是否有配置
  try {
    const hasConfig = await checkHasConfig()
    if (hasConfig) {
      // 加载配置
      const config = await fetchConfig()
      configStore.setConfig(config)
      next()
    } else {
      // 没有配置，重定向到欢迎页
      if (to.path !== '/welcome' && to.path !== '/loading') {
        next('/welcome')
      } else {
        next()
      }
    }
  } catch (error) {
    console.error('Failed to check config:', error)
    // 出错时也重定向到欢迎页
    if (to.path !== '/welcome' && to.path !== '/loading') {
      next('/welcome')
    } else {
      next()
    }
  }
})

export default router
