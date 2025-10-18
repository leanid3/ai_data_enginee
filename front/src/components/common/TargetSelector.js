import React from 'react';
import config from '../../services/config';

const TargetSelector = ({ selectedTarget, onTargetSelect, disabled = false }) => {
  const targetSystems = Object.entries(config.targetSystems).map(([key, system]) => ({
    key,
    ...system
  }));

  const containerStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
    gap: '20px',
    marginTop: '16px',
  };

  const cardStyle = (isSelected) => ({
    border: `2px solid ${isSelected ? '#27ae60' : '#e0e0e0'}`,
    borderRadius: '12px',
    padding: '20px',
    cursor: disabled ? 'not-allowed' : 'pointer',
    backgroundColor: isSelected ? '#f0fff4' : '#ffffff',
    opacity: disabled ? 0.6 : 1,
    transition: 'all 0.3s ease',
    textAlign: 'center',
    boxShadow: isSelected ? '0 4px 12px rgba(39, 174, 96, 0.2)' : '0 2px 4px rgba(0,0,0,0.1)',
  });

  const iconStyle = {
    fontSize: '48px',
    marginBottom: '12px',
  };

  const nameStyle = {
    fontSize: '20px',
    fontWeight: 'bold',
    marginBottom: '8px',
    color: '#2c3e50',
  };

  const descriptionStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
    marginBottom: '16px',
    lineHeight: '1.4',
  };

  const featuresContainerStyle = {
    textAlign: 'left',
  };

  const featuresTitleStyle = {
    fontSize: '14px',
    fontWeight: 'bold',
    color: '#34495e',
    marginBottom: '8px',
  };

  const featuresListStyle = {
    listStyle: 'none',
    padding: 0,
    margin: 0,
  };

  const featureItemStyle = {
    fontSize: '12px',
    color: '#7f8c8d',
    marginBottom: '4px',
    paddingLeft: '16px',
    position: 'relative',
  };

  const featureIconStyle = {
    position: 'absolute',
    left: '0',
    color: '#27ae60',
  };

  return (
    <div>
      <h3 style={{ marginBottom: '16px', color: '#2c3e50' }}>
        Выберите целевую систему
      </h3>
      <div style={containerStyle}>
        {targetSystems.map((system) => (
          <div
            key={system.key}
            style={cardStyle(selectedTarget === system.key)}
            onClick={() => !disabled && onTargetSelect(system.key)}
          >
            <div style={iconStyle}>{system.icon}</div>
            <div style={nameStyle}>{system.name}</div>
            <div style={descriptionStyle}>{system.description}</div>
            <div style={featuresContainerStyle}>
              <div style={featuresTitleStyle}>Особенности:</div>
              <ul style={featuresListStyle}>
                {system.features.map((feature, index) => (
                  <li key={index} style={featureItemStyle}>
                    <span style={featureIconStyle}>✓</span>
                    {feature}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TargetSelector;
