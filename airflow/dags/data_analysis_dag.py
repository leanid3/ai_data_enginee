"""
DAG для анализа данных после загрузки файла
"""
from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from airflow.providers.http.operators.http import SimpleHttpOperator
from airflow.providers.http.sensors.http import HttpSensor
import json
import requests
import os

# Конфигурация по умолчанию
default_args = {
    'owner': 'aien-backend',
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
    description='Анализ структуры и качества данных',
    schedule_interval=None,  # Запускается вручную
    catchup=False,
    tags=['data-analysis', 'llm', 'quality'],
)

def analyze_data_structure(**context):
    """
    Анализ структуры данных с использованием LLM
    """
    # Получение параметров из контекста
    file_id = context['dag_run'].conf.get('file_id')
    user_id = context['dag_run'].conf.get('user_id')
    file_path = context['dag_run'].conf.get('file_path')
    
    print(f"Анализ файла: {file_id} для пользователя: {user_id}")
    print(f"Путь к файлу: {file_path}")
    
    # Получение файла из MinIO
    minio_endpoint = os.getenv('MINIO_ENDPOINT', 'http://minio:9000')
    minio_access_key = os.getenv('MINIO_ACCESS_KEY', 'minioadmin')
    minio_secret_key = os.getenv('MINIO_SECRET_KEY', 'minioadmin')
    
    # Здесь должна быть логика получения файла из MinIO
    # Для демонстрации создаем фиктивные данные анализа
    
    analysis_result = {
        'file_id': file_id,
        'user_id': user_id,
        'analysis_timestamp': datetime.now().isoformat(),
        'data_structure': {
            'total_rows': 1000,
            'total_columns': 5,
            'column_types': {
                'id': 'integer',
                'name': 'string',
                'email': 'string',
                'age': 'integer',
                'city': 'string'
            },
            'missing_values': {
                'id': 0,
                'name': 5,
                'email': 2,
                'age': 10,
                'city': 3
            },
            'data_quality_score': 0.85
        },
        'recommendations': [
            'Обнаружены пропущенные значения в полях name, age, city',
            'Рекомендуется очистка данных перед загрузкой в целевую систему',
            'Поле email содержит некорректные форматы',
            'Данные подходят для аналитической обработки'
        ],
        'storage_recommendations': {
            'primary_storage': 'PostgreSQL',
            'analytics_storage': 'ClickHouse',
            'backup_storage': 'HDFS',
            'reasoning': 'Данные структурированы и подходят для реляционной БД'
        }
    }
    
    # Сохранение результата анализа
    context['task_instance'].xcom_push(key='analysis_result', value=analysis_result)
    
    return analysis_result

def generate_llm_analysis(**context):
    """
    Генерация анализа с использованием LLM (Ollama)
    """
    file_id = context['dag_run'].conf.get('file_id')
    user_id = context['dag_run'].conf.get('user_id')
    
    # Получение результата предыдущего анализа
    analysis_result = context['task_instance'].xcom_pull(key='analysis_result')
    
    # Подготовка данных для LLM
    llm_prompt = f"""
    Проанализируй следующие данные и дай рекомендации:
    
    Структура данных: {analysis_result['data_structure']}
    
    Дай рекомендации по:
    1. Качеству данных
    2. Оптимальному хранилищу
    3. Предобработке данных
    4. Возможностям аналитики
    """
    
    # Отправка запроса к Ollama
    ollama_url = os.getenv('OLLAMA_URL', 'http://ollama:11434')
    
    try:
        response = requests.post(
            f"{ollama_url}/api/generate",
            json={
                "model": "llama2",
                "prompt": llm_prompt,
                "stream": False
            },
            timeout=30
        )
        
        if response.status_code == 200:
            llm_response = response.json()
            llm_analysis = llm_response.get('response', 'Анализ недоступен')
        else:
            llm_analysis = "Ошибка подключения к LLM"
            
    except Exception as e:
        print(f"Ошибка при обращении к LLM: {e}")
        llm_analysis = "LLM недоступен"
    
    # Объединение результатов
    final_analysis = {
        **analysis_result,
        'llm_analysis': llm_analysis,
        'llm_timestamp': datetime.now().isoformat()
    }
    
    context['task_instance'].xcom_push(key='final_analysis', value=final_analysis)
    
    return final_analysis

def save_analysis_result(**context):
    """
    Сохранение результата анализа в базу данных
    """
    final_analysis = context['task_instance'].xcom_pull(key='final_analysis')
    
    # Здесь должна быть логика сохранения в PostgreSQL
    # Для демонстрации выводим результат
    
    print("=== РЕЗУЛЬТАТ АНАЛИЗА ДАННЫХ ===")
    print(f"Файл ID: {final_analysis['file_id']}")
    print(f"Пользователь: {final_analysis['user_id']}")
    print(f"Время анализа: {final_analysis['analysis_timestamp']}")
    print(f"Качество данных: {final_analysis['data_structure']['data_quality_score']}")
    print(f"Рекомендуемое хранилище: {final_analysis['storage_recommendations']['primary_storage']}")
    print(f"LLM анализ: {final_analysis['llm_analysis']}")
    
    # Сохранение в файл для демонстрации
    with open(f"/tmp/analysis_{final_analysis['file_id']}.json", "w") as f:
        json.dump(final_analysis, f, indent=2, ensure_ascii=False)
    
    return final_analysis

# Определение задач
analyze_structure = PythonOperator(
    task_id='analyze_data_structure',
    python_callable=analyze_data_structure,
    dag=dag,
)

generate_llm_insights = PythonOperator(
    task_id='generate_llm_analysis',
    python_callable=generate_llm_analysis,
    dag=dag,
)

save_results = PythonOperator(
    task_id='save_analysis_result',
    python_callable=save_analysis_result,
    dag=dag,
)

# Определение зависимостей
analyze_structure >> generate_llm_insights >> save_results