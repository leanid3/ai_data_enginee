import config from '../services/config';

// Валидация файлов
export const validateFile = (file) => {
  const errors = [];

  if (!file) {
    errors.push('Файл не выбран');
    return { isValid: false, errors };
  }

  // Проверка размера файла
  if (file.size > config.app.maxFileSize) {
    errors.push(`Размер файла превышает ${config.app.maxFileSize / (1024 * 1024)}MB`);
  }

  // Проверка типа файла
  const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
  if (!config.app.supportedFileTypes.includes(fileExtension)) {
    errors.push(`Неподдерживаемый тип файла. Поддерживаются: ${config.app.supportedFileTypes.join(', ')}`);
  }

  return {
    isValid: errors.length === 0,
    errors
  };
};

// Валидация данных пайплайна
export const validatePipelineData = (pipelineData) => {
  const errors = [];

  if (!pipelineData.source) {
    errors.push('Не указан источник данных');
  }

  if (!pipelineData.target) {
    errors.push('Не указана целевая система');
  }

  if (!pipelineData.transformations || pipelineData.transformations.length === 0) {
    errors.push('Не указаны трансформации данных');
  }

  return {
    isValid: errors.length === 0,
    errors
  };
};

// Валидация конфигурации ETL
export const validateETLConfig = (config) => {
  const errors = [];

  if (!config.schedule) {
    errors.push('Не указано расписание выполнения');
  }

  if (config.schedule && !config.schedule.cron && !config.schedule.interval) {
    errors.push('Не указан формат расписания (cron или интервал)');
  }

  return {
    isValid: errors.length === 0,
    errors
  };
};

// Валидация URL
export const isValidUrl = (string) => {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
};

// Валидация email
export const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};
