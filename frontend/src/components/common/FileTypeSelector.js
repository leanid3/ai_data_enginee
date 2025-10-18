import React from 'react';

const FileTypeSelector = ({ selectedType, onTypeSelect, disabled = false }) => {
  const fileTypes = [
    {
      type: 'csv',
      name: 'CSV',
      description: 'Comma Separated Values',
      icon: '📊',
      features: ['Табличные данные', 'Простой формат', 'Широко поддерживается']
    },
    {
      type: 'json',
      name: 'JSON',
      description: 'JavaScript Object Notation',
      icon: '📋',
      features: ['Структурированные данные', 'Вложенные объекты', 'API интеграция']
    },
    {
      type: 'xml',
      name: 'XML',
      description: 'eXtensible Markup Language',
      icon: '📄',
      features: ['Иерархические данные', 'Метаданные', 'Валидация схемы']
    }
  ];

  const containerStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: '16px',
    marginTop: '16px',
  };

  const cardStyle = (isSelected) => ({
    border: `2px solid ${isSelected ? '#3498db' : '#e0e0e0'}`,
    borderRadius: '8px',
    padding: '16px',
    cursor: disabled ? 'not-allowed' : 'pointer',
    backgroundColor: isSelected ? '#f0f8ff' : '#ffffff',
    opacity: disabled ? 0.6 : 1,
    transition: 'all 0.3s ease',
    textAlign: 'center',
  });

  const iconStyle = {
    fontSize: '32px',
    marginBottom: '8px',
  };

  const nameStyle = {
    fontSize: '18px',
    fontWeight: 'bold',
    marginBottom: '4px',
    color: '#2c3e50',
  };

  const descriptionStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
    marginBottom: '8px',
  };

  const featuresStyle = {
    fontSize: '12px',
    color: '#95a5a6',
    textAlign: 'left',
  };

  const featureItemStyle = {
    marginBottom: '2px',
  };

  return (
    <div>
      <h3 style={{ marginBottom: '16px', color: '#2c3e50' }}>
        Выберите тип файла
      </h3>
      <div style={containerStyle}>
        {fileTypes.map((fileType) => (
          <div
            key={fileType.type}
            style={cardStyle(selectedType === fileType.type)}
            onClick={() => !disabled && onTypeSelect(fileType.type)}
          >
            <div style={iconStyle}>{fileType.icon}</div>
            <div style={nameStyle}>{fileType.name}</div>
            <div style={descriptionStyle}>{fileType.description}</div>
            <div style={featuresStyle}>
              {fileType.features.map((feature, index) => (
                <div key={index} style={featureItemStyle}>
                  • {feature}
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default FileTypeSelector;
