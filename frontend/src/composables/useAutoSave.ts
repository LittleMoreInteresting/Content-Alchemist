import { ref, onUnmounted } from 'vue'
import { debounce } from '@/utils/performance'

interface AutoSaveOptions {
  delay?: number
  onSave?: () => Promise<void>
  onError?: (error: Error) => void
}

export function useAutoSave(options: AutoSaveOptions = {}) {
  const { delay = 3000, onSave, onError } = options
  
  const saveStatus = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')
  const lastSavedAt = ref<Date | null>(null)
  const error = ref<string | null>(null)
  
  let saveTimeout: number | null = null
  
  const save = debounce(async () => {
    if (!onSave) return
    
    saveStatus.value = 'saving'
    error.value = null
    
    try {
      await onSave()
      saveStatus.value = 'saved'
      lastSavedAt.value = new Date()
    } catch (err) {
      saveStatus.value = 'error'
      error.value = err instanceof Error ? err.message : '保存失败'
      onError?.(err as Error)
    }
  }, delay)
  
  function triggerSave() {
    saveStatus.value = 'idle'
    save()
  }
  
  function cancelSave() {
    if (saveTimeout) {
      clearTimeout(saveTimeout)
      saveTimeout = null
    }
  }
  
  onUnmounted(() => {
    cancelSave()
  })
  
  return {
    saveStatus,
    lastSavedAt,
    error,
    triggerSave,
    cancelSave
  }
}
