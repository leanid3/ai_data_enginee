import React, { useState, useEffect } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';
import { useAnalysis } from '../../../hooks/useAnalysis';
import LoadingSpinner from '../../common/LoadingSpinner';

const Step2Analysis = () => {
  const { wizardData, updateWizardData, setError, clearError } = usePipelineContext();
  const { analysisResults, loading, error, startAnalysis, checkAnalysisStatus, isAnalysisCompleted } = useAnalysis();
  
  const [analysisStarted, setAnalysisStarted] = useState(false);
  const [statusMessage, setStatusMessage] = useState('');

  useEffect(() => {
    if (wizardData.analysis) {
      // Если анализ уже выполнен, показываем результаты
      setAnalysisStarted(true);
    }
  }, [wizardData.analysis]);

  const handleStartAnalysis = async () => {
    if (!wizardData.source?.storagePath) {
      setError('analysis', 'Сначала загрузите файл');
      return;
    }

    setAnalysisStarted(true);
    setStatusMessage('Запуск анализа...');
    clearError('analysis');

    const fileId = wizardData.source.storagePath.split('/')[2];
    const result = await startAnalysis(fileId, wizardData.source.storagePath);
    
    if (result) {
      setStatusMessage('✅ Анализ запущен');
    } else {
      setStatusMessage('❌ Ошибка запуска анализа');
    }
  };

  const handleCheckStatus = async () => {
    setStatusMessage('Проверка статуса анализа...');
    const result = await checkAnalysisStatus();
    
    if (result && isAnalysisCompleted()) {
      setStatusMessage('✅ Анализ завершен');
      updateWizardData({ analysis: analysisResults });
    } else if (result) {
      setStatusMessage(`Статус: ${result.status}`);
    } else {
      setStatusMessage('❌ Ошибка проверки статуса');
    }
  };

  const QualityScore = ({ score }) => {
    const percentage = (score * 100).toFixed(1);
    let color = "#4CAF50";
    if (score < 0.7) color = "#f44336";
    else if (score < 0.9) color = "#FF9800";

    return (
      <div style={{ marginBottom: '20px' }}>
        <div style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '8px' }}>
          Качество данных
        </div>
        <div style={{
          width: '100%',
          height: '20px',
          backgroundColor: '#f0f0f0',
          borderRadius: '10px',
          overflow: 'hidden',
        }}>
          <div
            style={{
              width: `${percentage}%`,
              height: '100%',
              backgroundColor: color,
              transition: 'width 0.3s ease',
            }}
          />
        </div>
        <div style={{ textAlign: 'right', marginTop: '4px', fontSize: '14px' }}>
          {percentage}%
        </div>
      </div>
    );
  };

  const DataCharacteristics = ({ data }) => {
    if (!data) return null;

    return (
      <div style={{ marginBottom: '20px' }}>
        <h4 style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '12px' }}>
          📊 Характеристики данных
        </h4>
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))',
          gap: '12px',
        }}>
          <div style={{ textAlign: 'center', padding: '12px', backgroundColor: '#f8f9fa', borderRadius: '8px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#3498db' }}>
              {data.row_count}
            </div>
            <div style={{ fontSize: '12px', color: '#7f8c8d' }}>строк</div>
          </div>
          <div style={{ textAlign: 'center', padding: '12px', backgroundColor: '#f8f9fa', borderRadius: '8px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#3498db' }}>
              {data.column_count}
            </div>
            <div style={{ fontSize: '12px', color: '#7f8c8d' }}>колонок</div>
          </div>
          <div style={{ textAlign: 'center', padding: '12px', backgroundColor: '#f8f9fa', borderRadius: '8px' }}>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#3498db' }}>
              {data.estimated_size}
            </div>
            <div style={{ fontSize: '12px', color: '#7f8c8d' }}>размер</div>
          </div>
        </div>
      </div>
    );
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
    marginRight: '12px',
    marginBottom: '12px',
  };

  const statusStyle = {
    padding: '12px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    marginBottom: '16px',
    fontSize: '14px',
  };

  const resultsStyle = {
    padding: '20px',
    backgroundColor: '#f8f9fa',
    borderRadius: '12px',
    border: '1px solid #e0e0e0',
  };

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Анализ данных</h2>
        
        {!analysisStarted && (
          <div>
            <p style={{ color: '#7f8c8d', marginBottom: '20px' }}>
              Запустите анализ загруженного файла для получения информации о структуре и качестве данных.
            </p>
            <button
              style={buttonStyle}
              onClick={handleStartAnalysis}
              disabled={loading || !wizardData.source?.storagePath}
            >
              🔍 Запустить анализ
            </button>
          </div>
        )}

        {analysisStarted && !isAnalysisCompleted() && (
          <div>
            <div style={statusStyle}>
              {statusMessage || 'Анализ выполняется...'}
            </div>
            <button
              style={buttonStyle}
              onClick={handleCheckStatus}
              disabled={loading}
            >
              📊 Проверить статус
            </button>
          </div>
        )}

        {isAnalysisCompleted() && analysisResults && (
          <div style={resultsStyle}>
            <h3 style={{ fontSize: '18px', fontWeight: 'bold', marginBottom: '16px' }}>
              ✅ Результаты анализа
            </h3>
            
            <QualityScore score={analysisResults.result.data_quality_score} />
            
            <DataCharacteristics data={analysisResults.result.ddl_metadata.data_characteristics} />
            
            <div style={{ marginBottom: '20px' }}>
              <h4 style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '12px' }}>
                💡 Рекомендации
              </h4>
              <ul style={{ paddingLeft: '20px' }}>
                {analysisResults.result.recommendations.map((rec, index) => (
                  <li key={index} style={{ marginBottom: '4px', color: '#2c3e50' }}>
                    {rec}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}

        {loading && <LoadingSpinner message="Выполнение анализа..." />}
        {error && <div style={{ color: '#e74c3c', fontSize: '14px', marginTop: '8px' }}>{error}</div>}
      </div>
    </div>
  );
};

export default Step2Analysis;
