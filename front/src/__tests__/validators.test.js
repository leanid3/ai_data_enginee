import { validateFile, validatePipelineData, validateETLConfig } from '../utils/validators';

describe('validators', () => {
  describe('validateFile', () => {
    test('валидирует корректный CSV файл', () => {
      const file = new File(['test content'], 'test.csv', { type: 'text/csv' });
      Object.defineProperty(file, 'size', { value: 1024 });
      
      const result = validateFile(file);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('валидирует корректный JSON файл', () => {
      const file = new File(['{"test": "data"}'], 'test.json', { type: 'application/json' });
      Object.defineProperty(file, 'size', { value: 1024 });
      
      const result = validateFile(file);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('валидирует корректный XML файл', () => {
      const file = new File(['<test>data</test>'], 'test.xml', { type: 'application/xml' });
      Object.defineProperty(file, 'size', { value: 1024 });
      
      const result = validateFile(file);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('возвращает ошибку для пустого файла', () => {
      const result = validateFile(null);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Файл не выбран');
    });

    test('возвращает ошибку для неподдерживаемого типа файла', () => {
      const file = new File(['test content'], 'test.txt', { type: 'text/plain' });
      Object.defineProperty(file, 'size', { value: 1024 });
      
      const result = validateFile(file);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Неподдерживаемый тип файла. Поддерживаются: .csv, .json, .xml');
    });

    test('возвращает ошибку для слишком большого файла', () => {
      const file = new File(['test content'], 'test.csv', { type: 'text/csv' });
      Object.defineProperty(file, 'size', { value: 11 * 1024 * 1024 }); // 11MB
      
      const result = validateFile(file);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Размер файла превышает 10MB');
    });
  });

  describe('validatePipelineData', () => {
    test('валидирует корректные данные пайплайна', () => {
      const pipelineData = {
        source: { type: 'file', file: { name: 'test.csv' } },
        target: { system: 'postgresql' },
        transformations: [{ type: 'filter', name: 'test' }],
      };
      
      const result = validatePipelineData(pipelineData);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('возвращает ошибку при отсутствии источника', () => {
      const pipelineData = {
        target: { system: 'postgresql' },
        transformations: [{ type: 'filter', name: 'test' }],
      };
      
      const result = validatePipelineData(pipelineData);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Не указан источник данных');
    });

    test('возвращает ошибку при отсутствии целевой системы', () => {
      const pipelineData = {
        source: { type: 'file', file: { name: 'test.csv' } },
        transformations: [{ type: 'filter', name: 'test' }],
      };
      
      const result = validatePipelineData(pipelineData);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Не указана целевая система');
    });

    test('возвращает ошибку при отсутствии трансформаций', () => {
      const pipelineData = {
        source: { type: 'file', file: { name: 'test.csv' } },
        target: { system: 'postgresql' },
        transformations: [],
      };
      
      const result = validatePipelineData(pipelineData);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Не указаны трансформации данных');
    });
  });

  describe('validateETLConfig', () => {
    test('валидирует корректную конфигурацию ETL с cron', () => {
      const config = {
        schedule: {
          type: 'cron',
          cron: '0 0 * * *',
        },
      };
      
      const result = validateETLConfig(config);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('валидирует корректную конфигурацию ETL с интервалом', () => {
      const config = {
        schedule: {
          type: 'interval',
          interval: 3600,
        },
      };
      
      const result = validateETLConfig(config);
      
      expect(result.isValid).toBe(true);
      expect(result.errors).toHaveLength(0);
    });

    test('возвращает ошибку при отсутствии расписания', () => {
      const config = {};
      
      const result = validateETLConfig(config);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Не указано расписание выполнения');
    });

    test('возвращает ошибку при отсутствии формата расписания', () => {
      const config = {
        schedule: {
          type: 'cron',
        },
      };
      
      const result = validateETLConfig(config);
      
      expect(result.isValid).toBe(false);
      expect(result.errors).toContain('Не указан формат расписания (cron или интервал)');
    });
  });
});
