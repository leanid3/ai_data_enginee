import React from 'react';

const ProgressIndicator = ({ steps, currentStep }) => {
  const containerStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '30px',
    padding: '0 20px',
  };

  const stepStyle = (stepIndex, isCompleted, isCurrent) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    flex: 1,
    position: 'relative',
  });

  const stepNumberStyle = (isCompleted, isCurrent) => ({
    width: '40px',
    height: '40px',
    borderRadius: '50%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontSize: '16px',
    fontWeight: 'bold',
    marginBottom: '8px',
    backgroundColor: isCompleted ? '#27ae60' : isCurrent ? '#3498db' : '#bdc3c7',
    color: 'white',
    transition: 'all 0.3s ease',
  });

  const stepTitleStyle = (isCompleted, isCurrent) => ({
    fontSize: '12px',
    textAlign: 'center',
    color: isCompleted || isCurrent ? '#2c3e50' : '#7f8c8d',
    fontWeight: isCurrent ? 'bold' : 'normal',
    maxWidth: '80px',
    lineHeight: '1.2',
  });

  const connectorStyle = (isCompleted) => ({
    position: 'absolute',
    top: '20px',
    left: '50%',
    width: '100%',
    height: '2px',
    backgroundColor: isCompleted ? '#27ae60' : '#bdc3c7',
    zIndex: -1,
  });

  return (
    <div style={containerStyle}>
      {steps.map((step, index) => {
        const isCompleted = index < currentStep;
        const isCurrent = index === currentStep;
        const isLast = index === steps.length - 1;

        return (
          <div key={step.id} style={stepStyle(index, isCompleted, isCurrent)}>
            <div style={stepNumberStyle(isCompleted, isCurrent)}>
              {isCompleted ? '✓' : index + 1}
            </div>
            <div style={stepTitleStyle(isCompleted, isCurrent)}>
              {step.title}
            </div>
            {!isLast && (
              <div style={connectorStyle(isCompleted)} />
            )}
          </div>
        );
      })}
    </div>
  );
};

export default ProgressIndicator;
