import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { PipelineProvider } from '../contexts/PipelineContext';
import PipelineWizard from '../components/wizard/PipelineWizard';

// Мокаем хуки
jest.mock('../hooks/useFileUpload', () => ({
  useFileUpload: () => ({
    selectedFile: null,
    storagePath: null,
    loading: false,
    error: null,
    handleFileChange: jest.fn(),
    uploadFile: jest.fn(),
    resetFile: jest.fn(),
    getFileId: jest.fn(),
  }),
}));

jest.mock('../hooks/useAnalysis', () => ({
  useAnalysis: () => ({
    analysisResults: null,
    loading: false,
    error: null,
    startAnalysis: jest.fn(),
    checkAnalysisStatus: jest.fn(),
    resetAnalysis: jest.fn(),
    isAnalysisCompleted: jest.fn(() => false),
  }),
}));

jest.mock('../hooks/usePipeline', () => ({
  usePipeline: () => ({
    currentPipeline: null,
    loading: false,
    error: null,
    generatePipeline: jest.fn(),
    getPipeline: jest.fn(),
    executePipeline: jest.fn(),
    deletePipeline: jest.fn(),
    getPipelines: jest.fn(),
    resetPipeline: jest.fn(),
  }),
}));

const renderWithProvider = (component) => {
  return render(
    <PipelineProvider>
      {component}
    </PipelineProvider>
  );
};

describe('PipelineWizard', () => {
  test('рендерится без ошибок', () => {
    renderWithProvider(<PipelineWizard />);
    
    expect(screen.getByText('Интеллектуальный цифровой инженер данных')).toBeInTheDocument();
    expect(screen.getByText('Создание ETL-пайплайна: Источник данных')).toBeInTheDocument();
  });

  test('отображает индикатор прогресса', () => {
    renderWithProvider(<PipelineWizard />);
    
    expect(screen.getByText('Источник данных')).toBeInTheDocument();
    expect(screen.getByText('Анализ')).toBeInTheDocument();
    expect(screen.getByText('Целевая система')).toBeInTheDocument();
    expect(screen.getByText('Конфигурация ETL')).toBeInTheDocument();
    expect(screen.getByText('Визуализация')).toBeInTheDocument();
    expect(screen.getByText('Мониторинг')).toBeInTheDocument();
  });

  test('отображает навигацию', () => {
    renderWithProvider(<PipelineWizard />);
    
    expect(screen.getByText('Шаг 1 из 6')).toBeInTheDocument();
    expect(screen.getByText('← Назад')).toBeInTheDocument();
    expect(screen.getByText('Далее →')).toBeInTheDocument();
  });

  test('кнопка "Назад" отключена на первом шаге', () => {
    renderWithProvider(<PipelineWizard />);
    
    const backButton = screen.getByText('← Назад');
    expect(backButton).toBeDisabled();
  });
});
