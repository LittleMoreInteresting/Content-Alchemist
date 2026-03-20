import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Config } from '@/types'

export const useConfigStore = defineStore('config', () => {
  // State
  const config = ref<Config>({
    id: '',
    apiBaseUrl: 'https://api.deepseek.com',
    apiKey: '',
    model: 'deepseek-chat',
    temperature: 0.7,
    styleTags: [],
    audience: '',
    persona: '我',
    createdAt: '',
    updatedAt: ''
  })
  
  const hasConfig = ref(false)
  const isLoading = ref(false)
  const error = ref('')

  // Getters
  const isLocalMode = computed(() => !config.value.apiKey)
  
  const styleDescription = computed(() => {
    const parts: string[] = []
    if (config.value.styleTags.length > 0) {
      parts.push(config.value.styleTags.join('、'))
    }
    if (config.value.audience) {
      parts.push(`面向${config.value.audience}`)
    }
    if (config.value.persona) {
      parts.push(`以${config.value.persona}的视角`)
    }
    return parts.join('，')
  })

  // Actions
  function setConfig(newConfig: Config) {
    config.value = newConfig
    hasConfig.value = true
  }

  function updateConfig(updates: Partial<Config>) {
    config.value = { ...config.value, ...updates }
  }

  function setHasConfig(value: boolean) {
    hasConfig.value = value
  }

  function setLoading(value: boolean) {
    isLoading.value = value
  }

  function setError(msg: string) {
    error.value = msg
  }

  function reset() {
    config.value = {
      id: '',
      apiBaseUrl: 'https://api.deepseek.com',
      apiKey: '',
      model: 'deepseek-chat',
      temperature: 0.7,
      styleTags: [],
      audience: '',
      persona: '我',
      createdAt: '',
      updatedAt: ''
    }
    hasConfig.value = false
    error.value = ''
  }

  return {
    config,
    hasConfig,
    isLoading,
    error,
    isLocalMode,
    styleDescription,
    setConfig,
    updateConfig,
    setHasConfig,
    setLoading,
    setError,
    reset
  }
})
