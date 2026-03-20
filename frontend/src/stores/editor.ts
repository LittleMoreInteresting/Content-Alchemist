import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Article, OutlineNode } from '@/types'

export const useEditorStore = defineStore('editor', () => {
  // State
  const article = ref<Article>({
    id: '',
    title: '',
    content: '# 请输入标题\n\n开始创作...',
    outline: [],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    status: 'draft'
  })
  
  const currentOutlineId = ref('')
  const editorView = ref<'edit' | 'preview' | 'split'>('split')
  const isAIGenerating = ref(false)
  const selectedText = ref('')
  const saveStatus = ref<'saved' | 'saving' | 'unsaved'>('saved')
  const currentTheme = ref('default')

  // Getters
  const wordCount = computed(() => {
    if (!article.value.content) return 0
    // 中文字符 + 英文单词
    const chinese = (article.value.content.match(/[\u4e00-\u9fa5]/g) || []).length
    const english = (article.value.content.match(/[a-zA-Z]+/g) || []).length
    return chinese + english
  })

  const outlineFromContent = computed((): OutlineNode[] => {
    const content = article.value.content || ''
    const lines = content.split('\n')
    const outline: OutlineNode[] = []
    let idCounter = 1

    for (const line of lines) {
      const trimmed = line.trim()
      if (trimmed.startsWith('# ')) {
        outline.push({
          id: String(idCounter++),
          level: 1,
          title: trimmed.substring(2),
          status: 'draft',
          wordCount: 0,
          targetWords: 0
        })
      } else if (trimmed.startsWith('## ')) {
        outline.push({
          id: String(idCounter++),
          level: 2,
          title: trimmed.substring(3),
          status: 'draft',
          wordCount: 0,
          targetWords: 0
        })
      } else if (trimmed.startsWith('### ')) {
        outline.push({
          id: String(idCounter++),
          level: 3,
          title: trimmed.substring(4),
          status: 'draft',
          wordCount: 0,
          targetWords: 0
        })
      }
    }

    return outline
  })

  // Actions
  function setArticle(newArticle: Article) {
    article.value = newArticle
    markUnsaved()
  }

  function updateContent(content: string) {
    article.value.content = content
    article.value.updatedAt = new Date().toISOString()
    markUnsaved()
  }

  function updateTitle(title: string) {
    article.value.title = title
    article.value.updatedAt = new Date().toISOString()
    markUnsaved()
  }

  function updateOutline(outline: OutlineNode[]) {
    article.value.outline = outline
    markUnsaved()
  }

  function setCurrentOutlineId(id: string) {
    currentOutlineId.value = id
  }

  function setEditorView(view: 'edit' | 'preview' | 'split') {
    editorView.value = view
  }

  function setAIGenerating(generating: boolean) {
    isAIGenerating.value = generating
  }

  function setSelectedText(text: string) {
    selectedText.value = text
  }

  function markSaved() {
    saveStatus.value = 'saved'
  }

  function markSaving() {
    saveStatus.value = 'saving'
  }

  function markUnsaved() {
    saveStatus.value = 'unsaved'
  }

  function setTheme(theme: string) {
    currentTheme.value = theme
  }

  function reset() {
    article.value = {
      id: '',
      title: '',
      content: '# 请输入标题\n\n开始创作...',
      outline: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      status: 'draft'
    }
    currentOutlineId.value = ''
    editorView.value = 'split'
    isAIGenerating.value = false
    selectedText.value = ''
    saveStatus.value = 'saved'
  }

  return {
    article,
    currentOutlineId,
    editorView,
    isAIGenerating,
    selectedText,
    saveStatus,
    currentTheme,
    wordCount,
    outlineFromContent,
    setArticle,
    updateContent,
    updateTitle,
    updateOutline,
    setCurrentOutlineId,
    setEditorView,
    setAIGenerating,
    setSelectedText,
    markSaved,
    markSaving,
    markUnsaved,
    setTheme,
    reset
  }
})
