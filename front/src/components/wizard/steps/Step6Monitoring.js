import React, { useState, useEffect } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';
import { usePipeline } from '../../../hooks/usePipeline';
import LoadingSpinner from '../../common/LoadingSpinner';

const Step6Monitoring = () => {
  const { wizardData, updateWizardData, completeWizard } = usePipelineContext();
  const { currentPipeline, loading, error, executePipeline, deletePipeline } = usePipeline();
  
  const [executionStatus, setExecutionStatus] = useState('ready');
  const [executionLogs, setExecutionLogs] = useState([]);

  useEffect(() => {
    if (wizardData.pipeline) {
      setExecutionStatus('ready');
    }
  }, [wizardData.pipeline]);

  const handleExecutePipeline = async () => {
    if (!currentPipeline?.id) {
      console.error('Pipeline ID not found');
      return;
    }

    setExecutionStatus('running');
    setExecutionLogs(prev => [...prev, { timestamp: new Date(), message: 'Запуск пайплайна...', type: 'info' }]);

    const result = await executePipeline(currentPipeline.id);
    
    if (result) {
      setExecutionStatus('completed');
      setExecutionLogs(prev => [...prev, { timestamp: new Date(), message: 'Пайплайн успешно выполнен', type: 'success' }]);
    } else {
      setExecutionStatus('failed');
      setExecutionLogs(prev => [...prev, { timestamp: new Date(), message: 'Ошибка выполнения пайплайна', type: 'error' }]);
    }
  };

  const handleDeletePipeline = async () => {
    if (!currentPipeline?.id) {
      console.error('Pipeline ID not found');
      return;
    }

    const result = await deletePipeline(currentPipeline.id);
    if (result) {
      setExecutionStatus('deleted');
      setExecutionLogs(prev => [...prev, { timestamp: new Date(), message: 'Пайплайн удален', type: 'info' }]);
    }
  };

  const handleCompleteWizard = () => {
    completeWizard();
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'ready': return '#3498db';
      case 'running': return '#f39c12';
      case 'completed': return '#27ae60';
      case 'failed': return '#e74c3c';
      case 'deleted': return '#95a5a6';
      default: return '#7f8c8d';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'ready': return '⏳';
      case 'running': return '🔄';
      case 'completed': return '✅';
      case 'failed': return '❌';
      case 'deleted': return '🗑️';
      default: return '❓';
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

  const buttonStyle = (color = '#3498db') => ({
    padding: '12px 24px',
    backgroundColor: color,
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontSize: '16px',
    fontWeight: 'bold',
    marginRight: '12px',
    marginBottom: '12px',
  });

  const statusStyle = {
    padding: '16px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    border: '1px solid #e0e0e0',
    marginBottom: '20px',
    textAlign: 'center',
  };

  const statusTextStyle = {
    fontSize: '18px',
    fontWeight: 'bold',
    color: getStatusColor(executionStatus),
    marginBottom: '8px',
  };

  const logsStyle = {
    padding: '16px',
    backgroundColor: '#2c3e50',
    borderRadius: '8px',
    color: '#ecf0f1',
    fontFamily: 'monospace',
    fontSize: '14px',
    maxHeight: '300px',
    overflowY: 'auto',
  };

  const logEntryStyle = (type) => ({
    marginBottom: '4px',
    color: type === 'error' ? '#e74c3c' : type === 'success' ? '#27ae60' : '#3498db',
  });

  const summaryStyle = {
    padding: '20px',
    backgroundColor: '#f8f9fa',
    borderRadius: '12px',
    border: '1px solid #e0e0e0',
    marginBottom: '20px',
  };

  const summaryItemStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '8px',
    fontSize: '14px',
  };

  const summaryLabelStyle = {
    fontWeight: 'bold',
    color: '#2c3e50',
  };

  const summaryValueStyle = {
    color: '#7f8c8d',
  };

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Запуск и мониторинг пайплайна</h2>
        
        <div style={summaryStyle}>
          <h3 style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '16px' }}>
            📋 Сводка пайплайна
          </h3>
          
          <div style={summaryItemStyle}>
            <span style={summaryLabelStyle}>Источник:</span>
            <span style={summaryValueStyle}>
              {wizardData.source?.type === 'file' ? `Файл: ${wizardData.source?.file?.name}` : wizardData.source?.type}
            </span>
          </div>
          
          <div style={summaryItemStyle}>
            <span style={summaryLabelStyle}>Целевая система:</span>
            <span style={summaryValueStyle}>
              {wizardData.target?.system?.toUpperCase()}
            </span>
          </div>
          
          <div style={summaryItemStyle}>
            <span style={summaryLabelStyle}>Трансформаций:</span>
            <span style={summaryValueStyle}>
              {wizardData.etlConfig?.transformations?.length || 0}
            </span>
          </div>
          
          <div style={summaryItemStyle}>
            <span style={summaryLabelStyle}>Расписание:</span>
            <span style={summaryValueStyle}>
              {wizardData.etlConfig?.schedule?.type === 'manual' ? 'Ручной запуск' : 
               wizardData.etlConfig?.schedule?.type === 'cron' ? `Cron: ${wizardData.etlConfig?.schedule?.cron}` :
               `Интервал: ${wizardData.etlConfig?.schedule?.interval}с`}
            </span>
          </div>
        </div>

        <div style={statusStyle}>
          <div style={statusTextStyle}>
            {getStatusIcon(executionStatus)} Статус: {executionStatus.toUpperCase()}
          </div>
          <div style={{ fontSize: '14px', color: '#7f8c8d' }}>
            {executionStatus === 'ready' && 'Пайплайн готов к выполнению'}
            {executionStatus === 'running' && 'Пайплайн выполняется...'}
            {executionStatus === 'completed' && 'Пайплайн успешно завершен'}
            {executionStatus === 'failed' && 'Ошибка выполнения пайплайна'}
            {executionStatus === 'deleted' && 'Пайплайн удален'}
          </div>
        </div>

        <div style={{ marginBottom: '20px' }}>
          <button
            style={buttonStyle('#27ae60')}
            onClick={handleExecutePipeline}
            disabled={loading || executionStatus === 'running' || executionStatus === 'deleted'}
          >
            🚀 Запустить пайплайн
          </button>

          <button
            style={buttonStyle('#e74c3c')}
            onClick={handleDeletePipeline}
            disabled={loading || executionStatus === 'running'}
          >
            🗑️ Удалить пайплайн
          </button>

          <button
            style={buttonStyle('#9b59b6')}
            onClick={handleCompleteWizard}
          >
            ✅ Завершить создание
          </button>
        </div>

        {executionLogs.length > 0 && (
          <div>
            <h3 style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '12px' }}>
              📝 Логи выполнения
            </h3>
            <div style={logsStyle}>
              {executionLogs.map((log, index) => (
                <div key={index} style={logEntryStyle(log.type)}>
                  [{log.timestamp.toLocaleTimeString()}] {log.message}
                </div>
              ))}
            </div>
          </div>
        )}

        {loading && <LoadingSpinner message="Выполнение операции..." />}
        {error && <div style={{ color: '#e74c3c', fontSize: '14px', marginTop: '8px' }}>{error}</div>}
      </div>
    </div>
  );
};

export default Step6Monitoring;
