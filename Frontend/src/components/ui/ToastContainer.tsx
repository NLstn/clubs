import { FC } from 'react';
import { useToast } from '@/hooks/useToast';
import ToastItem from './ToastItem';
import './Toast.css';

const ToastContainer: FC = () => {
  const { toasts, removeToast } = useToast();

  if (toasts.length === 0) {
    return null;
  }

  return (
    <div className="toast-container" aria-live="polite" aria-atomic="false">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={removeToast} />
      ))}
    </div>
  );
};

export default ToastContainer;
