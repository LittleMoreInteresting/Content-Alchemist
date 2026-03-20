import { marked, type RendererObject } from 'marked'

// 配置 marked
marked.setOptions({
  gfm: true,
  breaks: true
})

// 主题样式定义
const themeStyles: Record<string, Record<string, string>> = {
  default: {
    '--primary-color': '#1677ff',
    '--heading-font': '"PingFang SC", "Microsoft YaHei", sans-serif',
    '--body-font': '"PingFang SC", "Microsoft YaHei", sans-serif',
    '--code-bg': '#f6f8fa',
    '--quote-bg': 'linear-gradient(120deg, #e0c3fc 0%, #8ec5fc 100%)',
    '--h1-style': 'center',
    '--h1-border': '2px solid #1677ff'
  },
  minimal: {
    '--primary-color': '#333',
    '--heading-font': '"PingFang SC", sans-serif',
    '--body-font': '"PingFang SC", sans-serif',
    '--code-bg': '#f5f5f5',
    '--quote-bg': '#f5f5f5',
    '--h1-style': 'left',
    '--h1-border': 'none'
  },
  vibrant: {
    '--primary-color': '#ff6b35',
    '--heading-font': '"PingFang SC", sans-serif',
    '--body-font': '"PingFang SC", sans-serif',
    '--code-bg': '#fff5f0',
    '--quote-bg': 'linear-gradient(120deg, #fccb90 0%, #d57eeb 100%)',
    '--h1-style': 'center',
    '--h1-border': '3px solid #ff6b35'
  }
}

// 自定义渲染器
const renderer: RendererObject = {
  // 标题渲染
  heading({ tokens, depth }) {
    const text = this.parser.parseInline(tokens)
    if (depth === 1) {
      return `<h1 style="
        font-size: 24px;
        font-weight: 700;
        color: var(--primary-color);
        text-align: center;
        padding: 20px 0;
        margin: 0 0 24px 0;
        border-bottom: var(--h1-border);
        position: relative;
      ">${text}</h1>`
    } else if (depth === 2) {
      return `<h2 style="
        font-size: 20px;
        font-weight: 600;
        color: #07c160;
        padding: 16px;
        margin: 24px 0 16px;
        background: rgba(7, 193, 96, 0.08);
        border-left: 4px solid #07c160;
        border-radius: 0 8px 8px 0;
      ">${text}</h2>`
    } else {
      return `<h3 style="
        font-size: 17px;
        font-weight: 600;
        color: #576b95;
        padding: 12px 0;
        margin: 20px 0 12px;
        border-bottom: 1px solid #e8e8e8;
      ">${text}</h3>`
    }
  },

  // 段落渲染
  paragraph({ tokens }) {
    const text = this.parser.parseInline(tokens)
    return `<p style="
      font-size: 16px;
      line-height: 1.8;
      color: #333;
      margin: 16px 0;
      text-align: justify;
    ">${text}</p>`
  },

  // 代码块渲染
  code({ text, lang }) {
    const language = lang || 'code'
    return `<div style="
      margin: 16px 0;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    ">
      <div style="
        background: #2d2d2d;
        padding: 8px 16px;
        display: flex;
        gap: 6px;
        align-items: center;
      ">
        <span style="width: 12px; height: 12px; border-radius: 50%; background: #ff5f56;"></span>
        <span style="width: 12px; height: 12px; border-radius: 50%; background: #ffbd2e;"></span>
        <span style="width: 12px; height: 12px; border-radius: 50%; background: #27c93f;"></span>
        <span style="color: #999; font-size: 12px; margin-left: 8px;">${language}</span>
      </div>
      <pre style="
        background: #1e1e1e;
        color: #d4d4d4;
        padding: 16px;
        margin: 0;
        overflow-x: auto;
        font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
        font-size: 14px;
        line-height: 1.6;
      "><code>${escapeHtml(text)}</code></pre>
    </div>`
  },

  // 引用块渲染
  blockquote({ tokens }) {
    const text = this.parser.parse(tokens)
    return `<blockquote style="
      margin: 20px 0;
      padding: 20px 24px;
      background: var(--quote-bg);
      border-radius: 8px;
      position: relative;
      font-style: italic;
    ">
      <div style="
        position: absolute;
        top: 10px;
        left: 16px;
        font-size: 48px;
        color: rgba(0,0,0,0.1);
        font-family: Georgia, serif;
        line-height: 1;
      ">"</div>
      <div style="position: relative; z-index: 1;">${text}</div>
    </blockquote>`
  },

  // 列表渲染
  list({ items, ordered }) {
    const tag = ordered ? 'ol' : 'ul'
    const style = ordered 
      ? 'list-style: decimal; padding-left: 24px;'
      : 'list-style: disc; padding-left: 24px;'
    
    const listItems = items.map(item => `<li style="
      margin: 8px 0;
      line-height: 1.8;
    ">${this.parser.parse(item.tokens)}</li>`).join('')
    
    return `<${tag} style="
      margin: 16px 0;
      ${style}
      color: #333;
    ">${listItems}</${tag}>`
  },

  // 表格渲染
  table({ header, rows }) {
    const headerCells = header.map(cell => `<th style="
      padding: 12px;
      text-align: left;
      font-weight: 600;
      background: #07c160;
      color: white;
    ">${this.parser.parseInline(cell.tokens)}</th>`).join('')
    
    const bodyRows = rows.map(row => `<tr style="
      border-bottom: 1px solid #e8e8e8;
    ">${row.map(cell => `<td style="
      padding: 12px;
      text-align: left;
      background: #fafafa;
    ">${this.parser.parseInline(cell.tokens)}</td>`).join('')}</tr>`).join('')
    
    return `<div style="
      overflow-x: auto;
      margin: 20px 0;
    ">
      <table style="
        width: 100%;
        border-collapse: collapse;
        font-size: 14px;
      ">
        <thead>${headerCells}</thead>
        <tbody>${bodyRows}</tbody>
      </table>
    </div>`
  },

  // 链接渲染
  link({ href, title, tokens }) {
    const text = this.parser.parseInline(tokens)
    return `<a href="${href}" title="${title || ''}" style="
      color: #576b95;
      text-decoration: none;
      border-bottom: 1px solid #576b95;
    " target="_blank">${text}</a>`
  },

  // 图片渲染
  image({ href, title, text }) {
    return `<img src="${href}" alt="${text}" title="${title || ''}" style="
      max-width: 100%;
      height: auto;
      border-radius: 8px;
      margin: 16px 0;
    ">`
  },

  // 行内代码
  codespan({ text }) {
    return `<code style="
      background: var(--code-bg);
      padding: 2px 6px;
      border-radius: 4px;
      font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
      font-size: 0.9em;
      color: #e83e8c;
    ">${escapeHtml(text)}</code>`
  },

  // 粗体
  strong({ tokens }) {
    const text = this.parser.parseInline(tokens)
    return `<strong style="font-weight: 700; color: #000;">${text}</strong>`
  },

  // 斜体
  em({ tokens }) {
    const text = this.parser.parseInline(tokens)
    return `<em style="font-style: italic; color: #666;">${text}</em>`
  },

  // 分割线
  hr() {
    return `<hr style="
      border: none;
      height: 1px;
      background: linear-gradient(to right, transparent, #ddd, transparent);
      margin: 32px 0;
    ">`
  }
}

// 注册自定义渲染器
marked.use({ renderer })

// HTML转义
function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  }
  return text.replace(/[&<>"']/g, m => map[m])
}

// 渲染微信HTML
export function renderWechatHTML(markdown: string, theme: string = 'default'): string {
  const themeStyle = themeStyles[theme] || themeStyles.default
  
  // 构建CSS变量
  const cssVars = Object.entries(themeStyle)
    .map(([key, value]) => `${key}: ${value};`)
    .join(' ')
  
  const html = marked.parse(markdown) as string
  
  return `<div style="${cssVars}; font-family: var(--body-font);">${html}</div>`
}

// 导出纯文本（用于字数统计等）
export function extractText(markdown: string): string {
  return markdown
    .replace(/#+ /g, '')
    .replace(/\*\*/g, '')
    .replace(/\*/g, '')
    .replace(/`/g, '')
    .replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1')
    .replace(/!\[([^\]]*)\]\([^\)]+\)/g, '')
}
