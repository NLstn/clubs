import React, { useState, useId, useCallback } from 'react';
import Markdown from 'react-markdown';
import { useT } from '@/hooks/useTranslation';
import './MarkdownEditor.css';

export interface MarkdownEditorProps {
  label?: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  disabled?: boolean;
  error?: string;
  rows?: number;
  id?: string;
}

type EditorMode = 'write' | 'preview';

export const MarkdownEditor: React.FC<MarkdownEditorProps> = ({
  label,
  value,
  onChange,
  placeholder = '',
  disabled = false,
  error,
  rows = 8,
  id,
}) => {
  const { t } = useT();
  const [mode, setMode] = useState<EditorMode>('write');
  const generatedId = useId();
  const editorId = id || generatedId;
  const errorId = `${editorId}-error`;

  const handleTextChange = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      onChange(e.target.value);
    },
    [onChange]
  );

  const handleModeChange = (newMode: EditorMode) => {
    setMode(newMode);
  };

  return (
    <div className="markdown-editor">
      <div className="markdown-editor-header">
        {label && (
          <label htmlFor={editorId} className="markdown-editor-label">
            {label}
          </label>
        )}
        <div className="markdown-editor-tabs">
          <button
            type="button"
            className={`markdown-editor-tab ${mode === 'write' ? 'active' : ''}`}
            onClick={() => handleModeChange('write')}
            disabled={disabled}
          >
            {t('markdown.write')}
          </button>
          <button
            type="button"
            className={`markdown-editor-tab ${mode === 'preview' ? 'active' : ''}`}
            onClick={() => handleModeChange('preview')}
            disabled={disabled}
          >
            {t('markdown.preview')}
          </button>
        </div>
      </div>

      <div className="markdown-editor-content">
        {mode === 'write' ? (
          <textarea
            id={editorId}
            className="markdown-editor-textarea"
            value={value}
            onChange={handleTextChange}
            placeholder={placeholder}
            disabled={disabled}
            rows={rows}
            aria-invalid={error ? 'true' : undefined}
            aria-describedby={error ? errorId : undefined}
          />
        ) : (
          <div
            className="markdown-editor-preview"
            data-placeholder={placeholder || t('markdown.noContent')}
            style={{ minHeight: `${rows * 1.5}rem` }}
          >
            {value ? (
              <Markdown>{value}</Markdown>
            ) : null}
          </div>
        )}
      </div>

      {error && <span id={errorId} className="markdown-editor-error">{error}</span>}

      <div className="markdown-editor-footer">
        <a
          href="https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax"
          target="_blank"
          rel="noopener noreferrer"
          className="markdown-editor-help-link"
        >
          {t('markdown.helpLink')}
        </a>
      </div>
    </div>
  );
};

export default MarkdownEditor;
