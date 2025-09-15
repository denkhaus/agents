/**
 * Markdown utilities
 * Using established libraries for secure and complete markdown rendering
 */

import { marked } from 'marked';
import DOMPurify from 'dompurify';

// Configure marked for security and consistency
marked.setOptions({
  breaks: true,
  gfm: true,
  headerIds: false,
  mangle: false
});

/**
 * Parse markdown to HTML using marked library
 */
export function parseMarkdown(markdown: string): string {
  if (!markdown) return '';
  
  try {
    const html = marked.parse(markdown);
    return sanitizeHTML(html);
  } catch (error) {
    console.error('Markdown parsing error:', error);
    return markdown; // Fallback to plain text
  }
}

/**
 * Sanitize HTML using DOMPurify to prevent XSS
 */
export function sanitizeHTML(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'p', 'br', 'strong', 'em', 'u', 's',
      'ul', 'ol', 'li',
      'blockquote', 'pre', 'code',
      'a', 'img',
      'table', 'thead', 'tbody', 'tr', 'th', 'td'
    ],
    ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class'],
    ALLOW_DATA_ATTR: false,
    ADD_ATTR: ['target', 'rel'],
    FORBID_TAGS: ['script', 'style', 'iframe', 'object', 'embed'],
    FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover']
  });
}

// Extract plain text from markdown (for previews)
export function markdownToPlainText(markdown: string): string {
  if (!markdown) return '';
  
  return markdown
    // Remove headers
    .replace(/^#{1,6}\s+/gm, '')
    // Remove bold/italic
    .replace(/\*\*(.*?)\*\*/g, '$1')
    .replace(/\*(.*?)\*/g, '$1')
    .replace(/__(.*?)__/g, '$1')
    .replace(/_(.*?)_/g, '$1')
    // Remove code blocks
    .replace(/```[\s\S]*?```/g, '')
    .replace(/`(.*?)`/g, '$1')
    // Remove links
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    // Remove list markers
    .replace(/^[\*\-\+]\s+/gm, '')
    // Clean up whitespace
    .replace(/\n+/g, ' ')
    .trim();
}

// Truncate text for previews
export function truncateText(text: string, maxLength: number = 100): string {
  if (text.length <= maxLength) return text;
  
  return text.substring(0, maxLength).trim() + '...';
}

// Check if text contains markdown
export function hasMarkdown(text: string): boolean {
  if (!text) return false;
  
  const markdownPatterns = [
    /^#{1,6}\s+/m,     // Headers
    /\*\*.*?\*\*/,     // Bold
    /\*.*?\*/,         // Italic
    /`.*?`/,           // Inline code
    /```[\s\S]*?```/,  // Code blocks
    /\[.*?\]\(.*?\)/,  // Links
    /^[\*\-\+]\s+/m    // Lists
  ];
  
  return markdownPatterns.some(pattern => pattern.test(text));
}