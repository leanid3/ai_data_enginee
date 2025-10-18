import React, { useState, useEffect } from 'react';
import { usePipeline } from '../../hooks/usePipeline';
import PipelineCard from './PipelineCard';
import LoadingSpinner from '../common/LoadingSpinner';

const PipelineList = () => {
  const { pipelines, loading, error, getPipelines } = usePipeline();
  const [filter, setFilter] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    getPipelines();
  }, [getPipelines]);

  const filteredPipelines = pipelines.filter(pipeline => {
    const matchesFilter = filter === 'all' || pipeline.status === filter;
    const matchesSearch = pipeline.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         pipeline.source?.type?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         pipeline.target?.system?.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesFilter && matchesSearch;
  });

  const containerStyle = {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '20px',
  };

  const headerStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '30px',
    flexWrap: 'wrap',
    gap: '16px',
  };

  const titleStyle = {
    fontSize: '28px',
    fontWeight: 'bold',
    color: '#2c3e50',
    margin: 0,
  };

  const controlsStyle = {
    display: 'flex',
    gap: '12px',
    alignItems: 'center',
    flexWrap: 'wrap',
  };

  const inputStyle = {
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontSize: '14px',
    minWidth: '200px',
  };

  const selectStyle = {
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontSize: '14px',
    backgroundColor: 'white',
  };

  const buttonStyle = {
    padding: '8px 16px',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '14px',
  };

  const gridStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))',
    gap: '20px',
  };

  const emptyStateStyle = {
    textAlign: 'center',
    padding: '60px 20px',
    color: '#7f8c8d',
  };

  const emptyIconStyle = {
    fontSize: '64px',
    marginBottom: '16px',
  };

  const statsStyle = {
    display: 'flex',
    gap: '20px',
    marginBottom: '20px',
    flexWrap: 'wrap',
  };

  const statItemStyle = {
    padding: '12px 16px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    textAlign: 'center',
    minWidth: '100px',
  };

  const statValueStyle = {
    fontSize: '24px',
    fontWeight: 'bold',
    color: '#2c3e50',
  };

  const statLabelStyle = {
    fontSize: '12px',
    color: '#7f8c8d',
    marginTop: '4px',
  };

  const getStatusCount = (status) => {
    return pipelines.filter(p => p.status === status).length;
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

  return (
    <div style={containerStyle}>
      <div style={headerStyle}>
        <h1 style={titleStyle}>Управление пайплайнами</h1>
        <div style={controlsStyle}>
          <input
            type="text"
            placeholder="Поиск пайплайнов..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            style={inputStyle}
          />
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            style={selectStyle}
          >
            <option value="all">Все статусы</option>
            <option value="pending">Ожидание</option>
            <option value="running">Выполняется</option>
            <option value="completed">Завершен</option>
            <option value="failed">Ошибка</option>
          </select>
          <button
            style={buttonStyle}
            onClick={() => getPipelines()}
            disabled={loading}
          >
            🔄 Обновить
          </button>
        </div>
      </div>

      {pipelines.length > 0 && (
        <div style={statsStyle}>
          <div style={statItemStyle}>
            <div style={statValueStyle}>{pipelines.length}</div>
            <div style={statLabelStyle}>Всего</div>
          </div>
          <div style={statItemStyle}>
            <div style={{ ...statValueStyle, color: getStatusColor('running') }}>
              {getStatusCount('running')}
            </div>
            <div style={statLabelStyle}>Выполняется</div>
          </div>
          <div style={statItemStyle}>
            <div style={{ ...statValueStyle, color: getStatusColor('completed') }}>
              {getStatusCount('completed')}
            </div>
            <div style={statLabelStyle}>Завершено</div>
          </div>
          <div style={statItemStyle}>
            <div style={{ ...statValueStyle, color: getStatusColor('failed') }}>
              {getStatusCount('failed')}
            </div>
            <div style={statLabelStyle}>Ошибки</div>
          </div>
        </div>
      )}

      {loading && <LoadingSpinner message="Загрузка пайплайнов..." />}

      {error && (
        <div style={{
          padding: '16px',
          backgroundColor: '#f8d7da',
          color: '#721c24',
          borderRadius: '8px',
          marginBottom: '20px',
        }}>
          ❌ Ошибка загрузки пайплайнов: {error}
        </div>
      )}

      {!loading && !error && filteredPipelines.length === 0 && (
        <div style={emptyStateStyle}>
          <div style={emptyIconStyle}>📊</div>
          <h3 style={{ fontSize: '20px', marginBottom: '8px' }}>
            {pipelines.length === 0 ? 'Пайплайны не найдены' : 'Пайплайны не найдены по фильтру'}
          </h3>
          <p style={{ fontSize: '14px', marginBottom: '20px' }}>
            {pipelines.length === 0 
              ? 'Создайте свой первый ETL-пайплайн с помощью мастера'
              : 'Попробуйте изменить фильтры или поисковый запрос'
            }
          </p>
        </div>
      )}

      {!loading && !error && filteredPipelines.length > 0 && (
        <div style={gridStyle}>
          {filteredPipelines.map((pipeline) => (
            <PipelineCard key={pipeline.id} pipeline={pipeline} />
          ))}
        </div>
      )}
    </div>
  );
};

export default PipelineList;
