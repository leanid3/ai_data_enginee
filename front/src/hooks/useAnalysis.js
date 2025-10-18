import { useState, useCallback } from 'react';
import apiService from '../services/api';
import { handleApiError } from '../utils/errorHandler';

export const useAnalysis = () => {
  const [analysisId, setAnalysisId] = useState(null);
  const [analysisResults, setAnalysisResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const startAnalysis = useCallback(async (fileId, filePath) => {
    if (!fileId || !filePath) {
      setError('Не указаны параметры для анализа');
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await apiService.startAnalysis(fileId, filePath);
      setAnalysisId(result.analysis_id);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const checkAnalysisStatus = useCallback(async () => {
    if (!analysisId) {
      setError('ID анализа не найден');
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await apiService.getAnalysisStatus(analysisId);
      setAnalysisResults(result);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, [analysisId]);

  const resetAnalysis = useCallback(() => {
    setAnalysisId(null);
    setAnalysisResults(null);
    setError(null);
    setLoading(false);
  }, []);

  const isAnalysisCompleted = useCallback(() => {
    return analysisResults && analysisResults.status === 'completed';
  }, [analysisResults]);

  return {
    analysisId,
    analysisResults,
    loading,
    error,
    startAnalysis,
    checkAnalysisStatus,
    resetAnalysis,
    isAnalysisCompleted,
  };
};
