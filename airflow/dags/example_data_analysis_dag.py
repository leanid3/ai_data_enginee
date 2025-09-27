"""
Пример DAG для анализа данных
Этот DAG демонстрирует базовую структуру для анализа загруженных файлов
"""

from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from airflow.models import Variable
import requests
import json

# Параметры по умолчанию
default_args = {
    'owner': 'data-engineer',
    'depends_on_past': False,
    'start_date': datetime(2024, 1, 1),
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

# Создание DAG
dag = DAG(
    'data_analysis_pipeline',
    default_args=default_args,
    description='Пайплайн для анализа данных',
    schedule_interval=None,  # Запускается вручную
    catchup=False,
    tags=['data-analysis', 'etl'],
)

def analyze_file_structure(**context):
    """
    Анализирует структуру загруженного файла
    """
    # Получение параметров из контекста
    file_path = context['dag_run'].conf.get('file_path')
    user_id = context['dag_run'].conf.get('user_id')
    
    print(f"Анализируем файл: {file_path} для пользователя: {user_id}")
    
    # Здесь будет логика анализа файла
    # Пока что просто логируем
    analysis_result = {
        'file_path': file_path,
        'user_id': user_id,
        'analysis_time': datetime.now().isoformat(),
        'status': 'completed'
    }
    
    # Сохранение результата в XCom для передачи следующей задаче
    return analysis_result

def generate_recommendations(**context):
    """
    Генерирует рекомендации на основе анализа
    """
    # Получение результата анализа из предыдущей задачи
    analysis_result = context['task_instance'].xcom_pull(task_ids='analyze_structure')
    
    print(f"Генерируем рекомендации для: {analysis_result}")
    
    # Здесь будет логика генерации рекомендаций через LLM
    recommendations = {
        'storage_recommendation': 'PostgreSQL',
        'optimization_tips': ['Добавить индексы', 'Партиционирование'],
        'data_quality_score': 85
    }
    
    return recommendations

def notify_completion(**context):
    """
    Уведомляет о завершении анализа
    """
    # Получение всех результатов
    analysis_result = context['task_instance'].xcom_pull(task_ids='analyze_structure')
    recommendations = context['task_instance'].xcom_pull(task_ids='generate_recommendations')
    
    print("Анализ завершен!")
    print(f"Результат анализа: {analysis_result}")
    print(f"Рекомендации: {recommendations}")
    
    # Здесь можно отправить уведомление в Chat Service
    # или обновить статус в базе данных
    
    return "Анализ успешно завершен"

# Определение задач
analyze_structure = PythonOperator(
    task_id='analyze_structure',
    python_callable=analyze_file_structure,
    dag=dag,
)

generate_recommendations = PythonOperator(
    task_id='generate_recommendations',
    python_callable=generate_recommendations,
    dag=dag,
)

notify_completion = PythonOperator(
    task_id='notify_completion',
    python_callable=notify_completion,
    dag=dag,
)

# Определение зависимостей
analyze_structure >> generate_recommendations >> notify_completion