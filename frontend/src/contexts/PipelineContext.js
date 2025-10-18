import React, { createContext, useContext, useReducer, useCallback, useMemo } from 'react';

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
        currentStep: Math.min(state.currentStep + 1, 5),
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
      const newErrors = { ...state.errors };
      if (action.payload.message) {
        newErrors[action.payload.field] = action.payload.message;
      } else {
        delete newErrors[action.payload.field];
      }
      return {
        ...state,
        errors: newErrors,
      };

    case ACTIONS.CLEAR_ERROR:
      const filteredErrors = { ...state.errors };
      delete filteredErrors[action.payload];
      return {
        ...state,
        errors: filteredErrors,
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

  // Валидация шага (не зависит от state в зависимостях)
  const validateStep = useCallback((step, wizardData) => {
    switch (step) {
      case 0: // Источник данных
        if (!wizardData.source) {
          return { isValid: false, errorField: 'source', errorMessage: 'Не выбран источник данных' };
        }
        break;
      case 1: // Анализ
        if (!wizardData.analysis) {
          return { isValid: false, errorField: 'analysis', errorMessage: 'Анализ не выполнен' };
        }
        break;
      case 2: // Целевая система
        if (!wizardData.target) {
          return { isValid: false, errorField: 'target', errorMessage: 'Не выбрана целевая система' };
        }
        break;
      case 3: // Конфигурация ETL
        if (!wizardData.etlConfig) {
          return { isValid: false, errorField: 'etlConfig', errorMessage: 'Не настроена конфигурация ETL' };
        }
        break;
      case 4: // Визуализация
        if (!wizardData.pipeline) {
          return { isValid: false, errorField: 'pipeline', errorMessage: 'Пайплайн не сгенерирован' };
        }
        break;
      case 5: // Мониторинг
        return { isValid: true };
      default:
        return { isValid: false };
    }
    
    return { isValid: true };
  }, []);

  // Мемоизированное значение canProceedToNextStep
  const canProceedToNextStep = useMemo(() => {
    const validation = validateStep(state.currentStep, state.wizardData);
    
    if (!validation.isValid) {
      // Если есть ошибка и она еще не установлена - устанавливаем
      if (validation.errorField && !state.errors[validation.errorField]) {
        setError(validation.errorField, validation.errorMessage);
      }
      return false;
    }
    
    // Если валидация прошла успешно - очищаем ошибки для текущего шага
    const currentStepErrors = {
      0: ['source'],
      1: ['analysis'],
      2: ['target'],
      3: ['etlConfig'],
      4: ['pipeline'],
      5: []
    };
    
    const fieldsToClear = currentStepErrors[state.currentStep] || [];
    fieldsToClear.forEach(field => {
      if (state.errors[field]) {
        clearError(field);
      }
    });
    
    return true;
  }, [state.currentStep, state.wizardData, state.errors, validateStep, setError, clearError]);

  // Мемоизированное значение контекста
  const contextValue = useMemo(() => ({
    // State
    currentStep: state.currentStep,
    wizardData: state.wizardData,
    isCompleted: state.isCompleted,
    errors: state.errors,
    
    // Actions
    setStep,
    nextStep,
    prevStep,
    setWizardData,
    updateWizardData,
    setError,
    clearError,
    resetWizard,
    completeWizard,
    
    // Validation
    canProceedToNextStep,
    
    // Helper functions
    validateStep: (step, data) => validateStep(step, data),
  }), [
    state.currentStep,
    state.wizardData,
    state.isCompleted,
    state.errors,
    setStep,
    nextStep,
    prevStep,
    setWizardData,
    updateWizardData,
    setError,
    clearError,
    resetWizard,
    completeWizard,
    canProceedToNextStep,
    validateStep
  ]);

  return (
    <PipelineContext.Provider value={contextValue}>
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

export default PipelineContext;