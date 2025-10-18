import config from './config';

class ApiService {
  constructor() {
    this.baseUrl = config.api.baseUrl;
    this.userId = config.app.userId;
  }

  // Общий метод для выполнения запросов
  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const defaultOptions = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    };

    const response = await fetch(url, { ...defaultOptions, ...options });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Загрузка файла
  async uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(
      `${this.baseUrl}${config.api.endpoints.files.upload}?user_id=${this.userId}`,
      {
        method: 'POST',
        body: formData,
      }
    );

    if (!response.ok) {
      throw new Error(`Ошибка загрузки: ${response.status}`);
    }

    return response.json();
  }

  // Запуск анализа
  async startAnalysis(fileId, filePath) {
    const requestData = {
      file_id: fileId,
      user_id: this.userId,
      file_path: filePath,
    };

    return this.request(config.api.endpoints.analysis.start, {
      method: 'POST',
      body: JSON.stringify(requestData),
    });
  }

  // Проверка статуса анализа
  async getAnalysisStatus(analysisId) {
    return this.request(`${config.api.endpoints.analysis.status}/${analysisId}`);
  }

  // Генерация пайплайна
  async generatePipeline(pipelineData) {
    return this.request(config.api.endpoints.pipeline.generate, {
      method: 'POST',
      body: JSON.stringify(pipelineData),
    });
  }

  // Получение пайплайна
  async getPipeline(pipelineId) {
    return this.request(`${config.api.endpoints.pipeline.get}/${pipelineId}`);
  }

  // Выполнение пайплайна
  async executePipeline(pipelineId) {
    return this.request(`${config.api.endpoints.pipeline.execute}/${pipelineId}/execute`, {
      method: 'POST',
    });
  }

  // Удаление пайплайна
  async deletePipeline(pipelineId) {
    return this.request(`${config.api.endpoints.pipeline.delete}/${pipelineId}`, {
      method: 'DELETE',
    });
  }

  // Получение списка пайплайнов
  async getPipelines() {
    return this.request(config.api.endpoints.pipeline.list);
  }
}

export default new ApiService();
