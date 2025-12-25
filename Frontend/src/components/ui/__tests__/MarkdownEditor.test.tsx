import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { MarkdownEditor } from '../MarkdownEditor';

vi.mock('../../../hooks/useTranslation', () => ({
  useT: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'markdown.write': 'Write',
        'markdown.preview': 'Preview',
        'markdown.helpLink': 'Markdown is supported',
        'markdown.noContent': 'Nothing to preview'
      };
      return translations[key] || key;
    }
  })
}));

describe('MarkdownEditor Component', () => {
  const defaultProps = {
    value: '',
    onChange: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with Write and Preview tabs', () => {
    render(<MarkdownEditor {...defaultProps} />);
    
    expect(screen.getByRole('button', { name: /write/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /preview/i })).toBeInTheDocument();
  });

  it('renders with label when provided', () => {
    render(<MarkdownEditor {...defaultProps} label="Content" />);
    
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('starts in write mode by default', () => {
    render(<MarkdownEditor {...defaultProps} />);
    
    const writeTab = screen.getByRole('button', { name: /write/i });
    const previewTab = screen.getByRole('button', { name: /preview/i });
    
    expect(writeTab).toHaveClass('active');
    expect(previewTab).not.toHaveClass('active');
    expect(screen.getByRole('textbox')).toBeInTheDocument();
  });

  it('switches to preview mode when Preview tab is clicked', async () => {
    render(<MarkdownEditor {...defaultProps} value="**Bold text**" />);
    
    const previewTab = screen.getByRole('button', { name: /preview/i });
    fireEvent.click(previewTab);
    
    expect(previewTab).toHaveClass('active');
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByText('Bold text')).toBeInTheDocument();
    });
  });

  it('switches back to write mode when Write tab is clicked', () => {
    render(<MarkdownEditor {...defaultProps} value="Test content" />);
    
    const writeTab = screen.getByRole('button', { name: /write/i });
    const previewTab = screen.getByRole('button', { name: /preview/i });
    
    fireEvent.click(previewTab);
    expect(screen.queryByRole('textbox')).not.toBeInTheDocument();
    
    fireEvent.click(writeTab);
    expect(screen.getByRole('textbox')).toBeInTheDocument();
    expect(writeTab).toHaveClass('active');
  });

  it('calls onChange when text is entered', () => {
    const handleChange = vi.fn();
    render(<MarkdownEditor {...defaultProps} onChange={handleChange} />);
    
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: 'New content' } });
    
    expect(handleChange).toHaveBeenCalledWith('New content');
  });

  it('displays placeholder text', () => {
    render(<MarkdownEditor {...defaultProps} placeholder="Enter content here" />);
    
    const textarea = screen.getByRole('textbox');
    expect(textarea).toHaveAttribute('placeholder', 'Enter content here');
  });

  it('disables tabs and textarea when disabled prop is true', () => {
    render(<MarkdownEditor {...defaultProps} disabled={true} />);
    
    const writeTab = screen.getByRole('button', { name: /write/i });
    const previewTab = screen.getByRole('button', { name: /preview/i });
    const textarea = screen.getByRole('textbox');
    
    expect(writeTab).toBeDisabled();
    expect(previewTab).toBeDisabled();
    expect(textarea).toBeDisabled();
  });

  it('displays error message when error prop is provided', () => {
    render(<MarkdownEditor {...defaultProps} error="This field is required" />);
    
    expect(screen.getByText('This field is required')).toBeInTheDocument();
  });

  it('sets aria-invalid on textarea when there is an error', () => {
    render(<MarkdownEditor {...defaultProps} error="This field is required" />);
    
    const textarea = screen.getByRole('textbox');
    expect(textarea).toHaveAttribute('aria-invalid', 'true');
  });

  it('associates error message with textarea via aria-describedby', () => {
    render(<MarkdownEditor {...defaultProps} error="This field is required" id="test-editor" />);
    
    const textarea = screen.getByRole('textbox');
    const errorMessage = screen.getByText('This field is required');
    
    expect(textarea).toHaveAttribute('aria-describedby', 'test-editor-error');
    expect(errorMessage).toHaveAttribute('id', 'test-editor-error');
  });

  it('does not set aria-invalid when there is no error', () => {
    render(<MarkdownEditor {...defaultProps} />);
    
    const textarea = screen.getByRole('textbox');
    expect(textarea).not.toHaveAttribute('aria-invalid');
  });

  it('renders markdown help link', () => {
    render(<MarkdownEditor {...defaultProps} />);
    
    const link = screen.getByRole('link', { name: /markdown is supported/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', expect.stringContaining('github.com'));
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('applies custom id when provided', () => {
    render(<MarkdownEditor {...defaultProps} id="custom-id" label="Custom" />);
    
    const textarea = screen.getByRole('textbox');
    expect(textarea).toHaveAttribute('id', 'custom-id');
  });

  it('renders markdown content in preview mode', async () => {
    render(<MarkdownEditor {...defaultProps} value="# Heading" />);
    
    const previewTab = screen.getByRole('button', { name: /preview/i });
    fireEvent.click(previewTab);
    
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: 'Heading' })).toBeInTheDocument();
    });
  });

  it('renders bold text in preview mode', async () => {
    render(<MarkdownEditor {...defaultProps} value="**bold text**" />);
    
    const previewTab = screen.getByRole('button', { name: /preview/i });
    fireEvent.click(previewTab);
    
    await waitFor(() => {
      expect(screen.getByText('bold text')).toBeInTheDocument();
    });
  });

  it('sets custom number of rows', () => {
    render(<MarkdownEditor {...defaultProps} rows={12} />);
    
    const textarea = screen.getByRole('textbox');
    expect(textarea).toHaveAttribute('rows', '12');
  });

  it('applies correct CSS classes', () => {
    const { container } = render(<MarkdownEditor {...defaultProps} />);
    
    expect(container.querySelector('.markdown-editor')).toBeInTheDocument();
    expect(container.querySelector('.markdown-editor-header')).toBeInTheDocument();
    expect(container.querySelector('.markdown-editor-tabs')).toBeInTheDocument();
    expect(container.querySelector('.markdown-editor-content')).toBeInTheDocument();
    expect(container.querySelector('.markdown-editor-footer')).toBeInTheDocument();
  });

  it('shows preview container with data-placeholder when content is empty', () => {
    const { container } = render(<MarkdownEditor {...defaultProps} value="" />);
    
    const previewTab = screen.getByRole('button', { name: /preview/i });
    fireEvent.click(previewTab);
    
    const preview = container.querySelector('.markdown-editor-preview');
    expect(preview).toHaveAttribute('data-placeholder', 'Nothing to preview');
  });
});
