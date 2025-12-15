import { FC, useState, useEffect, useCallback } from 'react';
import { Toast as ToastType } from '@/context/ToastContext';
import './Toast.css';

interface ToastItemProps {
  toast: ToastType;
  onRemove: (id: string) => void;
}

const ToastItem: FC<ToastItemProps> = ({ toast, onRemove }) => {
  const [isExiting, setIsExiting] = useState(false);

  const handleClose = useCallback(() => {
    setIsExiting(true);
    // Wait for animation to complete before removing
    setTimeout(() => {
      onRemove(toast.id);
    }, 300);
  }, [onRemove, toast.id]);

  // Auto-close animation when duration expires
  useEffect(() => {
    if (toast.duration && toast.duration > 0) {
      const timer = setTimeout(() => {
        handleClose();
      }, toast.duration - 300); // Start exit animation 300ms before removal

      return () => clearTimeout(timer);
    }
  }, [toast.duration, handleClose]);

  const getIcon = () => {
    switch (toast.type) {
      case 'success':
        return '✓';
      case 'error':
        return '✕';
      case 'warning':
        return '⚠';
      case 'info':
      default:
        return 'ℹ';
    }
  };

  return (
    <div className={`toast ${toast.type} ${isExiting ? 'exiting' : ''}`} role="alert">
      <div className="toast-icon" aria-hidden="true">
        {getIcon()}
      </div>
      <div className="toast-message">{toast.message}</div>
      <button
        className="toast-close"
        onClick={handleClose}
        aria-label="Close notification"
        type="button"
      >
        ×
      </button>
    </div>
  );
};

export default ToastItem;
