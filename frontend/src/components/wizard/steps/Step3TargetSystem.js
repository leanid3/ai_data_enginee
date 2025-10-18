import React, { useState } from 'react';
import { usePipelineContext } from '../../../contexts/PipelineContext';
import TargetSelector from '../../common/TargetSelector';

const Step3TargetSystem = () => {
  const { wizardData, updateWizardData, setError, clearError } = usePipelineContext();
  const [selectedTarget, setSelectedTarget] = useState(wizardData.target?.system || '');

  const handleTargetSelect = (targetSystem) => {
    setSelectedTarget(targetSystem);
    clearError('target');
    updateWizardData({
      target: {
        system: targetSystem,
        config: getDefaultConfig(targetSystem),
      }
    });
  };

  const getDefaultConfig = (system) => {
    const configs = {
      postgresql: {
        host: 'localhost',
        port: 5432,
        database: 'analytics',
        username: 'user',
        password: '',
        ssl: false,
      },
      clickhouse: {
        host: 'localhost',
        port: 9000,
        database: 'analytics',
        username: 'default',
        password: '',
        cluster: '',
      },
      hdfs: {
        namenode: 'localhost:9000',
        path: '/data/analytics',
        replication: 3,
        blockSize: 134217728, // 128MB
      }
    };
    return configs[system] || {};
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

  const configSectionStyle = {
    marginTop: '20px',
    padding: '20px',
    backgroundColor: '#f8f9fa',
    borderRadius: '8px',
    border: '1px solid #e0e0e0',
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

  const checkboxStyle = {
    marginRight: '8px',
  };

  return (
    <div style={containerStyle}>
      <div style={sectionStyle}>
        <h2 style={titleStyle}>Выбор целевой системы</h2>
        
        <p style={{ color: '#7f8c8d', marginBottom: '20px' }}>
          Выберите систему, в которую будут загружены обработанные данные.
        </p>

        <TargetSelector
          selectedTarget={selectedTarget}
          onTargetSelect={handleTargetSelect}
        />

        {selectedTarget && (
          <div style={configSectionStyle}>
            <h3 style={{ fontSize: '16px', fontWeight: 'bold', marginBottom: '16px' }}>
              ⚙️ Конфигурация подключения
            </h3>
            
            {selectedTarget === 'postgresql' && (
              <div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Хост:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.host || 'localhost'}
                    placeholder="localhost"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Порт:</label>
                  <input
                    type="number"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.port || 5432}
                    placeholder="5432"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>База данных:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.database || 'analytics'}
                    placeholder="analytics"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Пользователь:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.username || 'user'}
                    placeholder="user"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Пароль:</label>
                  <input
                    type="password"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.password || ''}
                    placeholder="password"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label>
                    <input
                      type="checkbox"
                      style={checkboxStyle}
                      defaultChecked={wizardData.target?.config?.ssl || false}
                    />
                    Использовать SSL
                  </label>
                </div>
              </div>
            )}

            {selectedTarget === 'clickhouse' && (
              <div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Хост:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.host || 'localhost'}
                    placeholder="localhost"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Порт:</label>
                  <input
                    type="number"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.port || 9000}
                    placeholder="9000"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>База данных:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.database || 'analytics'}
                    placeholder="analytics"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Пользователь:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.username || 'default'}
                    placeholder="default"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Пароль:</label>
                  <input
                    type="password"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.password || ''}
                    placeholder="password"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Кластер (опционально):</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.cluster || ''}
                    placeholder="cluster_name"
                  />
                </div>
              </div>
            )}

            {selectedTarget === 'hdfs' && (
              <div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>NameNode:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.namenode || 'localhost:9000'}
                    placeholder="localhost:9000"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Путь:</label>
                  <input
                    type="text"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.path || '/data/analytics'}
                    placeholder="/data/analytics"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Репликация:</label>
                  <input
                    type="number"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.replication || 3}
                    placeholder="3"
                    min="1"
                    max="10"
                  />
                </div>
                <div style={inputGroupStyle}>
                  <label style={labelStyle}>Размер блока (байты):</label>
                  <input
                    type="number"
                    style={inputStyle}
                    defaultValue={wizardData.target?.config?.blockSize || 134217728}
                    placeholder="134217728"
                  />
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default Step3TargetSystem;
