import React from 'react';
import { usePipelineContext } from '../../contexts/PipelineContext';
import WizardNavigation from './WizardNavigation';
import ProgressIndicator from './ProgressIndicator';
import Step1DataSource from './steps/Step1DataSource';
import Step2Analysis from './steps/Step2Analysis';
import Step3TargetSystem from './steps/Step3TargetSystem';
import Step4ETLConfig from './steps/Step4ETLConfig';
import Step5Visualization from './steps/Step5Visualization';
import Step6Monitoring from './steps/Step6Monitoring';

const PipelineWizard = () => {
  const { currentStep, wizardData } = usePipelineContext();

  const steps = [
    { id: 0, title: 'Источник данных', component: Step1DataSource },
    { id: 1, title: 'Анализ', component: Step2Analysis },
    { id: 2, title: 'Целевая система', component: Step3TargetSystem },
    { id: 3, title: 'Конфигурация ETL', component: Step4ETLConfig },
    { id: 4, title: 'Визуализация', component: Step5Visualization },
    { id: 5, title: 'Мониторинг', component: Step6Monitoring },
  ];

  const CurrentStepComponent = steps[currentStep]?.component;

  const containerStyle = {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '20px',
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
  };

  const headerStyle = {
    textAlign: 'center',
    marginBottom: '30px',
    paddingBottom: '20px',
    borderBottom: '2px solid #f0f0f0',
  };

  const titleStyle = {
    fontSize: '28px',
    fontWeight: 'bold',
    color: '#2c3e50',
    marginBottom: '8px',
  };

  const subtitleStyle = {
    fontSize: '16px',
    color: '#7f8c8d',
  };

  const stepContentStyle = {
    minHeight: '400px',
    padding: '20px 0',
  };

  return (
    <div style={containerStyle}>
      <div style={headerStyle}>
        <h1 style={titleStyle}>Интеллектуальный цифровой инженер данных</h1>
        <p style={subtitleStyle}>
          Создание ETL-пайплайна: {steps[currentStep]?.title}
        </p>
      </div>

      <ProgressIndicator steps={steps} currentStep={currentStep} />

      <div style={stepContentStyle}>
        {CurrentStepComponent && <CurrentStepComponent />}
      </div>

      <WizardNavigation steps={steps} currentStep={currentStep} />
    </div>
  );
};

export default PipelineWizard;
