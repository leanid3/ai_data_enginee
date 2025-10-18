import React, { useState, useEffect } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';
import { usePipeline } from '../../../hooks/usePipeline';
import LoadingSpinner from '../../common/LoadingSpinner';

const Step5Visualization = () => {
  const { wizardData, updateWizardData, setError, clearError } = usePipelineContext();
  const { currentPipeline, loading, error, generatePipeline } = usePipeline();
  
  const [pipelineGenerated, setPipelineGenerated] = useState(false);

  useEffect(() => {
    if (wizardData.pipeline) {
      setPipelineGenerated(true);
    }
  }, [wizardData.pipeline]);

  const handleGeneratePipeline = async () => {
    if (!wizardData.source || !wizardData.target || !wizardData.etlConfig) {
      setError('pipeline', 'Не все шаги завершены');
      return;
    }

    clearError('pipeline');
    
    const pipelineData = {
      source: wizardData.source,
      target: wizardData.target,
      etlConfig: wizardData.etlConfig,
      analysis: wizardData.analysis,
    };

    const result = await generatePipeline(pipelineData);
    if (result) {
      setPipelineGenerated(true);
      updateWizardData({ pipeline: result });
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

  const buttonStyle = {
    padding: '12px 24px',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontSize: '16px',
    fontWeight: 'bold',
    marginBottom: '20px',
  };

  const pipelineStyle = {
    padding: '20px',
    backgroundColor: '#f8f9fa',
    borderRadius: '12px',
    border: '1px solid #e0e0e0',
  };

  const nodeStyle = (type) => ({
    padding: '12px 16px',
    borderRadius: '8px',
    margin: '8px 0',
    textAlign: 'center',
    fontWeight: 'bold',
    color: 'white',
    backgroundColor: type === 'source' ? '#27ae60' : type === 'target' ? '#e74c3c' : '#3498db',
  });

  const arrowStyle = {
    textAlign: 'center',
    fontSize: '24px',
    color: '#7f8c8d',
    margin: '8px 0',
  };

  const infoStyle = {
    padding: '12px',
    backgroundColor: '#e8f4fd',
    borderRadius: '8px',
    marginBottom: '16px',
    fontSize: '14px',
  };

  const errorStyle = {
    color: '#e74c3c',
    fontSize: '14px',
    marginTop: '8px',
  };

  const renderPipelineVisualization = () => {
    if (!currentPipeline) return null;

    return (
      <div style={pipelineStyle}>
        <h3 style={{ fontSize: '18px', fontWeight: 'bold', marginBottom: '16px' }}>
          📊 Визуализация пайплайна
        </h3>
        
        {/* Источник данных */}
        <div style={nodeStyle('source')}>
          📁 {wizardData.source?.type === 'file' ? 'Файл' : wizardData.source?.type}
          <br />
          <small style={{ fontSize: '12px', opacity: 0.8 }}>
            {wizardData.source?.file?.name || 'Источник данных'}
          </small>
        </div>

        <div style={arrowStyle}>⬇️</div>

        {/* Трансформации */}
        {wizardData.etlConfig?.transformations?.map((transformation, index) => (
          <div key={index}>
            <div style={nodeStyle('transform')}>
              🔄 {transformation.name || `Трансформация ${index + 1}`}
              <br />
              <small style={{ fontSize: '12px', opacity: 0.8 }}>
                {transformation.type}
              </small>
            </div>
            <div style={arrowStyle}>⬇️</div>
          </div>
        ))}

        {/* Целевая система */}
        <div style={nodeStyle('target')}>
          🎯 {wizardData.target?.system?.toUpperCase() || 'Целевая система'}
          <br />
          <small style={{ fontSize: '12px', opacity: 0.8 }}>
            {wizardData.target?.config?.database || wizardData.target?.config?.path || 'Назначение'}
          </small>
        </div>

        <div style={infoStyle}>
          <strong>ID пайплайна:</strong> {currentPipeline.id || 'Генерируется...'}
          <br />
          <strong>Статус:</strong> {currentPipeline.status || 'Создан'}
          <br />
          <strong>Создан:</strong> {new Date().toLocaleString()}
        </div>
      </div>
    );
  };

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Генерация и визуализация пайплайна</h2>
        
        <p style={{ color: '#7f8c8d', marginBottom: '20px' }}>
          Сгенерируйте ETL-пайплайн на основе настроенных параметров и просмотрите его визуализацию.
        </p>

        {!pipelineGenerated && (
          <div>
            <button
              style={buttonStyle}
              onClick={handleGeneratePipeline}
              disabled={loading || !wizardData.source || !wizardData.target || !wizardData.etlConfig}
            >
              🚀 Сгенерировать пайплайн
            </button>
            
            {(!wizardData.source || !wizardData.target || !wizardData.etlConfig) && (
              <div style={infoStyle}>
                ⚠️ Для генерации пайплайна необходимо завершить все предыдущие шаги:
                <ul style={{ marginTop: '8px', paddingLeft: '20px' }}>
                  {!wizardData.source && <li>Источник данных</li>}
                  {!wizardData.target && <li>Целевая система</li>}
                  {!wizardData.etlConfig && <li>Конфигурация ETL</li>}
                </ul>
              </div>
            )}
          </div>
        )}

        {pipelineGenerated && renderPipelineVisualization()}

        {loading && <LoadingSpinner message="Генерация пайплайна..." />}
        {error && <div style={errorStyle}>{error}</div>}
      </div>
    </div>
  );
};

export default Step5Visualization;
