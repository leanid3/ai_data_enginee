import { useState, useCallback } from 'react';
import apiService from '../services/api';
import { validatePipelineData } from '../utils/validators';
import { handleApiError } from '../utils/errorHandler';

export const usePipeline = () => {
  const [pipelines, setPipelines] = useState([]);
  const [currentPipeline, setCurrentPipeline] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const generatePipeline = useCallback(async (pipelineData) => {
    const validation = validatePipelineData(pipelineData);
    
    if (!validation.isValid) {
      setError(validation.errors.join(', '));
      return null;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await apiService.generatePipeline(pipelineData);
      setCurrentPipeline(result);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const getPipeline = useCallback(async (pipelineId) => {
    setLoading(true);
    setError(null);

    try {
      const result = await apiService.getPipeline(pipelineId);
      setCurrentPipeline(result);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const executePipeline = useCallback(async (pipelineId) => {
    setLoading(true);
    setError(null);

    try {
      const result = await apiService.executePipeline(pipelineId);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const deletePipeline = useCallback(async (pipelineId) => {
    setLoading(true);
    setError(null);

    try {
      const result = await apiService.deletePipeline(pipelineId);
      // Удаляем из локального списка
      setPipelines(prev => prev.filter(p => p.id !== pipelineId));
      if (currentPipeline && currentPipeline.id === pipelineId) {
        setCurrentPipeline(null);
      }
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, [currentPipeline]);

  const getPipelines = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await apiService.getPipelines();
      setPipelines(result);
      return result;
    } catch (error) {
      const appError = handleApiError(error);
      setError(appError.message);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const resetPipeline = useCallback(() => {
    setCurrentPipeline(null);
    setError(null);
    setLoading(false);
  }, []);

  return {
    pipelines,
    currentPipeline,
    loading,
    error,
    generatePipeline,
    getPipeline,
    executePipeline,
    deletePipeline,
    getPipelines,
    resetPipeline,
  };
};
