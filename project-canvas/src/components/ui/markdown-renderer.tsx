/**
 * Markdown Renderer Component
 * Safely renders markdown content from task/project descriptions
 */

import React from 'react';
import { parseMarkdown, hasMarkdown, markdownToPlainText } from '../../utils/markdown';
import { cn } from '../../lib/utils';

interface MarkdownRendererProps {
  content: string;
  className?: string;
  maxLength?: number;
  showPreview?: boolean;
}

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
  content,
  className = '',
  maxLength,
  showPreview = false
}) => {
  const [htmlContent, setHtmlContent] = React.useState<string>('');
  
  React.useEffect(() => {
    if (content && !showPreview && hasMarkdown(content)) {
      parseMarkdown(content).then(setHtmlContent);
    }
  }, [content, showPreview]);

  if (!content) {
    return <span className={`text-muted-foreground ${className}`}>No description</span>;
  }

  // If showing preview or content is short, show plain text
  if (showPreview || !hasMarkdown(content)) {
    const plainText = hasMarkdown(content) ? markdownToPlainText(content) : content;
    const displayText = maxLength ? 
      (plainText.length > maxLength ? plainText.substring(0, maxLength) + '...' : plainText) : 
      plainText;
    
    return (
      <span className={`${className}`}>
        {displayText}
      </span>
    );
  }

  // Render full markdown
  return (
    <div 
      className={cn(
        "prose prose-sm dark:prose-invert max-w-none",
        "prose-headings:font-semibold prose-headings:text-foreground",
        "prose-p:text-foreground prose-p:leading-relaxed",
        "prose-a:text-primary hover:prose-a:text-primary/80",
        "prose-strong:text-foreground prose-em:text-foreground",
        "prose-code:bg-muted prose-code:px-1 prose-code:py-0.5 prose-code:rounded prose-code:text-sm",
        "prose-pre:bg-muted prose-pre:border",
        "prose-blockquote:border-l-primary prose-blockquote:bg-muted/50",
        "prose-ul:text-foreground prose-ol:text-foreground",
        "prose-li:text-foreground",
        "prose-table:text-foreground prose-th:border prose-td:border",
        className
      )}
      dangerouslySetInnerHTML={{ __html: htmlContent }}
    />
  );
};

interface EditableMarkdownProps {
  content: string;
  onSave: (newContent: string) => void;
  placeholder?: string;
  className?: string;
}

export const EditableMarkdown: React.FC<EditableMarkdownProps> = ({
  content,
  onSave,
  placeholder = "Enter description...",
  className = ''
}) => {
  const [isEditing, setIsEditing] = React.useState(false);
  const [editValue, setEditValue] = React.useState(content);

  const handleSave = () => {
    onSave(editValue);
    setIsEditing(false);
  };

  const handleCancel = () => {
    setEditValue(content);
    setIsEditing(false);
  };

  if (isEditing) {
    return (
      <div className={`space-y-2 ${className}`}>
        <textarea
          value={editValue}
          onChange={(e) => setEditValue(e.target.value)}
          placeholder={placeholder}
          className="w-full min-h-[100px] p-2 border rounded-md resize-vertical"
          autoFocus
        />
        <div className="flex gap-2">
          <button
            onClick={handleSave}
            className="px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600"
          >
            Save
          </button>
          <button
            onClick={handleCancel}
            className="px-3 py-1 bg-gray-500 text-white rounded text-sm hover:bg-gray-600"
          >
            Cancel
          </button>
        </div>
        {hasMarkdown(editValue) && (
          <div className="mt-2 p-2 bg-gray-50 dark:bg-gray-800 rounded">
            <p className="text-xs text-muted-foreground mb-1">Preview:</p>
            <MarkdownRenderer content={editValue} />
          </div>
        )}
      </div>
    );
  }

  return (
    <div 
      className={`cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 p-2 rounded ${className}`}
      onClick={() => setIsEditing(true)}
    >
      <MarkdownRenderer content={content} />
      {!content && (
        <span className="text-muted-foreground italic">{placeholder}</span>
      )}
    </div>
  );
};