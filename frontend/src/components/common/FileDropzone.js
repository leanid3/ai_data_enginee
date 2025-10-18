import React, { useCallback, useState } from 'react';
import { validateFile } from '../../utils/validators';

const FileDropzone = ({ onFileSelect, acceptedTypes = ['.csv', '.json', '.xml'], maxSize = 10 * 1024 * 1024 }) => {
  const [isDragOver, setIsDragOver] = useState(false);
  const [error, setError] = useState(null);

  const handleDragOver = useCallback((e) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    setIsDragOver(false);
    setError(null);

    const files = Array.from(e.dataTransfer.files);
    if (files.length > 1) {
      setError('Можно загрузить только один файл');
      return;
    }

    const file = files[0];
    const validation = validateFile(file);
    
    if (!validation.isValid) {
      setError(validation.errors.join(', '));
      return;
    }

    onFileSelect(file);
  }, [onFileSelect]);

  const handleFileInput = useCallback((e) => {
    const file = e.target.files[0];
    if (!file) return;

    setError(null);
    const validation = validateFile(file);
    
    if (!validation.isValid) {
      setError(validation.errors.join(', '));
      return;
    }

    onFileSelect(file);
  }, [onFileSelect]);

  const dropzoneStyle = {
    border: `2px dashed ${isDragOver ? '#3498db' : '#bdc3c7'}`,
    borderRadius: '8px',
    padding: '40px 20px',
    textAlign: 'center',
    backgroundColor: isDragOver ? '#f8f9fa' : '#ffffff',
    cursor: 'pointer',
    transition: 'all 0.3s ease',
    position: 'relative',
  };

  const iconStyle = {
    fontSize: '48px',
    color: isDragOver ? '#3498db' : '#bdc3c7',
    marginBottom: '16px',
  };

  const textStyle = {
    fontSize: '16px',
    color: '#2c3e50',
    marginBottom: '8px',
  };

  const subtextStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
    marginBottom: '16px',
  };

  const buttonStyle = {
    padding: '10px 20px',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
  };

  const errorStyle = {
    color: '#e74c3c',
    fontSize: '14px',
    marginTop: '8px',
  };

  return (
    <div>
      <div
        style={dropzoneStyle}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => document.getElementById('file-input').click()}
      >
        <div style={iconStyle}>📁</div>
        <div style={textStyle}>
          {isDragOver ? 'Отпустите файл здесь' : 'Перетащите файл сюда или нажмите для выбора'}
        </div>
        <div style={subtextStyle}>
          Поддерживаемые форматы: {acceptedTypes.join(', ')}
          <br />
          Максимальный размер: {Math.round(maxSize / (1024 * 1024))}MB
        </div>
        <button
          type="button"
          style={buttonStyle}
          onClick={(e) => {
            e.stopPropagation();
            document.getElementById('file-input').click();
          }}
        >
          Выбрать файл
        </button>
        <input
          id="file-input"
          type="file"
          accept={acceptedTypes.join(',')}
          onChange={handleFileInput}
          style={{ display: 'none' }}
        />
      </div>
      {error && <div style={errorStyle}>{error}</div>}
    </div>
  );
};

export default FileDropzone;
