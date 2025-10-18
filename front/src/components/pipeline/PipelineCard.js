import React, { useState } from 'react';
import { usePipeline } from '../../hooks/usePipeline';

const PipelineCard = ({ pipeline }) => {
  const { executePipeline, deletePipeline, loading } = usePipeline();
  const [showDetails, setShowDetails] = useState(false);

  const handleExecute = async () => {
    await executePipeline(pipeline.id);
  };

  const handleDelete = async () => {
    if (window.confirm('Вы уверены, что хотите удалить этот пайплайн?')) {
      await deletePipeline(pipeline.id);
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'completed': return '#27ae60';
      case 'running': return '#f39c12';
      case 'failed': return '#e74c3c';
      case 'pending': return '#3498db';
      default: return '#95a5a6';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'completed': return '✅';
      case 'running': return '🔄';
      case 'failed': return '❌';
      case 'pending': return '⏳';
      default: return '❓';
    }
  };

  const cardStyle = {
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    border: '1px solid #e0e0e0',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    overflow: 'hidden',
    transition: 'all 0.3s ease',
  };

  const headerStyle = {
    padding: '16px',
    borderBottom: '1px solid #f0f0f0',
  };

  const titleStyle = {
    fontSize: '16px',
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: '8px',
  };

  const statusStyle = (status) => ({
    display: 'inline-flex',
    alignItems: 'center',
    gap: '4px',
    padding: '4px 8px',
    borderRadius: '12px',
    fontSize: '12px',
    fontWeight: 'bold',
    backgroundColor: getStatusColor(status) + '20',
    color: getStatusColor(status),
  });

  const contentStyle = {
    padding: '16px',
  };

  const infoRowStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '8px',
    fontSize: '14px',
  };

  const labelStyle = {
    color: '#7f8c8d',
    fontWeight: 'bold',
  };

  const valueStyle = {
    color: '#2c3e50',
  };

  const actionsStyle = {
    display: 'flex',
    gap: '8px',
    marginTop: '16px',
    paddingTop: '16px',
    borderTop: '1px solid #f0f0f0',
  };

  const buttonStyle = (color = '#3498db') => ({
    flex: 1,
    padding: '8px 12px',
    backgroundColor: color,
    color: 'white',
    border: 'none',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '12px',
    fontWeight: 'bold',
    transition: 'all 0.3s ease',
  });

  const detailsStyle = {
    padding: '16px',
    backgroundColor: '#f8f9fa',
    borderTop: '1px solid #e0e0e0',
  };

  const detailsTitleStyle = {
    fontSize: '14px',
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: '12px',
  };

  const detailsGridStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
    gap: '12px',
  };

  const detailItemStyle = {
    padding: '8px',
    backgroundColor: '#ffffff',
    borderRadius: '6px',
    fontSize: '12px',
  };

  const detailLabelStyle = {
    fontWeight: 'bold',
    color: '#7f8c8d',
    marginBottom: '2px',
  };

  const detailValueStyle = {
    color: '#2c3e50',
  };

  return (
    <div style={cardStyle}>
      <div style={headerStyle}>
        <div style={titleStyle}>
          Пайплайн {pipeline.id}
        </div>
        <div style={statusStyle(pipeline.status)}>
          {getStatusIcon(pipeline.status)} {pipeline.status}
        </div>
      </div>

      <div style={contentStyle}>
        <div style={infoRowStyle}>
          <span style={labelStyle}>Источник:</span>
          <span style={valueStyle}>
            {pipeline.source?.type === 'file' ? 'Файл' : pipeline.source?.type || 'Неизвестно'}
          </span>
        </div>
        
        <div style={infoRowStyle}>
          <span style={labelStyle}>Назначение:</span>
          <span style={valueStyle}>
            {pipeline.target?.system?.toUpperCase() || 'Неизвестно'}
          </span>
        </div>
        
        <div style={infoRowStyle}>
          <span style={labelStyle}>Создан:</span>
          <span style={valueStyle}>
            {pipeline.created_at ? new Date(pipeline.created_at).toLocaleDateString() : 'Неизвестно'}
          </span>
        </div>

        {pipeline.updated_at && (
          <div style={infoRowStyle}>
            <span style={labelStyle}>Обновлен:</span>
            <span style={valueStyle}>
              {new Date(pipeline.updated_at).toLocaleDateString()}
            </span>
          </div>
        )}

        <div style={actionsStyle}>
          <button
            style={buttonStyle('#27ae60')}
            onClick={handleExecute}
            disabled={loading || pipeline.status === 'running'}
          >
            🚀 Запустить
          </button>
          
          <button
            style={buttonStyle('#3498db')}
            onClick={() => setShowDetails(!showDetails)}
          >
            {showDetails ? '📋 Скрыть' : '📋 Детали'}
          </button>
          
          <button
            style={buttonStyle('#e74c3c')}
            onClick={handleDelete}
            disabled={loading || pipeline.status === 'running'}
          >
            🗑️ Удалить
          </button>
        </div>
      </div>

      {showDetails && (
        <div style={detailsStyle}>
          <div style={detailsTitleStyle}>Детальная информация</div>
          <div style={detailsGridStyle}>
            <div style={detailItemStyle}>
              <div style={detailLabelStyle}>ID пайплайна</div>
              <div style={detailValueStyle}>{pipeline.id}</div>
            </div>
            
            <div style={detailItemStyle}>
              <div style={detailLabelStyle}>Статус</div>
              <div style={detailValueStyle}>{pipeline.status}</div>
            </div>
            
            <div style={detailItemStyle}>
              <div style={detailLabelStyle}>Трансформаций</div>
              <div style={detailValueStyle}>
                {pipeline.transformations?.length || 0}
              </div>
            </div>
            
            <div style={detailItemStyle}>
              <div style={detailLabelStyle}>Расписание</div>
              <div style={detailValueStyle}>
                {pipeline.schedule?.type || 'Ручной запуск'}
              </div>
            </div>
            
            {pipeline.source?.file?.name && (
              <div style={detailItemStyle}>
                <div style={detailLabelStyle}>Файл</div>
                <div style={detailValueStyle}>{pipeline.source.file.name}</div>
              </div>
            )}
            
            {pipeline.target?.config?.database && (
              <div style={detailItemStyle}>
                <div style={detailLabelStyle}>База данных</div>
                <div style={detailValueStyle}>{pipeline.target.config.database}</div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default PipelineCard;
