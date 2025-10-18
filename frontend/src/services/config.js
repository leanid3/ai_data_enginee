// Конфигурация API и приложения
const config = {
  api: {
    baseUrl: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080',
    endpoints: {
      files: {
        upload: '/api/v1/files/upload',
      },
      analysis: {
        start: '/api/v1/analysis/start',
        status: '/api/v1/analysis/status',
      },
      pipeline: {
        generate: '/api/v1/generate-pipeline',
        get: '/api/v1/pipeline',
        execute: '/api/v1/pipeline',
        delete: '/api/v1/pipeline',
        list: '/api/v1/pipelines',
      }
    }
  },
  app: {
    userId: process.env.REACT_APP_USER_ID || 'test-user-test',
    supportedFileTypes: ['.csv', '.json', '.xml'],
    maxFileSize: 10 * 1024 * 1024, // 10MB
  },
  targetSystems: {
    postgresql: {
      name: 'PostgreSQL',
      description: 'Реляционная база данных для OLTP',
      icon: '🐘',
      features: ['ACID', 'JSON поддержка', 'Полнотекстовый поиск']
    },
    clickhouse: {
      name: 'ClickHouse',
      description: 'Аналитическая СУБД для OLAP',
      icon: '⚡',
      features: ['Высокая производительность', 'Сжатие данных', 'Колоночное хранение']
    },
    hdfs: {
      name: 'HDFS',
      description: 'Распределенная файловая система',
      icon: '🗄️',
      features: ['Масштабируемость', 'Отказоустойчивость', 'Большие данные']
    }
  }
};

export default config;
