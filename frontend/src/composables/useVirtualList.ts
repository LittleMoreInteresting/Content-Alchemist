import { ref, computed, onMounted, onUnmounted } from 'vue'

interface UseVirtualListOptions {
  itemHeight: number
  overscan?: number
}

export function useVirtualList<T>(
  items: T[],
  options: UseVirtualListOptions
) {
  const { itemHeight, overscan = 5 } = options
  
  const containerRef = ref<HTMLElement | null>(null)
  const scrollTop = ref(0)
  const containerHeight = ref(0)
  
  const totalHeight = computed(() => items.length * itemHeight)
  
  const visibleRange = computed(() => {
    const start = Math.max(0, Math.floor(scrollTop.value / itemHeight) - overscan)
    const visibleCount = Math.ceil(containerHeight.value / itemHeight)
    const end = Math.min(items.length, start + visibleCount + overscan * 2)
    
    return { start, end }
  })
  
  const visibleItems = computed(() => {
    const { start, end } = visibleRange.value
    return items.slice(start, end).map((item, index) => ({
      item,
      index: start + index
    }))
  })
  
  const offsetY = computed(() => visibleRange.value.start * itemHeight)
  
  function onScroll(e: Event) {
    scrollTop.value = (e.target as HTMLElement).scrollTop
  }
  
  onMounted(() => {
    if (containerRef.value) {
      containerHeight.value = containerRef.value.clientHeight
      containerRef.value.addEventListener('scroll', onScroll, { passive: true })
    }
  })
  
  onUnmounted(() => {
    if (containerRef.value) {
      containerRef.value.removeEventListener('scroll', onScroll)
    }
  })
  
  return {
    containerRef,
    visibleItems,
    totalHeight,
    offsetY,
    itemHeight
  }
}
