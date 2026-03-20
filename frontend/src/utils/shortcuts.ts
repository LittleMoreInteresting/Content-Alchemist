import type { App } from 'vue'

// 快捷键管理器
export class ShortcutManager {
  private shortcuts: Map<string, () => void> = new Map()
  private enabled = true

  register(key: string, callback: () => void) {
    this.shortcuts.set(key.toLowerCase(), callback)
  }

  unregister(key: string) {
    this.shortcuts.delete(key.toLowerCase())
  }

  enable() {
    this.enabled = true
  }

  disable() {
    this.enabled = false
  }

  handle(e: KeyboardEvent) {
    if (!this.enabled) return

    const keys: string[] = []
    if (e.metaKey) keys.push('cmd')
    if (e.ctrlKey) keys.push('ctrl')
    if (e.altKey) keys.push('alt')
    if (e.shiftKey) keys.push('shift')
    keys.push(e.key.toLowerCase())

    const shortcut = keys.join('+')
    const callback = this.shortcuts.get(shortcut)
    
    if (callback) {
      e.preventDefault()
      callback()
    }
  }
}

// 全局快捷键实例
export const globalShortcuts = new ShortcutManager()

// Vue 插件
export function installShortcuts(_app: App) {
  window.addEventListener('keydown', (e) => globalShortcuts.handle(e))
  
  // 注册默认快捷键
  globalShortcuts.register('cmd+k', () => {
    // 打开命令面板 - 通过事件触发
    window.dispatchEvent(new CustomEvent('open-command-palette'))
  })
  
  globalShortcuts.register('ctrl+k', () => {
    window.dispatchEvent(new CustomEvent('open-command-palette'))
  })
}

// 快捷键定义
export const SHORTCUTS = {
  // 文件操作
  NEW_ARTICLE: 'Cmd/Ctrl + N',
  SAVE_ARTICLE: 'Cmd/Ctrl + S',
  SAVE_VERSION: 'Cmd/Ctrl + Shift + S',
  
  // 编辑操作
  BOLD: 'Cmd/Ctrl + B',
  ITALIC: 'Cmd/Ctrl + I',
  UNDO: 'Cmd/Ctrl + Z',
  REDO: 'Cmd/Ctrl + Shift + Z',
  
  // AI操作
  COMMAND_PALETTE: 'Cmd/Ctrl + K',
  SLASH_COMMAND: '/',
  
  // 导航
  PUBLISH: 'Cmd/Ctrl + Shift + P',
  SETTINGS: 'Cmd/Ctrl + ,',
  
  // 通用
  ESCAPE: 'Esc',
  SEARCH: 'Cmd/Ctrl + F'
} as const

// 快捷键说明
export const SHORTCUT_DESCRIPTIONS: Record<string, string> = {
  'Cmd/Ctrl + N': '新建文章',
  'Cmd/Ctrl + S': '保存文章',
  'Cmd/Ctrl + Shift + S': '创建版本',
  'Cmd/Ctrl + B': '加粗',
  'Cmd/Ctrl + I': '斜体',
  'Cmd/Ctrl + Z': '撤销',
  'Cmd/Ctrl + Shift + Z': '重做',
  'Cmd/Ctrl + K': '命令面板',
  'Cmd/Ctrl + Shift + P': '发布文章',
  'Cmd/Ctrl + ,': '打开设置',
  'Esc': '关闭弹窗',
  '/': 'AI斜杠命令'
}
