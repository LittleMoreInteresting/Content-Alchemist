// 防抖函数
export function debounce<T extends (...args: any[]) => void>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timer: number | null = null
  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (timer) {
      clearTimeout(timer)
    }
    timer = window.setTimeout(() => {
      fn.apply(this, args)
      timer = null
    }, delay)
  }
}

// 节流函数
export function throttle<T extends (...args: any[]) => void>(
  fn: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle = false
  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (!inThrottle) {
      fn.apply(this, args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

// 虚拟列表配置
export interface VirtualListConfig {
  itemHeight: number
  overscan: number
}

// 计算虚拟列表可见范围
export function getVirtualListRange(
  scrollTop: number,
  containerHeight: number,
  totalItems: number,
  config: VirtualListConfig
) {
  const { itemHeight, overscan } = config
  
  const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan)
  const visibleCount = Math.ceil(containerHeight / itemHeight)
  const endIndex = Math.min(totalItems, startIndex + visibleCount + overscan * 2)
  
  return {
    startIndex,
    endIndex,
    offsetY: startIndex * itemHeight,
    totalHeight: totalItems * itemHeight
  }
}

// 测量性能
export function measurePerformance<T>(name: string, fn: () => T): T {
  const start = performance.now()
  const result = fn()
  const end = performance.now()
  console.log(`[Performance] ${name}: ${(end - start).toFixed(2)}ms`)
  return result
}

// 懒加载图片
export function lazyLoadImages(container: HTMLElement) {
  const images = container.querySelectorAll('img[data-src]')
  
  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const img = entry.target as HTMLImageElement
        img.src = img.dataset.src || ''
        img.removeAttribute('data-src')
        observer.unobserve(img)
      }
    })
  })
  
  images.forEach((img) => observer.observe(img))
  
  return () => observer.disconnect()
}

// 使用 requestIdleCallback 执行低优先级任务
export function scheduleIdleTask(callback: () => void, timeout = 2000) {
  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(callback, { timeout })
  } else {
    setTimeout(callback, 1)
  }
}

// 内存监控（开发环境）
export function monitorMemory() {
  if (process.env.NODE_ENV === 'development' && 'memory' in performance) {
    setInterval(() => {
      const memory = (performance as any).memory
      console.log('[Memory]', {
        used: (memory.usedJSHeapSize / 1048576).toFixed(2) + ' MB',
        total: (memory.totalJSHeapSize / 1048576).toFixed(2) + ' MB',
        limit: (memory.jsHeapSizeLimit / 1048576).toFixed(2) + ' MB'
      })
    }, 30000)
  }
}
