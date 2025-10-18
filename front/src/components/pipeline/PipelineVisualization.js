import React from 'react';

const PipelineVisualization = ({ pipeline, style = {} }) => {
  if (!pipeline) {
    return (
      <div style={{
        padding: '40px',
        textAlign: 'center',
        color: '#7f8c8d',
        ...style
      }}>
        <div style={{ fontSize: '48px', marginBottom: '16px' }}>📊</div>
        <div>Пайплайн не найден</div>
      </div>
    );
  }

  const containerStyle = {
    padding: '20px',
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    border: '1px solid #e0e0e0',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    ...style
  };

  const nodeStyle = (type, isActive = false) => ({
    padding: '16px 20px',
    borderRadius: '12px',
    margin: '12px 0',
    textAlign: 'center',
    fontWeight: 'bold',
    color: 'white',
    backgroundColor: isActive ? '#2ecc71' : 
                     type === 'source' ? '#27ae60' : 
                     type === 'target' ? '#e74c3c' : 
                     type === 'transform' ? '#3498db' : '#95a5a6',
    boxShadow: isActive ? '0 4px 12px rgba(46, 204, 113, 0.3)' : '0 2px 4px rgba(0,0,0,0.1)',
    transition: 'all 0.3s ease',
    position: 'relative',
  });

  const arrowStyle = {
    textAlign: 'center',
    fontSize: '24px',
    color: '#7f8c8d',
    margin: '8px 0',
    position: 'relative',
  };

  const nodeIconStyle = {
    fontSize: '24px',
    marginBottom: '8px',
  };

  const nodeTitleStyle = {
    fontSize: '16px',
    marginBottom: '4px',
  };

  const nodeSubtitleStyle = {
    fontSize: '12px',
    opacity: 0.8,
  };

  const statusBadgeStyle = (status) => ({
    position: 'absolute',
    top: '-8px',
    right: '-8px',
    width: '20px',
    height: '20px',
    borderRadius: '50%',
    backgroundColor: status === 'running' ? '#f39c12' : 
                     status === 'completed' ? '#27ae60' : 
                     status === 'failed' ? '#e74c3c' : '#95a5a6',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontSize: '12px',
    color: 'white',
  });

  const getNodeIcon = (type) => {
    switch (type) {
      case 'source': return '📁';
      case 'target': return '🎯';
      case 'transform': return '🔄';
      default: return '📦';
    }
  };

  const getNodeTitle = (type, data) => {
    switch (type) {
      case 'source':
        return data?.type === 'file' ? 'Файл' : data?.type || 'Источник';
      case 'target':
        return data?.system?.toUpperCase() || 'Назначение';
      case 'transform':
        return data?.name || 'Трансформация';
      default:
        return 'Узел';
    }
  };

  const getNodeSubtitle = (type, data) => {
    switch (type) {
      case 'source':
        return data?.file?.name || data?.type || 'Источник данных';
      case 'target':
        return data?.config?.database || data?.config?.path || 'Целевая система';
      case 'transform':
        return data?.type || 'Обработка данных';
      default:
        return '';
    }
  };

  return (
    <div style={containerStyle}>
      <h3 style={{
        fontSize: '18px',
        fontWeight: 'bold',
        marginBottom: '20px',
        color: '#2c3e50',
        textAlign: 'center'
      }}>
        📊 Визуализация пайплайна
      </h3>

      {/* Источник данных */}
      <div style={nodeStyle('source', pipeline.status === 'running')}>
        {pipeline.status === 'running' && <div style={statusBadgeStyle('running')}>🔄</div>}
        <div style={nodeIconStyle}>{getNodeIcon('source')}</div>
        <div style={nodeTitleStyle}>{getNodeTitle('source', pipeline.source)}</div>
        <div style={nodeSubtitleStyle}>{getNodeSubtitle('source', pipeline.source)}</div>
      </div>

      <div style={arrowStyle}>⬇️</div>

      {/* Трансформации */}
      {pipeline.transformations?.map((transformation, index) => (
        <div key={index}>
          <div style={nodeStyle('transform', pipeline.status === 'running')}>
            {pipeline.status === 'running' && <div style={statusBadgeStyle('running')}>🔄</div>}
            <div style={nodeIconStyle}>{getNodeIcon('transform')}</div>
            <div style={nodeTitleStyle}>{getNodeTitle('transform', transformation)}</div>
            <div style={nodeSubtitleStyle}>{getNodeSubtitle('transform', transformation)}</div>
          </div>
          <div style={arrowStyle}>⬇️</div>
        </div>
      ))}

      {/* Целевая система */}
      <div style={nodeStyle('target', pipeline.status === 'completed')}>
        {pipeline.status === 'completed' && <div style={statusBadgeStyle('completed')}>✅</div>}
        {pipeline.status === 'failed' && <div style={statusBadgeStyle('failed')}>❌</div>}
        <div style={nodeIconStyle}>{getNodeIcon('target')}</div>
        <div style={nodeTitleStyle}>{getNodeTitle('target', pipeline.target)}</div>
        <div style={nodeSubtitleStyle}>{getNodeSubtitle('target', pipeline.target)}</div>
      </div>

      {/* Информация о пайплайне */}
      <div style={{
        marginTop: '20px',
        padding: '16px',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px',
        fontSize: '14px',
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
          <span style={{ fontWeight: 'bold', color: '#2c3e50' }}>ID:</span>
          <span style={{ color: '#7f8c8d' }}>{pipeline.id}</span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
          <span style={{ fontWeight: 'bold', color: '#2c3e50' }}>Статус:</span>
          <span style={{ 
            color: pipeline.status === 'completed' ? '#27ae60' : 
                   pipeline.status === 'failed' ? '#e74c3c' : 
                   pipeline.status === 'running' ? '#f39c12' : '#7f8c8d'
          }}>
            {pipeline.status}
          </span>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
          <span style={{ fontWeight: 'bold', color: '#2c3e50' }}>Создан:</span>
          <span style={{ color: '#7f8c8d' }}>
            {pipeline.created_at ? new Date(pipeline.created_at).toLocaleString() : 'Неизвестно'}
          </span>
        </div>
        {pipeline.updated_at && (
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span style={{ fontWeight: 'bold', color: '#2c3e50' }}>Обновлен:</span>
            <span style={{ color: '#7f8c8d' }}>
              {new Date(pipeline.updated_at).toLocaleString()}
            </span>
          </div>
        )}
      </div>
    </div>
  );
};

export default PipelineVisualization;
