import React from 'react';
import { usePipelineContext } from '../../contexts/PipelineContext';

const WizardNavigation = ({ steps, currentStep }) => {
  const { 
    nextStep, 
    prevStep, 
    canProceedToNextStep, 
    isCompleted 
  } = usePipelineContext();

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === steps.length - 1;

  const containerStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '20px 0',
    borderTop: '2px solid #f0f0f0',
    marginTop: '20px',
    backgroundColor: '#fafafa',
    borderRadius: '8px',
  };

  const buttonStyle = {
    padding: '12px 24px',
    borderRadius: '8px',
    border: 'none',
    fontSize: '16px',
    fontWeight: 'bold',
    cursor: 'pointer',
    transition: 'all 0.3s ease',
    minWidth: '120px',
  };

  const prevButtonStyle = {
    ...buttonStyle,
    backgroundColor: isFirstStep ? '#bdc3c7' : '#95a5a6',
    color: 'white',
    cursor: isFirstStep ? 'not-allowed' : 'pointer',
    opacity: isFirstStep ? 0.6 : 1,
  };

  const nextButtonStyle = {
    ...buttonStyle,
    backgroundColor: canProceedToNextStep ? '#3498db' : '#bdc3c7',
    color: 'white',
    cursor: canProceedToNextStep ? 'pointer' : 'not-allowed',
    opacity: canProceedToNextStep ? 1 : 0.6,
  };

  const completeButtonStyle = {
    ...buttonStyle,
    backgroundColor: '#27ae60',
    color: 'white',
  };

  const stepInfoStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
    fontWeight: 'bold',
    padding: '8px 16px',
    backgroundColor: 'white',
    borderRadius: '20px',
    border: '1px solid #e0e0e0',
  };

  const handleNext = () => {
    if (canProceedToNextStep) {
      nextStep();
    }
  };

  const handlePrev = () => {
    if (!isFirstStep) {
      prevStep();
    }
  };

  const handleComplete = () => {
    console.log('Wizard completed!');
    // Здесь можно добавить логику завершения
  };

  return (
    <div style={containerStyle}>
      <button
        style={prevButtonStyle}
        onClick={handlePrev}
        disabled={isFirstStep}
        type="button"
      >
        ← Назад
      </button>

      <div style={stepInfoStyle}>
        Шаг {currentStep + 1} из {steps.length}
      </div>

      {isLastStep ? (
        <button
          style={completeButtonStyle}
          onClick={handleComplete}
          type="button"
        >
          Завершить
        </button>
      ) : (
        <button
          style={nextButtonStyle}
          onClick={handleNext}
          disabled={!canProceedToNextStep}
          type="button"
        >
          Далее →
        </button>
      )}
    </div>
  );
};

export default WizardNavigation;