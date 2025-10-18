import React from 'react';

const LoadingSpinner = ({ size = 'medium', message = 'Загрузка...' }) => {
  const getSizeClass = () => {
    switch (size) {
      case 'small':
        return { width: '20px', height: '20px' };
      case 'large':
        return { width: '60px', height: '60px' };
      default:
        return { width: '40px', height: '40px' };
    }
  };

  const spinnerStyle = {
    ...getSizeClass(),
    border: '3px solid #f3f3f3',
    borderTop: '3px solid #3498db',
    borderRadius: '50%',
    animation: 'spin 1s linear infinite',
    margin: '0 auto',
  };

  const containerStyle = {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '20px',
  };

  const messageStyle = {
    marginTop: '10px',
    color: '#666',
    fontSize: '14px',
  };

  return (
    <div style={containerStyle}>
      <div style={spinnerStyle}></div>
      {message && <div style={messageStyle}>{message}</div>}
      <style jsx>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
};

export default LoadingSpinner;
