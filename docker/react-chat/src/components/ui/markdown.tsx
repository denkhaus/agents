'use client'

import React from 'react'

interface MarkdownProps {
  children: string
  className?: string
}

export function Markdown({ children, className = '' }: MarkdownProps) {
  // Simple markdown parsing for common patterns
  const parseMarkdown = (text: string): React.ReactNode[] => {
    const lines = text.split('\n')
    const elements: React.ReactNode[] = []
    let inCodeBlock = false
    let codeBlockContent: string[] = []
    let codeBlockLanguage = ''
    let inTable = false
    let tableRows: string[] = []

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i]

      // Code blocks
      if (line.startsWith('```')) {
        if (inCodeBlock) {
          // End code block
          elements.push(
            <div key={`code-${i}`} className="my-3">
              <div className="bg-gray-100 dark:bg-gray-800 rounded-t px-3 py-1 text-xs text-gray-600 dark:text-gray-400 border-b">
                {codeBlockLanguage || 'code'}
              </div>
              <pre className="bg-gray-50 dark:bg-gray-900 p-3 rounded-b overflow-x-auto">
                <code className="text-sm font-mono whitespace-pre">
                  {codeBlockContent.join('\n')}
                </code>
              </pre>
            </div>
          )
          inCodeBlock = false
          codeBlockContent = []
          codeBlockLanguage = ''
        } else {
          // Start code block
          inCodeBlock = true
          codeBlockLanguage = line.slice(3).trim()
        }
        continue
      }

      if (inCodeBlock) {
        codeBlockContent.push(line)
        continue
      }

      // Tables
      if (line.includes('|') && !inTable) {
        inTable = true
        tableRows = [line]
        continue
      }

      if (inTable) {
        if (line.includes('|')) {
          tableRows.push(line)
          continue
        } else {
          // End table
          elements.push(renderTable(tableRows, i))
          inTable = false
          tableRows = []
          // Continue processing this line as normal
        }
      }

      // Headers
      if (line.startsWith('### ')) {
        elements.push(
          <h3 key={`h3-${i}`} className="text-lg font-semibold mt-4 mb-2">
            {line.slice(4)}
          </h3>
        )
        continue
      }
      if (line.startsWith('## ')) {
        elements.push(
          <h2 key={`h2-${i}`} className="text-xl font-semibold mt-4 mb-2">
            {line.slice(3)}
          </h2>
        )
        continue
      }
      if (line.startsWith('# ')) {
        elements.push(
          <h1 key={`h1-${i}`} className="text-2xl font-bold mt-4 mb-2">
            {line.slice(2)}
          </h1>
        )
        continue
      }

      // Blockquotes
      if (line.startsWith('> ')) {
        elements.push(
          <blockquote key={`quote-${i}`} className="border-l-4 border-gray-300 dark:border-gray-600 pl-4 py-2 my-2 bg-gray-50 dark:bg-gray-800 italic">
            {parseInlineMarkdown(line.slice(2))}
          </blockquote>
        )
        continue
      }

      // Horizontal rule
      if (line.trim() === '---' || line.trim() === '***') {
        elements.push(
          <hr key={`hr-${i}`} className="my-4 border-gray-300 dark:border-gray-600" />
        )
        continue
      }

      // Lists
      if (line.startsWith('- ') || line.startsWith('* ')) {
        elements.push(
          <ul key={`ul-${i}`} className="ml-4 my-2">
            <li className="list-disc">
              {parseInlineMarkdown(line.slice(2))}
            </li>
          </ul>
        )
        continue
      }

      // Numbered lists
      if (/^\d+\.\s/.test(line)) {
        const match = line.match(/^\d+\.\s(.*)/)
        if (match) {
          elements.push(
            <ol key={`ol-${i}`} className="ml-4 my-2">
              <li className="list-decimal">
                {parseInlineMarkdown(match[1])}
              </li>
            </ol>
          )
          continue
        }
      }

      // Empty lines
      if (line.trim() === '') {
        elements.push(<br key={`br-${i}`} />)
        continue
      }

      // Regular paragraphs
      elements.push(
        <p key={`p-${i}`} className="mb-2">
          {parseInlineMarkdown(line)}
        </p>
      )
    }

    // Handle any remaining table at end of text
    if (inTable && tableRows.length > 0) {
      elements.push(renderTable(tableRows, lines.length))
    }

    return elements
  }

  const renderTable = (rows: string[], key: number): React.ReactNode => {
    if (rows.length < 2) return null

    const parseRow = (row: string) => 
      row.split('|').map(cell => cell.trim()).filter(cell => cell !== '')

    const headerRow = parseRow(rows[0])
    const dataRows = rows.slice(2).map(parseRow) // Skip header separator row

    return (
      <div key={`table-${key}`} className="my-3 overflow-x-auto">
        <table className="min-w-full border border-gray-200 dark:border-gray-700 rounded-lg">
          <thead className="bg-gray-50 dark:bg-gray-800">
            <tr>
              {headerRow.map((header, index) => (
                <th key={index} className="px-3 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider border-b border-gray-200 dark:border-gray-700">
                  {parseInlineMarkdown(header)}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            {dataRows.map((row, rowIndex) => (
              <tr key={rowIndex} className="hover:bg-gray-50 dark:hover:bg-gray-800">
                {row.map((cell, cellIndex) => (
                  <td key={cellIndex} className="px-3 py-2 text-sm text-gray-900 dark:text-gray-100 border-b border-gray-200 dark:border-gray-700">
                    {parseInlineMarkdown(cell)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
  }

  const parseInlineMarkdown = (text: string): React.ReactNode => {
    // Inline code
    text = text.replace(/`([^`]+)`/g, '<code class="bg-gray-100 dark:bg-gray-800 px-1 py-0.5 rounded text-sm font-mono text-red-600 dark:text-red-400">$1</code>')
    
    // Bold
    text = text.replace(/\*\*([^*]+)\*\*/g, '<strong class="font-semibold">$1</strong>')
    text = text.replace(/__([^_]+)__/g, '<strong class="font-semibold">$1</strong>')
    
    // Italic
    text = text.replace(/\*([^*]+)\*/g, '<em class="italic">$1</em>')
    text = text.replace(/_([^_]+)_/g, '<em class="italic">$1</em>')
    
    // Strikethrough
    text = text.replace(/~~([^~]+)~~/g, '<del class="line-through">$1</del>')
    
    // Links
    text = text.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" class="text-blue-600 hover:text-blue-800 underline" target="_blank" rel="noopener noreferrer">$1</a>')
    
    // Auto-links for URLs
    text = text.replace(/(https?:\/\/[^\s]+)/g, '<a href="$1" class="text-blue-600 hover:text-blue-800 underline" target="_blank" rel="noopener noreferrer">$1</a>')
    
    return <span dangerouslySetInnerHTML={{ __html: text }} />
  }

  return (
    <div className={`markdown-content ${className}`}>
      {parseMarkdown(children)}
    </div>
  )
}