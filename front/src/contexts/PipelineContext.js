import React, { createContext, useContext, useReducer, useCallback } from 'react';
import { validatePipelineData } from '../utils/validators';

// Начальное состояние
const initialState = {
  currentStep: 0,
  wizardData: {
    source: null,
    analysis: null,
    target: null,
    etlConfig: null,
    pipeline: null,
    monitoring: null,
  },
  isCompleted: false,
  errors: {},
};

// Типы действий
const ACTIONS = {
  SET_STEP: 'SET_STEP',
  NEXT_STEP: 'NEXT_STEP',
  PREV_STEP: 'PREV_STEP',
  SET_WIZARD_DATA: 'SET_WIZARD_DATA',
  UPDATE_WIZARD_DATA: 'UPDATE_WIZARD_DATA',
  SET_ERROR: 'SET_ERROR',
  CLEAR_ERROR: 'CLEAR_ERROR',
  RESET_WIZARD: 'RESET_WIZARD',
  COMPLETE_WIZARD: 'COMPLETE_WIZARD',
};

// Редьюсер
const pipelineReducer = (state, action) => {
  switch (action.type) {
    case ACTIONS.SET_STEP:
      return {
        ...state,
        currentStep: action.payload,
        errors: {},
      };

    case ACTIONS.NEXT_STEP:
      return {
        ...state,
        currentStep: Math.min(state.currentStep + 1, 5), // Максимум 6 шагов (0-5)
        errors: {},
      };

    case ACTIONS.PREV_STEP:
      return {
        ...state,
        currentStep: Math.max(state.currentStep - 1, 0),
        errors: {},
      };

    case ACTIONS.SET_WIZARD_DATA:
      return {
        ...state,
        wizardData: action.payload,
      };

    case ACTIONS.UPDATE_WIZARD_DATA:
      return {
        ...state,
        wizardData: {
          ...state.wizardData,
          ...action.payload,
        },
        errors: {},
      };

    case ACTIONS.SET_ERROR:
      return {
        ...state,
        errors: {
          ...state.errors,
          [action.payload.field]: action.payload.message,
        },
      };

    case ACTIONS.CLEAR_ERROR:
      return {
        ...state,
        errors: {
          ...state.errors,
          [action.payload]: undefined,
        },
      };

    case ACTIONS.RESET_WIZARD:
      return initialState;

    case ACTIONS.COMPLETE_WIZARD:
      return {
        ...state,
        isCompleted: true,
      };

    default:
      return state;
  }
};

// Создание контекста
const PipelineContext = createContext();

// Провайдер контекста
export const PipelineProvider = ({ children }) => {
  const [state, dispatch] = useReducer(pipelineReducer, initialState);

  const setStep = useCallback((step) => {
    dispatch({ type: ACTIONS.SET_STEP, payload: step });
  }, []);

  const nextStep = useCallback(() => {
    dispatch({ type: ACTIONS.NEXT_STEP });
  }, []);

  const prevStep = useCallback(() => {
    dispatch({ type: ACTIONS.PREV_STEP });
  }, []);

  const setWizardData = useCallback((data) => {
    dispatch({ type: ACTIONS.SET_WIZARD_DATA, payload: data });
  }, []);

  const updateWizardData = useCallback((data) => {
    dispatch({ type: ACTIONS.UPDATE_WIZARD_DATA, payload: data });
  }, []);

  const setError = useCallback((field, message) => {
    dispatch({ type: ACTIONS.SET_ERROR, payload: { field, message } });
  }, []);

  const clearError = useCallback((field) => {
    dispatch({ type: ACTIONS.CLEAR_ERROR, payload: field });
  }, []);

  const resetWizard = useCallback(() => {
    dispatch({ type: ACTIONS.RESET_WIZARD });
  }, []);

  const completeWizard = useCallback(() => {
    dispatch({ type: ACTIONS.COMPLETE_WIZARD });
  }, []);

  const validateCurrentStep = useCallback(() => {
    const { currentStep, wizardData } = state;
    
    switch (currentStep) {
      case 0: // Источник данных
        if (!wizardData.source) {
          setError('source', 'Не выбран источник данных');
          return false;
        }
        break;
      case 1: // Анализ
        if (!wizardData.analysis) {
          setError('analysis', 'Анализ не выполнен');
          return false;
        }
        break;
      case 2: // Целевая система
        if (!wizardData.target) {
          setError('target', 'Не выбрана целевая система');
          return false;
        }
        break;
      case 3: // Конфигурация ETL
        if (!wizardData.etlConfig) {
          setError('etlConfig', 'Не настроена конфигурация ETL');
          return false;
        }
        break;
      case 4: // Визуализация
        if (!wizardData.pipeline) {
          setError('pipeline', 'Пайплайн не сгенерирован');
          return false;
        }
        break;
      case 5: // Мониторинг
        // На последнем шаге валидация не требуется
        break;
      default:
        return false;
    }
    
    return true;
  }, [state, setError]);

  const canProceedToNextStep = useCallback(() => {
    return validateCurrentStep();
  }, [validateCurrentStep]);

  const value = {
    ...state,
    setStep,
    nextStep,
    prevStep,
    setWizardData,
    updateWizardData,
    setError,
    clearError,
    resetWizard,
    completeWizard,
    validateCurrentStep,
    canProceedToNextStep,
  };

  return (
    <PipelineContext.Provider value={value}>
      {children}
    </PipelineContext.Provider>
  );
};

// Хук для использования контекста
export const usePipelineContext = () => {
  const context = useContext(PipelineContext);
  if (!context) {
    throw new Error('usePipelineContext must be used within a PipelineProvider');
  }
  return context;
};
