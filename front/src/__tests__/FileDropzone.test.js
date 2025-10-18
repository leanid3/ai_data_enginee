import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import FileDropzone from '../components/common/FileDropzone';

describe('FileDropzone', () => {
  const mockOnFileSelect = jest.fn();

  beforeEach(() => {
    mockOnFileSelect.mockClear();
  });

  test('рендерится без ошибок', () => {
    render(<FileDropzone onFileSelect={mockOnFileSelect} />);
    
    expect(screen.getByText('Перетащите файл сюда или нажмите для выбора')).toBeInTheDocument();
    expect(screen.getByText('Выбрать файл')).toBeInTheDocument();
  });

  test('отображает поддерживаемые форматы', () => {
    const acceptedTypes = ['.csv', '.json', '.xml'];
    render(<FileDropzone onFileSelect={mockOnFileSelect} acceptedTypes={acceptedTypes} />);
    
    expect(screen.getByText('Поддерживаемые форматы: .csv, .json, .xml')).toBeInTheDocument();
  });

  test('отображает максимальный размер файла', () => {
    const maxSize = 5 * 1024 * 1024; // 5MB
    render(<FileDropzone onFileSelect={mockOnFileSelect} maxSize={maxSize} />);
    
    expect(screen.getByText('Максимальный размер: 5MB')).toBeInTheDocument();
  });

  test('вызывает onFileSelect при выборе файла', () => {
    render(<FileDropzone onFileSelect={mockOnFileSelect} />);
    
    const fileInput = screen.getByDisplayValue('');
    const file = new File(['test content'], 'test.csv', { type: 'text/csv' });
    
    fireEvent.change(fileInput, { target: { files: [file] } });
    
    expect(mockOnFileSelect).toHaveBeenCalledWith(file);
  });

  test('показывает ошибку при выборе неподдерживаемого типа файла', () => {
    render(<FileDropzone onFileSelect={mockOnFileSelect} acceptedTypes={['.csv']} />);
    
    const fileInput = screen.getByDisplayValue('');
    const file = new File(['test content'], 'test.txt', { type: 'text/plain' });
    
    fireEvent.change(fileInput, { target: { files: [file] } });
    
    expect(screen.getByText(/Неподдерживаемый тип файла/)).toBeInTheDocument();
  });

  test('показывает ошибку при выборе слишком большого файла', () => {
    const maxSize = 1024; // 1KB
    render(<FileDropzone onFileSelect={mockOnFileSelect} maxSize={maxSize} />);
    
    const fileInput = screen.getByDisplayValue('');
    const file = new File(['test content'], 'test.csv', { type: 'text/csv' });
    
    // Мокаем размер файла
    Object.defineProperty(file, 'size', { value: 2048 });
    
    fireEvent.change(fileInput, { target: { files: [file] } });
    
    expect(screen.getByText(/Размер файла превышает/)).toBeInTheDocument();
  });
});
