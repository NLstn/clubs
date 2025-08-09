import { FC, ReactNode } from 'react';
import './Modal.css';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  maxWidth?: string;
  showCloseButton?: boolean;
  className?: string;
}

export interface ModalBodyProps {
  children: ReactNode;
  className?: string;
}

export interface ModalActionsProps {
  children: ReactNode;
  className?: string;
}

export interface ModalErrorProps {
  error: string | null;
  onClose?: () => void;
}

const ModalBody: FC<ModalBodyProps> = ({ children, className = "" }) => (
  <div className={`modal-body ${className}`}>
    {children}
  </div>
);

const ModalActions: FC<ModalActionsProps> = ({ children, className = "" }) => (
  <div className={`modal-actions ${className}`}>
    {children}
  </div>
);

const ModalError: FC<ModalErrorProps> = ({ error, onClose }) => {
  if (!error) return null;
  
  return (
    <div className="modal-error-message">
      <span className="error-icon">⚠️</span>
      <span>{error}</span>
      {onClose && (
        <button 
          onClick={onClose} 
          className="error-close-button"
          aria-label="Close error"
        >
          ✕
        </button>
      )}
    </div>
  );
};

const LoadingSpinner: FC = () => (
  <span className="modal-loading-spinner"></span>
);

interface ModalComponent extends FC<ModalProps> {
  Body: typeof ModalBody;
  Actions: typeof ModalActions;
  Error: typeof ModalError;
  LoadingSpinner: typeof LoadingSpinner;
}

const ModalMain: FC<ModalProps> = ({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  maxWidth = "550px",
  showCloseButton = true,
  className = ""
}) => {
  if (!isOpen) return null;

  return (
    <div className="modal" onClick={onClose}>
      <div 
        className={`modal-content enhanced-modal ${className}`} 
        onClick={(e) => e.stopPropagation()}
        style={{ maxWidth }}
      >
        <div className="modal-header">
          <h2>{title}</h2>
          {showCloseButton && (
            <button 
              onClick={onClose} 
              className="modal-close-button"
              aria-label="Close modal"
            >
              ✕
            </button>
          )}
        </div>
        {children}
      </div>
    </div>
  );
};

// Compound component pattern
const Modal = ModalMain as ModalComponent;
Modal.Body = ModalBody;
Modal.Actions = ModalActions;
Modal.Error = ModalError;
Modal.LoadingSpinner = LoadingSpinner;

export default Modal;
