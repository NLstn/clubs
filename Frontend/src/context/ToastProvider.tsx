import { FC, ReactNode, useState, useCallback, useMemo } from 'react';
import { ToastContext, Toast, ToastType } from './ToastContext';

interface ToastProviderProps {
  children: ReactNode;
}

export const ToastProvider: FC<ToastProviderProps> = ({ children }) => {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const addToast = useCallback((message: string, type: ToastType = 'info', duration: number = 5000) => {
    const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const newToast: Toast = { id, message, type, duration };
    
    setToasts((prev) => [...prev, newToast]);

    // Auto-remove toast after duration
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id);
      }, duration);
    }
  }, [removeToast]);

  const success = useCallback((message: string, duration?: number) => {
    addToast(message, 'success', duration);
  }, [addToast]);

  const error = useCallback((message: string, duration?: number) => {
    addToast(message, 'error', duration);
  }, [addToast]);

  const info = useCallback((message: string, duration?: number) => {
    addToast(message, 'info', duration);
  }, [addToast]);

  const warning = useCallback((message: string, duration?: number) => {
    addToast(message, 'warning', duration);
  }, [addToast]);

  const value = useMemo(() => ({
    toasts,
    addToast,
    removeToast,
    success,
    error,
    info,
    warning,
  }), [toasts, addToast, removeToast, success, error, info, warning]);

  return (
    <ToastContext.Provider value={value}>
      {children}
    </ToastContext.Provider>
  );
};
