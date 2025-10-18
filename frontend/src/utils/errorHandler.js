// Обработка ошибок приложения
export class AppError extends Error {
  constructor(message, code = 'UNKNOWN_ERROR', details = null) {
    super(message);
    this.name = 'AppError';
    this.code = code;
    this.details = details;
  }
}

// Обработчик ошибок API
export const handleApiError = (error) => {
  console.error('API Error:', error);
  
  if (error instanceof AppError) {
    return error;
  }

  if (error.name === 'TypeError' && error.message.includes('fetch')) {
    return new AppError(
      'Ошибка соединения с сервером',
      'NETWORK_ERROR',
      { originalError: error }
    );
  }

  if (error.message.includes('HTTP error')) {
    const statusMatch = error.message.match(/status: (\d+)/);
    const status = statusMatch ? parseInt(statusMatch[1]) : 500;
    
    return new AppError(
      `Ошибка сервера: ${status}`,
      'HTTP_ERROR',
      { status, originalError: error }
    );
  }

  return new AppError(
    'Произошла неизвестная ошибка',
    'UNKNOWN_ERROR',
    { originalError: error }
  );
};

// Показ уведомлений об ошибках
export const showErrorNotification = (error, showNotification) => {
  const appError = handleApiError(error);
  const message = appError.message;
  
  if (showNotification) {
    showNotification(message, 'error');
  }
  
  return appError;
};

// Обработка ошибок с логированием
export const withErrorHandling = (fn, errorHandler = showErrorNotification) => {
  return async (...args) => {
    try {
      return await fn(...args);
    } catch (error) {
      return errorHandler(error);
    }
  };
};
