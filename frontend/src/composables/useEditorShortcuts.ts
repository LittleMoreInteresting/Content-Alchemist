import { onMounted, onUnmounted } from 'vue'

interface ShortcutHandler {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  meta?: boolean
  handler: (e: KeyboardEvent) => void
  preventDefault?: boolean
}

export function useEditorShortcuts(shortcuts: ShortcutHandler[]) {
  function handleKeyDown(e: KeyboardEvent) {
    for (const shortcut of shortcuts) {
      const keyMatch = e.key.toLowerCase() === shortcut.key.toLowerCase()
      const ctrlMatch = !!shortcut.ctrl === (e.ctrlKey || e.metaKey)
      const shiftMatch = !!shortcut.shift === e.shiftKey
      const altMatch = !!shortcut.alt === e.altKey
      
      if (keyMatch && ctrlMatch && shiftMatch && altMatch) {
        if (shortcut.preventDefault !== false) {
          e.preventDefault()
        }
        shortcut.handler(e)
        break
      }
    }
  }
  
  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown)
  })
  
  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
  })
  
  return {
    register: (shortcut: ShortcutHandler) => shortcuts.push(shortcut)
  }
}

// 常用编辑器快捷键
export const EDITOR_SHORTCUTS = {
  bold: { key: 'b', ctrl: true },
  italic: { key: 'i', ctrl: true },
  save: { key: 's', ctrl: true },
  find: { key: 'f', ctrl: true },
  undo: { key: 'z', ctrl: true },
  redo: { key: 'z', ctrl: true, shift: true },
} as const
