import React, { useState } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';
import { useFileUpload } from '../../../hooks/useFileUpload';
import FileDropzone from '../../common/FileDropzone';
import FileTypeSelector from '../../common/FileTypeSelector';
import LoadingSpinner from '../../common/LoadingSpinner';

const Step1DataSource = () => {
  const { wizardData, updateWizardData, setError, clearError } = usePipelineContext();
  const { selectedFile, storagePath, loading, error, handleFileChange, uploadFile, getFileId } = useFileUpload();
  
  const [dataSourceType, setDataSourceType] = useState(wizardData.source?.type || 'file');
  const [fileType, setFileType] = useState(wizardData.source?.fileType || 'csv');

  const handleDataSourceTypeChange = (type) => {
    setDataSourceType(type);
    clearError('source');
    updateWizardData({
      source: {
        type,
        fileType: type === 'file' ? fileType : null,
        file: type === 'file' ? selectedFile : null,
        storagePath: type === 'file' ? storagePath : null,
      }
    });
  };

  const handleFileTypeChange = (type) => {
    setFileType(type);
    clearError('source');
    updateWizardData({
      source: {
        type: dataSourceType,
        fileType: type,
        file: selectedFile,
        storagePath,
      }
    });
  };

  const handleFileSelect = (file) => {
    handleFileChange(file);
    clearError('source');
    updateWizardData({
      source: {
        type: dataSourceType,
        fileType,
        file,
        storagePath,
      }
    });
  };

  const handleUpload = async () => {
    const result = await uploadFile();
    if (result) {
      updateWizardData({
        source: {
          type: dataSourceType,
          fileType,
          file: selectedFile,
          storagePath: result.storage_path,
        }
      });
    }
  };

  const containerStyle = {
    maxWidth: '800px',
    margin: '0 auto',
  };

  const sectionStyle = {
    marginBottom: '30px',
  };

  const titleStyle = {
    fontSize: '20px',
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: '16px',
  };

  const optionStyle = (isSelected) => ({
    padding: '16px',
    border: `2px solid ${isSelected ? '#3498db' : '#e0e0e0'}`,
    borderRadius: '8px',
    cursor: 'pointer',
    backgroundColor: isSelected ? '#f0f8ff' : '#ffffff',
    marginBottom: '12px',
    transition: 'all 0.3s ease',
  });

  const optionTitleStyle = {
    fontSize: '16px',
    fontWeight: 'bold',
    marginBottom: '4px',
    color: '#2c3e50',
  };

  const optionDescriptionStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
  };

  const buttonStyle = {
    padding: '12px 24px',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontSize: '16px',
    fontWeight: 'bold',
    marginTop: '16px',
  };

  const fileInfoStyle = {
    padding: '12px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    marginTop: '16px',
  };

  const errorStyle = {
    color: '#e74c3c',
    fontSize: '14px',
    marginTop: '8px',
  };

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Выберите источник данных</h2>
        
        <div style={optionStyle(dataSourceType === 'file')} onClick={() => handleDataSourceTypeChange('file')}>
          <div style={optionTitleStyle}>📁 Файл</div>
          <div style={optionDescriptionStyle}>
            Загрузите файл с данными (CSV, JSON, XML)
          </div>
        </div>

        <div style={optionStyle(dataSourceType === 'database')} onClick={() => handleDataSourceTypeChange('database')}>
          <div style={optionTitleStyle}>🗄️ База данных</div>
          <div style={optionDescriptionStyle}>
            Подключение к существующей базе данных
          </div>
        </div>

        <div style={optionStyle(dataSourceType === 'stream')} onClick={() => handleDataSourceTypeChange('stream')}>
          <div style={optionTitleStyle}>🌊 Поток данных</div>
          <div style={optionDescriptionStyle}>
            Обработка данных в реальном времени
          </div>
        </div>
      </div>

      {dataSourceType === 'file' && (
        <>
          <div style={sectionStyle}>
            <h3 style={titleStyle}>Тип файла</h3>
            <FileTypeSelector
              selectedType={fileType}
              onTypeSelect={handleFileTypeChange}
            />
          </div>

          <div style={sectionStyle}>
            <h3 style={titleStyle}>Загрузка файла</h3>
            <FileDropzone
              onFileSelect={handleFileSelect}
              acceptedTypes={['.csv', '.json', '.xml']}
            />
            
            {selectedFile && (
              <div style={fileInfoStyle}>
                <strong>Выбранный файл:</strong> {selectedFile.name}
                <br />
                <strong>Размер:</strong> {(selectedFile.size / 1024).toFixed(2)} KB
                <br />
                <strong>Тип:</strong> {selectedFile.type || 'Неизвестно'}
              </div>
            )}

            {selectedFile && !storagePath && (
              <button
                style={buttonStyle}
                onClick={handleUpload}
                disabled={loading}
              >
                {loading ? 'Загрузка...' : '📤 Загрузить файл'}
              </button>
            )}

            {storagePath && (
              <div style={fileInfoStyle}>
                ✅ Файл успешно загружен
                <br />
                <strong>Путь:</strong> {storagePath}
              </div>
            )}
          </div>
        </>
      )}

      {dataSourceType === 'database' && (
        <div style={sectionStyle}>
          <h3 style={titleStyle}>Подключение к базе данных</h3>
          <p style={{ color: '#7f8c8d' }}>
            Функция подключения к базе данных будет реализована в следующих версиях.
          </p>
        </div>
      )}

      {dataSourceType === 'stream' && (
        <div style={sectionStyle}>
          <h3 style={titleStyle}>Настройка потока данных</h3>
          <p style={{ color: '#7f8c8d' }}>
            Функция обработки потоков данных будет реализована в следующих версиях.
          </p>
        </div>
      )}

      {loading && <LoadingSpinner message="Обработка файла..." />}
      {error && <div style={errorStyle}>{error}</div>}
    </div>
  );
};

export default Step1DataSource;
