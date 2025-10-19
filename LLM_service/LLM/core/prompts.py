

class PromptTemplates:
    @staticmethod
    def router() -> str:
        return (
            """
            Ты - интеллектуальный маршрутизатор для системы автоматизации инженерии данных.
            Анализируй запрос и определи последовательность агентов.

            АГЕНТЫ:
            - DATA_ANALYZER: анализ структуры данных
            - DB_SELECTOR: выбор БД
            - DDL_GENERATOR: создание DDL
            - ETL_BUILDER: построение пайплайнов
            - QUERY_OPTIMIZER: оптимизация
            - REPORT_GENERATOR: отчеты

            Запрос: {user_request}
            Тип операции: {operation_type}

            Верни JSON:
            {{
                "agents_sequence": ["AGENT1", "AGENT2"],
                "reasoning": "обоснование"
            }}
            """
        ).strip()

    @staticmethod
    def analyzer() -> str:
        return (
            """
            Проанализируй данные и определи их характеристики.

            Образец данных:
            {data_sample}

            Определи:
            1. Тип: OLTP/OLAP/TimeSeries/Logs
            2. Объем и частоту обновления
            3. Ключевые поля
            4. Качество данных

            Верни JSON:
            {{
                "data_type": "тип",
                "characteristics": {{
                    "volume": "small/medium/large",
                    "update_frequency": "частота"
                }},
                "structure": {{
                    "key_fields": [],
                    "partition_fields": []
                }},
                "quality_score": 0.0
            }}
            """
        ).strip()

    @staticmethod
    def db_selector() -> str:
        return (
            """
            Выбери оптимальную БД на основе анализа.

            Анализ данных:
            {analysis_results}

            Доступные БД:
            - PostgreSQL: OLTP, транзакции, <1TB
            - ClickHouse: OLAP, аналитика, >100GB
            - HDFS: Big Data, архив, >10TB

            Верни JSON:
            {{
                "recommended_storage": "БД",
                "reasoning": "обоснование",
                "config": {{
                    "partitioning": "стратегия",
                    "replication": 1
                }}
            }}
            """
        ).strip()

    @staticmethod
    def ddl_generator() -> str:
        return (
            """
            Создай DDL для целевой БД.

            БД: {target_db}
            Структура: {data_structure}
            Конфигурация: {storage_config}

            Сгенерируй оптимальную схему с индексами и партицированием.

            Верни JSON:
            {{
                "ddl_scripts": [
                    {{
                        "type": "TABLE",
                        "name": "имя",
                        "script": "CREATE TABLE..."
                    }}
                ]
            }}
            """
        ).strip()

    @staticmethod
    def etl_builder() -> str:
        return (
            """
            Построй Airflow DAG для ETL процесса.

            Источник: {source_config}
            Цель: {target_config}
            Расписание: {schedule}

            Создай DAG с операторами, retry логикой и мониторингом.

            Верни JSON:
            {{
                "dag_config": {{
                    "dag_id": "pipeline_id",
                    "schedule": "0 */1 * * *"
                }},
                "python_code": "from airflow import DAG..."
            }}
            """
        ).strip()

    @staticmethod
    def optimizer() -> str:
        return (
            """
            Оптимизируй запросы для {db_type}.

            Запросы: {queries}
            Объем данных: {data_volume}

            Предложи индексы и оптимизации.

            Верни JSON:
            {{
                "optimizations": [],
                "indexes": [],
                "estimated_improvement": "30%"
            }}
            """
        ).strip()

    @staticmethod
    def reporter() -> str:
        return (
            """
            Создай понятный отчет для пользователя.

            Результаты анализа:
            {all_results}

            Напиши краткое резюме и рекомендации простым языком.

            Верни JSON:
            {{
                "summary": "резюме",
                "recommendations": [],
                "next_steps": []
            }}
            """
        ).strip()