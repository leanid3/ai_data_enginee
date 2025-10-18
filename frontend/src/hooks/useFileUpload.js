import { useState, useCallback } from 'react';
import apiService from '../services/api';
import { validateFile } from '../utils/validators';
import { handleApiError } from '../utils/errorHandler';

export const useFileUpload = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [storagePath, setStoragePath] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleFileChange = useCallback((file) => {
    const validation = validateFile(file);
    
    if (!validation.isValid) {
      setError(validation.errors.join(', '));
      return false;
    }

    setSelectedFile(file);
    setError(null);
    setStoragePath(null); // Сбрасываем предыдущий путь
    return true;
  }, []);

  const uploadFile = useCallback(async () => {
    if (!selectedFile) {
      setError('Файл не выбран');
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await apiService.uploadFile(selectedFile);
      // API возвращает file_id, создаем storage_path для совместимости
      const storagePath = `/files/${result.file_id}/${selectedFile.name}`;
      setStoragePath(storagePath);
      return { ...result, storage_path: storagePath };
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, [selectedFile]);

  const resetFile = useCallback(() => {
    setSelectedFile(null);
    setStoragePath(null);
    setError(null);
    setLoading(false);
  }, []);

  const getFileId = useCallback(() => {
    if (!storagePath) return null;
    const pathParts = storagePath.split('/');
    return pathParts[2]; // Формат: /files/fileId/filename
  }, [storagePath]);

  return {
    selectedFile,
    storagePath,
    loading,
    error,
    handleFileChange,
    uploadFile,
    resetFile,
    getFileId,
  };
};
