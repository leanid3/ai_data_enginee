import React, { useState } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';

const Step4ETLConfig = () => {
  const { wizardData, updateWizardData, setError, clearError } = usePipelineContext();
  
  const [scheduleType, setScheduleType] = useState(wizardData.etlConfig?.schedule?.type || 'manual');
  const [transformations, setTransformations] = useState(wizardData.etlConfig?.transformations || []);

  const handleScheduleTypeChange = (type) => {
    setScheduleType(type);
    clearError('etlConfig');
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: {
          type,
          cron: type === 'cron' ? '0 0 * * *' : null,
          interval: type === 'interval' ? 3600 : null,
        },
        transformations,
      }
    });
  };

  const handleCronChange = (cron) => {
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: {
          type: scheduleType,
          cron,
          interval: null,
        },
        transformations,
      }
    });
  };

  const handleIntervalChange = (interval) => {
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: {
          type: scheduleType,
          cron: null,
          interval: parseInt(interval),
        },
        transformations,
      }
    });
  };

  const addTransformation = () => {
    const newTransformation = {
      id: Date.now(),
      type: 'filter',
      name: '',
      config: {},
    };
    const newTransformations = [...transformations, newTransformation];
    setTransformations(newTransformations);
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: wizardData.etlConfig?.schedule,
        transformations: newTransformations,
      }
    });
  };

  const removeTransformation = (id) => {
    const newTransformations = transformations.filter(t => t.id !== id);
    setTransformations(newTransformations);
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: wizardData.etlConfig?.schedule,
        transformations: newTransformations,
      }
    });
  };

  const updateTransformation = (id, field, value) => {
    const newTransformations = transformations.map(t => 
      t.id === id ? { ...t, [field]: value } : t
    );
    setTransformations(newTransformations);
    updateWizardData({
      etlConfig: {
        ...wizardData.etlConfig,
        schedule: wizardData.etlConfig?.schedule,
        transformations: newTransformations,
      }
    });
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

  const inputGroupStyle = {
    marginBottom: '16px',
  };

  const labelStyle = {
    display: 'block',
    fontSize: '14px',
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: '4px',
  };

  const inputStyle = {
    width: '100%',
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontSize: '14px',
  };

  const buttonStyle = {
    padding: '8px 16px',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
    marginRight: '8px',
    marginBottom: '8px',
  };

  const transformationStyle = {
    padding: '16px',
    border: '1px solid #e0e0e0',
    borderRadius: '8px',
    marginBottom: '12px',
    backgroundColor: '#f8f9fa',
  };

  const transformationTypes = [
    { value: 'filter', label: 'Фильтрация', description: 'Отфильтровать данные по условию' },
    { value: 'map', label: 'Преобразование', description: 'Преобразовать значения полей' },
    { value: 'aggregate', label: 'Агрегация', description: 'Группировка и агрегация данных' },
    { value: 'join', label: 'Соединение', description: 'Соединить с другими данными' },
  ];

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Расписание выполнения</h2>
        
        <div style={optionStyle(scheduleType === 'manual')} onClick={() => handleScheduleTypeChange('manual')}>
          <div style={optionTitleStyle}>🔧 Ручной запуск</div>
          <div style={optionDescriptionStyle}>
            Пайплайн будет запускаться только вручную
          </div>
        </div>

        <div style={optionStyle(scheduleType === 'cron')} onClick={() => handleScheduleTypeChange('cron')}>
          <div style={optionTitleStyle}>⏰ По расписанию (Cron)</div>
          <div style={optionDescriptionStyle}>
            Запуск по расписанию в формате Cron
          </div>
        </div>

        <div style={optionStyle(scheduleType === 'interval')} onClick={() => handleScheduleTypeChange('interval')}>
          <div style={optionTitleStyle}>🔄 По интервалу</div>
          <div style={optionDescriptionStyle}>
            Запуск через определенные интервалы времени
          </div>
        </div>

        {scheduleType === 'cron' && (
          <div style={inputGroupStyle}>
            <label style={labelStyle}>Cron выражение:</label>
            <input
              type="text"
              style={inputStyle}
              defaultValue={wizardData.etlConfig?.schedule?.cron || '0 0 * * *'}
              placeholder="0 0 * * *"
              onChange={(e) => handleCronChange(e.target.value)}
            />
            <small style={{ color: '#7f8c8d', fontSize: '12px' }}>
              Формат: минута час день месяц день_недели
            </small>
          </div>
        )}

        {scheduleType === 'interval' && (
          <div style={inputGroupStyle}>
            <label style={labelStyle}>Интервал (секунды):</label>
            <input
              type="number"
              style={inputStyle}
              defaultValue={wizardData.etlConfig?.schedule?.interval || 3600}
              placeholder="3600"
              onChange={(e) => handleIntervalChange(e.target.value)}
            />
            <small style={{ color: '#7f8c8d', fontSize: '12px' }}>
              3600 секунд = 1 час
            </small>
          </div>
        )}
      </div>

      <div style={sectionStyle}>
        <h2 style={titleStyle}>Трансформации данных</h2>
        
        <p style={{ color: '#7f8c8d', marginBottom: '16px' }}>
          Настройте трансформации, которые будут применены к данным.
        </p>

        {transformations.map((transformation, index) => (
          <div key={transformation.id} style={transformationStyle}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
              <h4 style={{ fontSize: '16px', fontWeight: 'bold', margin: 0 }}>
                Трансформация {index + 1}
              </h4>
              <button
                style={{ ...buttonStyle, backgroundColor: '#e74c3c' }}
                onClick={() => removeTransformation(transformation.id)}
              >
                Удалить
              </button>
            </div>
            
            <div style={inputGroupStyle}>
              <label style={labelStyle}>Тип трансформации:</label>
              <select
                style={inputStyle}
                value={transformation.type}
                onChange={(e) => updateTransformation(transformation.id, 'type', e.target.value)}
              >
                {transformationTypes.map(type => (
                  <option key={type.value} value={type.value}>
                    {type.label} - {type.description}
                  </option>
                ))}
              </select>
            </div>

            <div style={inputGroupStyle}>
              <label style={labelStyle}>Название:</label>
              <input
                type="text"
                style={inputStyle}
                value={transformation.name}
                onChange={(e) => updateTransformation(transformation.id, 'name', e.target.value)}
                placeholder="Введите название трансформации"
              />
            </div>
          </div>
        ))}

        <button
          style={buttonStyle}
          onClick={addTransformation}
        >
          ➕ Добавить трансформацию
        </button>
      </div>
    </div>
  );
};

export default Step4ETLConfig;
