import React from 'react';
import { usePipelineContext } from '../../contexts/PipelineContext';

const WizardNavigation = ({ steps, currentStep }) => {
  const { nextStep, prevStep, canProceedToNextStep, isCompleted } = usePipelineContext();

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === steps.length - 1;

  const containerStyle = {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '20px 0',
    borderTop: '2px solid #f0f0f0',
    marginTop: '20px',
  };

  const buttonStyle = (disabled = false) => ({
    padding: '12px 24px',
    borderRadius: '8px',
    border: 'none',
    fontSize: '16px',
    fontWeight: 'bold',
    cursor: disabled ? 'not-allowed' : 'pointer',
    transition: 'all 0.3s ease',
    opacity: disabled ? 0.6 : 1,
  });

  const prevButtonStyle = {
    ...buttonStyle(isFirstStep),
    backgroundColor: isFirstStep ? '#bdc3c7' : '#95a5a6',
    color: 'white',
  };

  const nextButtonStyle = {
    ...buttonStyle(!canProceedToNextStep()),
    backgroundColor: !canProceedToNextStep() ? '#bdc3c7' : '#3498db',
    color: 'white',
  };

  const completeButtonStyle = {
    ...buttonStyle(false),
    backgroundColor: '#27ae60',
    color: 'white',
  };

  const stepInfoStyle = {
    fontSize: '14px',
    color: '#7f8c8d',
    fontWeight: 'bold',
  };

  const handleNext = () => {
    if (canProceedToNextStep()) {
      nextStep();
    }
  };

  const handlePrev = () => {
    if (!isFirstStep) {
      prevStep();
    }
  };

  const handleComplete = () => {
    // Логика завершения wizard
    console.log('Wizard completed!');
  };

  return (
    <div style={containerStyle}>
      <button
        style={prevButtonStyle}
        onClick={handlePrev}
        disabled={isFirstStep}
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
        >
          Завершить
        </button>
      ) : (
        <button
          style={nextButtonStyle}
          onClick={handleNext}
          disabled={!canProceedToNextStep()}
        >
          Далее →
        </button>
      )}
    </div>
  );
};

export default WizardNavigation;
