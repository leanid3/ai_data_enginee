import React, { useState, useEffect } from 'react';

const Notification = ({ message, type = 'info', duration = 5000, onClose }) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => {
        setIsVisible(false);
        if (onClose) onClose();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [duration, onClose]);

  const getTypeStyles = () => {
    switch (type) {
      case 'success':
        return {
          backgroundColor: '#d4edda',
          borderColor: '#c3e6cb',
          color: '#155724',
        };
      case 'error':
        return {
          backgroundColor: '#f8d7da',
          borderColor: '#f5c6cb',
          color: '#721c24',
        };
      case 'warning':
        return {
          backgroundColor: '#fff3cd',
          borderColor: '#ffeaa7',
          color: '#856404',
        };
      default:
        return {
          backgroundColor: '#d1ecf1',
          borderColor: '#bee5eb',
          color: '#0c5460',
        };
    }
  };

  const getIcon = () => {
    switch (type) {
      case 'success':
        return '✅';
      case 'error':
        return '❌';
      case 'warning':
        return '⚠️';
      default:
        return 'ℹ️';
    }
  };

  const handleClose = () => {
    setIsVisible(false);
    if (onClose) onClose();
  };

  if (!isVisible) return null;

  return (
    <div
      style={{
        position: 'fixed',
        top: '20px',
        right: '20px',
        padding: '12px 16px',
        borderRadius: '8px',
        border: '1px solid',
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
        zIndex: 1000,
        maxWidth: '400px',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        ...getTypeStyles(),
      }}
    >
      <span style={{ fontSize: '16px' }}>{getIcon()}</span>
      <span style={{ flex: 1 }}>{message}</span>
      <button
        onClick={handleClose}
        style={{
          background: 'none',
          border: 'none',
          fontSize: '18px',
          cursor: 'pointer',
          padding: '0',
          color: 'inherit',
          opacity: 0.7,
        }}
      >
        ×
      </button>
    </div>
  );
};

export default Notification;
